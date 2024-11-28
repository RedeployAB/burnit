package session

import (
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewInMemoryStore(t *testing.T) {
	got := NewInMemoryStore()
	if diff := cmp.Diff(&inMemoryStore{sessions: make(sessions), mu: sync.RWMutex{}}, got, cmp.AllowUnexported(inMemoryStore{}), cmpopts.IgnoreFields(inMemoryStore{}, "mu", "stopCh")); diff != "" {
		t.Errorf("NewInMemoryStore() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestInMemoryStore_Get(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			sessions sessions
			id       string
		}
		want    Session
		wantErr error
	}{
		{
			name: "Get session",
			input: struct {
				sessions sessions
				id       string
			}{
				sessions: sessions{
					"test": Session{
						id:        "test",
						expiresAt: n.Add(1),
					},
				},
				id: "test",
			},
			want: Session{
				id:        "test",
				expiresAt: n.Add(1),
			},
		},
		{
			name: "Session not found",
			input: struct {
				sessions sessions
				id       string
			}{
				sessions: sessions{},
				id:       "test",
			},
			wantErr: ErrSessionNotFound,
		},
		{
			name: "Session expired",
			input: struct {
				sessions sessions
				id       string
			}{
				sessions: sessions{
					"test": Session{
						id:        "test",
						expiresAt: n.Add(-1),
					},
				},
				id: "test",
			},
			wantErr: ErrSessionExpired,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &inMemoryStore{
				sessions: test.input.sessions,
				mu:       sync.RWMutex{},
			}

			got, gotErr := s.Get(test.input.id)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSRF{})); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}

		})
	}
}

func TestInMemoryStore_Set(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions sessions
			session  Session
		}
		want    Session
		wantErr error
	}{
		{
			name: "Set session",
			input: struct {
				sessions sessions
				session  Session
			}{
				sessions: sessions{},
				session: Session{
					id: "test",
				},
			},
			want: Session{
				id: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &inMemoryStore{
				sessions: test.input.sessions,
				mu:       sync.RWMutex{},
			}

			gotErr := s.Set(test.input.session.id, test.input.session)
			got := s.sessions[test.input.session.id]

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() = unexpected error (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSRF{})); diff != "" {
				t.Errorf("Set() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestInMemoryStore_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions sessions
			id       string
		}
		wantErr error
	}{
		{
			name: "Delete session",
			input: struct {
				sessions sessions
				id       string
			}{
				sessions: sessions{
					"test": Session{
						id: "test",
					},
				},
				id: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &inMemoryStore{
				sessions: test.input.sessions,
				mu:       sync.RWMutex{},
			}

			gotErr := s.Delete(test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}

			_, ok := s.sessions[test.input.id]
			if ok {
				t.Errorf("Delete() = session not deleted")
			}

		})
	}
}
