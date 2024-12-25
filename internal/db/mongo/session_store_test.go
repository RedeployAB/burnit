package mongo

import (
	"context"
	"testing"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSessionStore_Get(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions []db.Session
			id       string
			err      error
		}
		want    db.Session
		wantErr error
	}{
		{
			name: "get session",
			input: struct {
				sessions []db.Session
				id       string
				err      error
			}{
				sessions: []db.Session{
					{
						ID: "1",
					},
				},
				id: "1",
			},
			want: db.Session{
				ID: "1",
			},
		},
		{
			name: "get session - not found",
			input: struct {
				sessions []db.Session
				id       string
				err      error
			}{
				sessions: []db.Session{},
				id:       "1",
			},
			wantErr: dberrors.ErrSessionNotFound,
		},
		{
			name: "get session - error",
			input: struct {
				sessions []db.Session
				id       string
				err      error
			}{
				sessions: []db.Session{},
				id:       "1",
				err:      errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &sessionStore{
				client: &stubMongoClient{
					sessions: test.input.sessions,
					err:      test.input.err,
				},
			}

			got, gotErr := store.Get(context.Background(), test.input.id)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}

		})
	}
}

func TestSessionStore_GetByCSRFToken(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sessions []db.Session
			token    string
			err      error
		}
		want    db.Session
		wantErr error
	}{
		{
			name: "get session",
			input: struct {
				sessions []db.Session
				token    string
				err      error
			}{
				sessions: []db.Session{
					{
						ID: "1",
						CSRF: db.CSRF{
							Token: "token",
						},
					},
				},
				token: "token",
			},
			want: db.Session{
				ID: "1",
				CSRF: db.CSRF{
					Token: "token",
				},
			},
		},
		{
			name: "get session - not found",
			input: struct {
				sessions []db.Session
				token    string
				err      error
			}{
				sessions: []db.Session{},
				token:    "1",
			},
			wantErr: dberrors.ErrSessionNotFound,
		},
		{
			name: "get session - error",
			input: struct {
				sessions []db.Session
				token    string
				err      error
			}{
				sessions: []db.Session{},
				token:    "1",
				err:      errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &sessionStore{
				client: &stubMongoClient{
					sessions: test.input.sessions,
					err:      test.input.err,
				},
			}

			got, gotErr := store.GetByCSRFToken(context.Background(), test.input.token)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("GetByCSRFToken() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("GetByCSRFToken() = unexpected error (-want +got)\n%s\n", diff)
			}

		})
	}
}
