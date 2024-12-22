package session

import (
	"sync"
	"time"
)

const (
	// defaultInMemoryStoreCleanupInterval is the default interval for cleaning up expired sessions
	// and CSRF tokens in the in-memory store.
	defaultInMemoryStoreCleanupInterval = time.Minute
)

// Store is an interface for storing sessions.
type Store interface {
	// Get a session by its ID.
	Get(id string) (Session, error)
	// Set a session.
	Set(id string, session Session) error
	// Delete a session by its ID.
	Delete(id string) error
	// Cleanup runs a cleanup routine to delete expired sessions.
	Cleanup() chan error
	// Close the store.
	Close() error
}

// sessions is a map of session IDs to sessions.
type sessions map[string]Session

// inMemoryStore is an in-memory store for sessions.
type inMemoryStore struct {
	sessions sessions
	stopCh   chan struct{}
	mu       sync.RWMutex
}

// NewInMemoryStore creates a new in-memory session store.
func NewInMemoryStore() *inMemoryStore {
	s := &inMemoryStore{
		sessions: make(sessions),
		mu:       sync.RWMutex{},
		stopCh:   make(chan struct{}),
	}
	return s
}

// Get a session by its ID.
func (s *inMemoryStore) Get(id string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return Session{}, ErrSessionNotFound
	}
	if session.ExpiresAt().Before(now()) {
		delete(s.sessions, id)
		return Session{}, ErrSessionExpired
	}
	return session, nil
}

// Set a session.
func (s *inMemoryStore) Set(id string, sessions Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[id] = sessions
	return nil
}

// Delete a session by its ID.
func (s *inMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
	return nil
}

// Close the store.
func (s *inMemoryStore) Close() error {
	s.stopCh <- struct{}{}
	return nil
}

// Cleanup deletes expired sessions. The channel returned
// by the in-memory store has no real use since it does not
// return any errors.
func (s *inMemoryStore) Cleanup() chan error {
	errCh := make(chan error)
	go func() {
		for {
			select {
			case <-time.After(defaultInMemoryStoreCleanupInterval):
				s.mu.Lock()
				for id := range s.sessions {
					if s.sessions[id].Expired() {
						delete(s.sessions, id)
					}
					if !s.sessions[id].CSRF().IsEmpty() && s.sessions[id].CSRF().Expired() {
						sess := s.sessions[id]
						(&sess).DeleteCSRF()
						s.sessions[id] = sess
					}
				}
				s.mu.Unlock()
			case <-s.stopCh:
				close(errCh)
				close(s.stopCh)
				return
			}
		}
	}()
	return errCh
}
