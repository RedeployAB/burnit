package session

import "sync"

// sessions is a map of session IDs to sessions.
type sessions map[string]Session

// inMemoryStore is an in-memory store for sessions.
type inMemoryStore struct {
	sessions sessions
	mu       sync.RWMutex
}

// NewInMemoryStore creates a new in-memory store.
func NewInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		sessions: make(sessions),
		mu:       sync.RWMutex{},
	}
}

// Get returns the session with the given ID.
func (s *inMemoryStore) Get(id string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.sessions[id]
	if !ok {
		return Session{}, ErrSessionNotFound
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
