package inmem

import (
	"context"
	"sync"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

// sessionStore is an in-memory store for sessions.
type sessionStore struct {
	sessions    map[string]db.Session
	sessionCSRF map[string]string
	mu          sync.RWMutex
}

// NewSessionStore creates a new in-memory session store.
func NewSessionStore() *sessionStore {
	s := &sessionStore{
		sessions:    make(map[string]db.Session),
		sessionCSRF: make(map[string]string),
		mu:          sync.RWMutex{},
	}
	return s
}

// Get a session by its ID.
func (s *sessionStore) Get(ctx context.Context, id string) (db.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return db.Session{}, dberrors.ErrSessionNotFound
	}

	return session, nil
}

// Get a session by its CSRF token.
func (s *sessionStore) GetByCSRFToken(ctx context.Context, token string) (db.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.sessionCSRF[token]
	if !ok {
		return db.Session{}, dberrors.ErrSessionNotFound
	}
	session, ok := s.sessions[id]
	if !ok {
		return db.Session{}, dberrors.ErrSessionNotFound
	}

	return session, nil
}

// Upsert a session. Create the session if it does not exist, otherwise
// update it.
func (s *sessionStore) Upsert(ctx context.Context, session db.Session) (db.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID] = db.Session{
		ID:        session.ID,
		ExpiresAt: session.ExpiresAt,
		CSRF: db.CSRF{
			Token:     session.CSRF.Token,
			ExpiresAt: session.CSRF.ExpiresAt,
		},
	}
	if len(session.CSRF.Token) > 0 {
		s.sessionCSRF[session.CSRF.Token] = session.ID
	}

	return s.sessions[session.ID], nil
}

// Delete a session by its ID.
func (s *sessionStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return dberrors.ErrSessionNotFound
	}

	delete(s.sessionCSRF, session.CSRF.Token)
	delete(s.sessions, session.ID)

	return nil
}

// DeleteByCSRFToken deletes a session by its CSRF token.
func (s *sessionStore) DeleteByCSRFToken(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.sessionCSRF[token]
	if !ok {
		return dberrors.ErrSessionNotFound
	}

	session, ok := s.sessions[id]
	if !ok {
		return dberrors.ErrSessionNotFound
	}

	delete(s.sessionCSRF, session.CSRF.Token)
	delete(s.sessions, session.ID)

	return nil
}

// DeleteExpired deletes all expired sessions.
func (s *sessionStore) DeleteExpired(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range s.sessions {
		if session.ExpiresAt.Before(now()) {
			delete(s.sessionCSRF, session.CSRF.Token)
			delete(s.sessions, session.ID)
		}
	}

	return nil
}

// Close the store and its underlying connections.
func (s *sessionStore) Close() error {
	return nil
}
