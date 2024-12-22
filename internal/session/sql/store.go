package sql

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/session"
)

const (
	// defaultStoreCleanupInterval is the default interval for cleaning up expired sessions
	// and CSRF tokens in the store.
	defaultStoreCleanupInterval = 5 * time.Second //time.Minute
	// defaultStoreTimeout is the default timeout for Store actions.
	defaultStoreTimeout = 5 * time.Second
)

// store is a SQL implementation of a session store.
type store struct {
	sessions        db.SessionRepository
	timeout         time.Duration
	cleanupInterval time.Duration
	stopCh          chan struct{}
}

// StoreOptions contains options for the store.
type StoreOptions struct {
	Timeout         time.Duration
	CleanupInterval time.Duration
	Logger          log.Logger
}

// StoreOption is a function that configures the Store.
type StoreOption func(o *StoreOptions)

// NewStore creates a new SQL session store.
func NewStore(repo db.SessionRepository, options ...StoreOption) (*store, error) {
	opts := StoreOptions{
		Timeout:         defaultStoreTimeout,
		CleanupInterval: defaultStoreCleanupInterval,
	}
	for _, option := range options {
		option(&opts)
	}
	if opts.Logger == nil {
		opts.Logger = log.New()
	}

	s := &store{
		sessions:        repo,
		timeout:         opts.Timeout,
		cleanupInterval: opts.CleanupInterval,
		stopCh:          make(chan struct{}),
	}

	//go s.cleanup()
	return s, nil
}

// Get a session by its ID.
func (s store) Get(id string) (session.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	_sess, err := s.sessions.Get(ctx, id)
	if err != nil {
		if errors.Is(err, dberrors.ErrSessionNotFound) {
			return session.Session{}, session.ErrSessionNotFound
		}
		return session.Session{}, err
	}

	sess := session.NewSession(withSession(_sess))
	if sess.ExpiresAt().Before(now()) {
		if err := s.sessions.Delete(ctx, sess.ID()); err != nil {
			return session.Session{}, err
		}
	}

	return sess, nil
}

// Set a session.
func (s store) Set(id string, sess session.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if _, err := s.sessions.Upsert(ctx, db.Session{
		ID:        sess.ID(),
		ExpiresAt: sess.ExpiresAt(),
		CSRF: db.CSRF{
			Token:     sess.CSRF().Token(),
			ExpiresAt: sess.CSRF().ExpiresAt(),
		},
	}); err != nil {
		return err
	}

	return nil
}

// Delete a session by its ID.
func (s store) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.sessions.Delete(ctx, id)
}

// Close the store.
func (s store) Close() error {
	s.stopCh <- struct{}{}
	return s.sessions.Close()
}

// Cleanup runs a cleanup routine to delete expired sessions.
// It returns a channel to receive errors. When the store is
// closed with Close, the channel is closed as it is not
// intended for further use.
func (s *store) Cleanup() chan error {
	errCh := make(chan error)
	go func() {
		for {
			select {
			case <-time.After(s.cleanupInterval):
				ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
				defer cancel()

				if err := s.sessions.DeleteExpired(ctx); err != nil {
					if !errors.Is(err, dberrors.ErrSessionsNotDeleted) {
						errCh <- err
					}
				}
			case <-s.stopCh:
				close(errCh)
				close(s.stopCh)
				return
			}
		}
	}()
	return errCh
}

// withSession is an option to create a new session from
// a db.Session.
func withSession(s db.Session) session.SessionOption {
	return func(o *session.SessionOptions) {
		o.ID = s.ID
		o.ExpiresAt = s.ExpiresAt
		o.CSRF = session.NewCSRF(func(o *session.CSRFOptions) {
			o.Token = s.CSRF.Token
			o.ExpiresAt = s.CSRF.ExpiresAt
		})
	}
}

var now = func() time.Time {
	return time.Now().UTC()
}
