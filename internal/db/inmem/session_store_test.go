package inmem

import (
	"context"
	"sync"
	"testing"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSessionStore(t *testing.T) {
	got := NewSessionStore()
	if diff := cmp.Diff(&sessionStore{sessions: make(map[string]db.Session), sessionCSRF: make(map[string]string), mu: sync.RWMutex{}}, got, cmp.AllowUnexported(sessionStore{}), cmpopts.IgnoreFields(sessionStore{}, "mu", "stopCh")); diff != "" {
		t.Errorf("NewInMemoryStore() = unexpected result (-want +got)\n%s\n", diff)
	}
}

func TestInMemoryStore_Get(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			sessions map[string]db.Session
			id       string
		}
		want    db.Session
		wantErr error
	}{
		{
			name: "Get session",
			input: struct {
				sessions map[string]db.Session
				id       string
			}{
				sessions: map[string]db.Session{
					"test": {
						ID:        "test",
						ExpiresAt: n.Add(1),
					},
				},
				id: "test",
			},
			want: db.Session{
				ID:        "test",
				ExpiresAt: n.Add(1),
			},
		},
		{
			name: "Session not found",
			input: struct {
				sessions map[string]db.Session
				id       string
			}{
				sessions: map[string]db.Session{},
				id:       "test",
			},
			wantErr: dberrors.ErrSessionNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &sessionStore{
				sessions: test.input.sessions,
				mu:       sync.RWMutex{},
			}

			got, gotErr := s.Get(context.Background(), test.input.id)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}

		})
	}
}

func TestInMemoryStore_GetByCSRFToken(t *testing.T) {
	n := now()
	var tests = []struct {
		name  string
		input struct {
			sessions    map[string]db.Session
			sessionCSRF map[string]string
			token       string
		}
		want    db.Session
		wantErr error
	}{
		{
			name: "Get session",
			input: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
				token       string
			}{
				sessions: map[string]db.Session{
					"test": {
						ID:        "test",
						ExpiresAt: n.Add(1),
						CSRF: db.CSRF{
							Token:     "token",
							ExpiresAt: n.Add(1),
						},
					},
				},
				sessionCSRF: map[string]string{
					"token": "test",
				},
				token: "token",
			},
			want: db.Session{
				ID:        "test",
				ExpiresAt: n.Add(1),
				CSRF: db.CSRF{
					Token:     "token",
					ExpiresAt: n.Add(1),
				},
			},
		},
		{
			name: "Session not found - no token",
			input: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
				token       string
			}{
				sessions:    map[string]db.Session{},
				sessionCSRF: map[string]string{},
				token:       "token",
			},
			wantErr: dberrors.ErrSessionNotFound,
		},
		{
			name: "Session not found - no session",
			input: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
				token       string
			}{
				sessions: map[string]db.Session{},
				sessionCSRF: map[string]string{
					"token": "test",
				},
				token: "token",
			},
			wantErr: dberrors.ErrSessionNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &sessionStore{
				sessions:    test.input.sessions,
				sessionCSRF: test.input.sessionCSRF,
				mu:          sync.RWMutex{},
			}

			got, gotErr := s.GetByCSRFToken(context.Background(), test.input.token)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}

		})
	}
}

func TestInMemoryStore_Upsert(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions    map[string]db.Session
			sessionCSRF map[string]string
			session     db.Session
		}
		want    db.Session
		wantErr error
	}{
		{
			name: "Upsert session",
			input: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
				session     db.Session
			}{
				sessions:    map[string]db.Session{},
				sessionCSRF: map[string]string{},
				session: db.Session{
					ID: "test",
					CSRF: db.CSRF{
						Token: "token",
					},
				},
			},
			want: db.Session{
				ID: "test",
				CSRF: db.CSRF{
					Token: "token",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &sessionStore{
				sessions:    test.input.sessions,
				sessionCSRF: test.input.sessionCSRF,
				mu:          sync.RWMutex{},
			}

			got, gotErr := s.Upsert(context.Background(), test.input.session)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Set() = unexpected error (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Set() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestInMemoryStore_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions map[string]db.Session
			id       string
		}
		wantErr error
	}{
		{
			name: "Delete session",
			input: struct {
				sessions map[string]db.Session
				id       string
			}{
				sessions: map[string]db.Session{
					"test": {
						ID: "test",
					},
				},
				id: "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &sessionStore{
				sessions: test.input.sessions,
				mu:       sync.RWMutex{},
			}

			gotErr := s.Delete(context.Background(), test.input.id)

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

func TestInMemoryStore_DeleteExpired(t *testing.T) {
	n := now()

	var tests = []struct {
		name  string
		input struct {
			sessions    map[string]db.Session
			sessionCSRF map[string]string
		}
		want struct {
			sessions    map[string]db.Session
			sessionCSRF map[string]string
		}
		wantErr error
	}{
		{
			name: "Delete expired sessions",
			input: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
			}{
				sessions: map[string]db.Session{
					"test": {
						ID:        "test",
						ExpiresAt: n.Add(-1),
						CSRF: db.CSRF{
							Token: "token",
						},
					},
				},
				sessionCSRF: map[string]string{
					"token": "test",
				},
			},
			want: struct {
				sessions    map[string]db.Session
				sessionCSRF map[string]string
			}{
				sessions:    map[string]db.Session{},
				sessionCSRF: map[string]string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &sessionStore{
				sessions:    test.input.sessions,
				sessionCSRF: test.input.sessionCSRF,
				mu:          sync.RWMutex{},
			}

			gotErr := s.DeleteExpired(context.Background())

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}

			gotSessions := s.sessions
			if diff := cmp.Diff(test.want.sessions, gotSessions); diff != "" {
				t.Errorf("DeleteExpired() = unexpected result (-want +got)\n%s\n", diff)
			}

			gotSessionCSRF := s.sessionCSRF
			if diff := cmp.Diff(test.want.sessionCSRF, gotSessionCSRF); diff != "" {
				t.Errorf("DeleteExpired() = unexpected result (-want +got)\n%s\n", diff)
			}

		})
	}
}
