package api

import (
	"fmt"
	"strings"
	"time"
)

var (
	// ErrInvalidTimeFormat is returned when the time format is invalid.
	ErrInvalidTimeFormat = fmt.Errorf("invalid time format: expected format is YYYY-MM-DDTHH:MM:SSZ or YYYY-MM-DDTHH:MM:SS±hh:mm (RFC3339)")
)

// Time is a wrapper around time.Time that supports JSON marshalling/unmarshalling
// with a more concise error message.
type Time struct {
	time.Time
}

// UnmarshalJSON unmarshals a JSON time string into a Time. It provides a more
// concise error message than the default time.Time implementation.
func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || len(s) == 0 {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ErrInvalidTimeFormat
	}
	t.Time = parsed
	return nil
}

// MarshalJSON marshals a Time into a JSON time string.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte{}, nil
	}
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}
