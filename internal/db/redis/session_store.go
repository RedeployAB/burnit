package redis

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

const (
	// sessionPrefix is the key prefix for sessions.
	sessionPrefix = "session:"
	// sessionCSRFPrefix is the key prefix for CSRF tokens.
	sessionCSRFPrefix = "session-csrf:"
)

// sessionStore is a Redis implementation of a SessionStore.
type sessionStore struct {
	client Client
}

// SessionStoreOptions is the options for the SessionStore.
type SessionStoreOptions struct{}

// SessionStoreOption is a function that sets options for the SessionStore.
type SessionStoreOption func(o *SessionStoreOptions)

// NewSessionStore creates and configures a new SessionStore.
func NewSessionStore(client Client, options ...SessionStoreOption) (*sessionStore, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}

	opts := SessionStoreOptions{}
	for _, option := range options {
		option(&opts)
	}

	return &sessionStore{
		client: client,
	}, nil
}

// Get a session by its ID.
func (s sessionStore) Get(ctx context.Context, id string) (db.Session, error) {
	data, err := s.client.HGet(ctx, sessionPrefix+id)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}
	return sessionFromMap(data)
}

// GetByCSRFToken gets a session by its CSRF token.
func (s sessionStore) GetByCSRFToken(ctx context.Context, token string) (db.Session, error) {
	data, err := s.client.Get(ctx, sessionCSRFPrefix+token)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}
	return s.Get(ctx, string(data))
}

// Upsert a session.
func (s sessionStore) Upsert(ctx context.Context, session db.Session) (db.Session, error) {
	var token string
	if sess, err := s.Get(ctx, session.ID); err == nil {
		token = sess.CSRF.Token
	} else if !errors.Is(err, dberrors.ErrSessionNotFound) {
		return db.Session{}, err
	}

	result, err := s.client.WithTransaction(ctx, func(tx Tx) {
		if len(token) > 0 {
			tx.Delete(ctx, sessionCSRFPrefix+token)
		}
		tx.HSet(ctx, sessionPrefix+session.ID, sessionToMap(&session))
		tx.Expire(ctx, sessionPrefix+session.ID, time.Until(session.ExpiresAt))
		if len(session.CSRF.Token) > 0 {
			tx.Set(ctx, sessionCSRFPrefix+session.CSRF.Token, []byte(session.ID), time.Until(session.CSRF.ExpiresAt))
		}
		tx.HGet(ctx, sessionPrefix+session.ID)
	})
	if err != nil {
		return db.Session{}, err
	}

	data := result.LastMap()
	if data == nil {
		return db.Session{}, dberrors.ErrSessionNotFound
	}
	return sessionFromMap(data)
}

// Delete a session by its ID.
func (s sessionStore) Delete(ctx context.Context, id string) error {
	session, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if _, err = s.client.WithTransaction(ctx, func(tx Tx) {
		if len(session.CSRF.Token) > 0 {
			tx.Delete(ctx, sessionCSRFPrefix+session.CSRF.Token)
		}
		tx.Delete(ctx, sessionPrefix+session.ID)
	}); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return dberrors.ErrSessionNotFound
		}
		return err
	}
	return nil
}

// DeleteByCSRFToken deletes a session by its CSRF token.
func (s sessionStore) DeleteByCSRFToken(ctx context.Context, token string) error {
	session, err := s.GetByCSRFToken(ctx, token)
	if err != nil {
		return err
	}

	if _, err := s.client.WithTransaction(ctx, func(tx Tx) {
		tx.Delete(ctx, sessionCSRFPrefix+token)
		tx.Delete(ctx, sessionPrefix+session.ID)
	}); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return dberrors.ErrSessionNotFound
		}
		return err
	}
	return nil
}

// DeleteExpired deletes all expired sessions. This is a no-op for Redis
// since Redis handles expiration automatically.
func (s sessionStore) DeleteExpired(ctx context.Context) error {
	return nil
}

// Close the store and its underlying connections.
func (s sessionStore) Close() error {
	return s.client.Close()
}

// sessionToMap creates a map from the provided session.
func sessionToMap(session *db.Session) map[string]any {
	return map[string]any{
		"id":              session.ID,
		"expires_at":      session.ExpiresAt,
		"csrf_token":      session.CSRF.Token,
		"csrf_expires_at": session.CSRF.ExpiresAt,
	}
}

// sessionFromMap creates a db.Session from the provided map.
func sessionFromMap(session map[string]string) (db.Session, error) {
	expiresAt, err := time.Parse(time.RFC3339, session["expires_at"])
	if err != nil {
		return db.Session{}, err
	}
	csrfExpiresAt, err := time.Parse(time.RFC3339, session["csrf_expires_at"])
	if err != nil {
		return db.Session{}, err
	}
	return db.Session{
		ID:        session["id"],
		ExpiresAt: expiresAt,
		CSRF: db.CSRF{
			Token:     session["csrf_token"],
			ExpiresAt: csrfExpiresAt,
		},
	}, nil
}
