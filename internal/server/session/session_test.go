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
			name: "with options - provided CSFR",
			options: []SessionOption{
				WithExpiresAt(time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)),
				WithCSFR(CSFR{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				}),
			},
			want: Session{
				id:        "test-uuid",
				expiresAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				csfr: CSFR{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "with options - provided CSFR options",
			options: []SessionOption{
				WithExpiresAt(time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)),
				WithCSFROptions(WithCSFRExpiresAt(time.Date(2024, 1, 1, 0, 20, 0, 0, time.UTC))),
			},
			want: Session{
				id:        "test-uuid",
				expiresAt: time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
				csfr: CSFR{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 20, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewSession(test.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSFR{})); diff != "" {
				t.Errorf("NewSession(%v) = unexpected result (-want +got)\n%s\n", test.options, diff)
			}
		})
	}
}

func TestSession_SetCSFR(t *testing.T) {
	var tests = []struct {
		name string
		csfr CSFR
		want Session
	}{
		{
			name: "default",
			csfr: CSFR{
				token:     "test-random",
				expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
			},
			want: Session{
				csfr: CSFR{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := Session{}
			got := s.SetCSFR(test.csfr)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSFR{})); diff != "" {
				t.Errorf("SetCSFR(%v) = unexpected result (-want +got)\n%s\n", test.csfr, diff)
			}
		})
	}
}

func TestSession_DeleteCSFR(t *testing.T) {
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
				csfr: CSFR{
					token:     "test-random",
					expiresAt: time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC),
				},
			}
			got := s.DeleteCSFR()

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(Session{}, CSFR{})); diff != "" {
				t.Errorf("DeleteCSFR() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
