package session

import (
	"sync"
	"time"
)

// Store is an interface for storing sessions.
type Store interface {
	Get(id string) (Session, error)
	Set(id string, session Session) error
	Delete(id string) error
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

// NewInMemoryStore creates a new in-memory store.
func NewInMemoryStore() *inMemoryStore {
	s := &inMemoryStore{
		sessions: make(sessions),
		mu:       sync.RWMutex{},
		stopCh:   make(chan struct{}),
	}

	go s.cleanup()
	return s
}

// Get returns the session with the given ID.
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

// Set sets the session with the given ID.
func (s *inMemoryStore) Set(id string, sessions Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[id] = sessions
	return nil
}

// Delete deletes the session with the given ID.
func (s *inMemoryStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)
	return nil
}

// Close closes the in-memory store.
func (s *inMemoryStore) Close() error {
	s.stopCh <- struct{}{}
	return nil
}

// cleanup deletes expired sessions.
func (s *inMemoryStore) cleanup() {
	for {
		select {
		case <-time.After(time.Second * 5):
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
			close(s.stopCh)
			return
		}
	}
}
