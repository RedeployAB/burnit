package auth

// TokenStore is an interface that represents
// a method of getting, setting and comparing tokens.
type TokenStore interface {
	Get(token string) (string, bool)
	Set(token, user string)
	Verify(token string) bool
}

// MemoryTokenStore implements TokenStore.
// It stores tokens and APIKeys in memory.
type MemoryTokenStore struct {
	Tokens map[string]Token
}

// Token represents an Key/Token.
type Token struct {
	User string
}

// Get searchs the memory token store for a key.
func (m *MemoryTokenStore) Get(token string) (string, bool) {
	t, ok := m.Tokens[token]
	if !ok {
		return "", ok
	}
	return t.User, ok
}

// Set sets a token in the TokenStore.
func (m *MemoryTokenStore) Set(token, user string) {
	m.Tokens[token] = Token{User: user}
}

// Verify compares the key in the memory token store with
// the key passed in as inKey.
func (m *MemoryTokenStore) Verify(token string) bool {
	_, ok := m.Get(token)
	if ok {
		return true
	}
	return false
}

// NewMemoryTokenStore initializes and returns a MemoryTokenStore.
func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{Tokens: map[string]Token{}}
}
