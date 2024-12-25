package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

const (
	// defaultTimeout is the default timeout for Store actions.
	defaultTimeout = 5 * time.Second
	// defaultCleanupInterval is the default interval for cleaning up expired sessions
	// and CSRF tokens in the store.
	defaultCleanupInterval = time.Minute
)

// Service is an interface for handling sessions.
type Service interface {
	// Get a session by its ID.
	Get(options ...GetOption) (Session, error)
	// Set a session.
	Set(session Session) error
	// Delete a session by its ID.
	Delete(id string) error
	// Cleanup runs a cleanup routine to delete expired sessions.
	Cleanup() chan error
	// Close the service.
	Close() error
}

// service is a session service.
type service struct {
	sessions        db.SessionStore
	timeout         time.Duration
	cleanupInterval time.Duration
	stopCh          chan struct{}
}

// ServiceOption is a function that sets options for the service.
type ServiceOption func(s *service)

// NewService creats a new session service.
func NewService(store db.SessionStore, options ...ServiceOption) (*service, error) {
	if store == nil {
		return nil, ErrNilStore
	}

	svc := &service{
		sessions:        store,
		timeout:         defaultTimeout,
		cleanupInterval: defaultCleanupInterval,
		stopCh:          make(chan struct{}),
	}

	for _, option := range options {
		option(svc)
	}

	return svc, nil
}

// GetOptions contains options for getting a session.
type GetOptions struct {
	ID        string
	CSRFToken string
}

// GetOption is a function that configures the GetOptions.
type GetOption func(o *GetOptions)

// Get a session by its ID.
func (s service) Get(options ...GetOption) (Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	opts := GetOptions{}
	for _, option := range options {
		option(&opts)
	}

	var getFunc func(context.Context, string) (db.Session, error)
	var param string
	switch {
	case len(opts.ID) > 0:
		getFunc = s.sessions.Get
		param = opts.ID
	case len(opts.CSRFToken) > 0:
		getFunc = s.sessions.GetByCSRFToken
		param = opts.CSRFToken
	default:
		return Session{}, errors.New("no ID or CSRF token provided")
	}

	sess, err := getFunc(ctx, param)
	if err != nil {
		if errors.Is(err, dberrors.ErrSessionNotFound) {
			return Session{}, ErrSessionNotFound
		}
		return Session{}, err
	}

	session := NewSession(withSession(sess))
	if session.ExpiresAt().Before(now()) {
		if err := s.sessions.Delete(ctx, session.ID()); err != nil {
			return Session{}, err
		}
	}

	return session, nil
}

// Set a session.
func (s service) Set(session Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if _, err := s.sessions.Upsert(ctx, db.Session{
		ID:        session.ID(),
		ExpiresAt: session.ExpiresAt(),
		CSRF: db.CSRF{
			Token:     session.CSRF().Token(),
			ExpiresAt: session.CSRF().ExpiresAt(),
		},
	}); err != nil {
		return err
	}

	return nil
}

// Delete a session by its ID.
func (s service) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.sessions.Delete(ctx, id)
}

// Cleanup runs a cleanup routine to delete expired sessions.
// It returns a channel to receive errors. When the service is
// closed with Close, the channel is closed as it is not
// intended for further use.
func (s *service) Cleanup() chan error {
	errCh := make(chan error)
	go func() {
		defer func() {
			close(errCh)
			close(s.stopCh)
		}()
		for {
			select {
			case <-time.After(s.cleanupInterval):
				ctx, cancel := context.WithTimeout(context.Background(), s.timeout)

				if err := s.sessions.DeleteExpired(ctx); err != nil {
					if !errors.Is(err, dberrors.ErrSessionsNotDeleted) {
						errCh <- fmt.Errorf("session store: %w", err)
					}
				}
				cancel()
			case <-s.stopCh:
				return
			}
		}
	}()
	return errCh
}

// Close the service.
func (s *service) Close() error {
	s.stopCh <- struct{}{}
	return s.sessions.Close()
}

// withSession is an option to create a new session from
// a db.Session.
func withSession(s db.Session) SessionOption {
	return func(o *SessionOptions) {
		o.ID = s.ID
		o.ExpiresAt = s.ExpiresAt
		o.CSRF = NewCSRF(func(o *CSRFOptions) {
			o.Token = s.CSRF.Token
			o.ExpiresAt = s.CSRF.ExpiresAt
		})
	}
}
