package session

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewSession(t *testing.T) {
	newUUID = func() string {
		return "test-uuid"
	}

	now = func() time.Time {
		return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	randomString = func() string {
		return "test-random"
	}

	var tests = []struct {
		name    string
		options []SessionOption
		want    Session
	}{
		{
			name:    "default",
			options: nil,
			want: Session{
				id:        "test-uuid",
				expiresAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "with options - provided CSRF",
			options: []SessionOption{
				WithExpiresAt(time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)),
				WithCSRF(CSRF{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				}),
			},
			want: Session{
				id:        "test-uuid",
				expiresAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				csrf: CSRF{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "with options - provided CSRF options",
			options: []SessionOption{
				WithExpiresAt(time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)),
				WithCSRFOptions(WithCSRFExpiresAt(time.Date(2024, 1, 1, 0, 20, 0, 0, time.UTC))),
			},
			want: Session{
				id:        "test-uuid",
				expiresAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				csrf: CSRF{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 20, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewSession(test.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSRF{})); diff != "" {
				t.Errorf("NewSession(%v) = unexpected result (-want +got)\n%s\n", test.options, diff)
			}
		})
	}
}

func TestSession_SetCSRF(t *testing.T) {
	var tests = []struct {
		name string
		csrf CSRF
		want Session
	}{
		{
			name: "default",
			csrf: CSRF{
				token:     "test-random",
				expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
			},
			want: Session{
				csrf: CSRF{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Session{}
			got := s.SetCSRF(test.csrf)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSRF{})); diff != "" {
				t.Errorf("SetCSRF(%v) = unexpected result (-want +got)\n%s\n", test.csrf, diff)
			}
		})
	}
}

func TestSession_DeleteCSRF(t *testing.T) {
	var tests = []struct {
		name string
		want Session
	}{
		{
			name: "default",
			want: Session{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Session{
				csrf: CSRF{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			}
			got := s.DeleteCSRF()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSRF{})); diff != "" {
				t.Errorf("DeleteCSRF() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
