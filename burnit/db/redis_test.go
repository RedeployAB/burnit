package db

import (
	"crypto/tls"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
)

func TestFromURI(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  *redis.Options
	}{
		{
			name:  "Address only",
			input: "localhost:6379",
			want: &redis.Options{
				Addr: "localhost:6379",
			},
		},
		{
			name:  "Address with protocol",
			input: "redis://localhost:6379",
			want: &redis.Options{
				Addr: "localhost:6379",
			},
		},
		{
			name:  "Address with protocol (secure)",
			input: "rediss://localhost:6379",
			want: &redis.Options{
				Addr: "localhost:6379",
			},
		},
		{
			name:  "Address with password and SSL",
			input: "localhost:6379,password=1234,ssl=true",
			want: &redis.Options{
				Addr:      "localhost:6379",
				Password:  "1234",
				TLSConfig: &tls.Config{},
			},
		},
		{
			name:  "Address with protocol, password and SSL",
			input: "redis://localhost:6379,password=1234,ssl=true",
			want: &redis.Options{
				Addr:      "localhost:6379",
				Password:  "1234",
				TLSConfig: &tls.Config{},
			},
		},
		{
			name:  "Address with protocol (secure), password and SSL",
			input: "rediss://localhost:6379,password=1234,ssl=true",
			want: &redis.Options{
				Addr:      "localhost:6379",
				Password:  "1234",
				TLSConfig: &tls.Config{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := fromURI(test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(redis.Options{}, tls.Config{}, sync.Mutex{}, sync.RWMutex{})); diff != "" {
				t.Errorf("fromURI(%q) = unexpected result (-want, +got)\n%s\n", test.input, diff)
			}
		})
	}
}
