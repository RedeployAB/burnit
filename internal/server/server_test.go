package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/internal/log"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			secrets secret.Service
			options []Option
		}
		want *server
	}{
		{
			name: "default",
			input: struct {
				secrets secret.Service
				options []Option
			}{
				secrets: &stubSecretService{},
				options: nil,
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         defaultHost + ":" + defaultPort,
					Handler:      &router{ServeMux: http.NewServeMux()},
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				secrets: &stubSecretService{},
				router:  &router{ServeMux: http.NewServeMux()},
				log:     log.New(),
			},
		},
		{
			name: "with options",
			input: struct {
				secrets secret.Service
				options []Option
			}{
				secrets: &stubSecretService{},
				options: []Option{
					WithOptions(Options{
						Router:       NewRouter(),
						Logger:       log.New(),
						Host:         "localhost",
						Port:         8081,
						ReadTimeout:  10 * time.Second,
						WriteTimeout: 10 * time.Second,
						IdleTimeout:  15 * time.Second,
					}),
				},
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         "localhost:8081",
					Handler:      NewRouter(),
					ReadTimeout:  10 * time.Second,
					WriteTimeout: 10 * time.Second,
					IdleTimeout:  15 * time.Second,
				},
				secrets: &stubSecretService{},
				router:  NewRouter(),
				log:     log.New(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := New(test.input.secrets, test.input.options...)
			if got == nil {
				t.Errorf("New(%v) = nil; want %v", test.input, test.want)
			}

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}, stubSecretService{}), cmpopts.IgnoreUnexported(http.Server{}, http.ServeMux{}), cmpopts.IgnoreFields(server{}, "stopCh", "errCh", "log")); diff != "" {
				t.Errorf("New(%v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			router:  &router{ServeMux: http.NewServeMux()},
			secrets: &stubSecretService{},
			log: &stubLogger{
				logs: &logs,
			},
			stopCh: make(chan os.Signal),
			errCh:  make(chan error),
		}
		go func() {
			time.Sleep(time.Millisecond * 100)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		srv.Start()

		want := []string{
			"Server started.",
			"address",
			"localhost:8080",
			"Server stopped.",
			"reason",
			"interrupt",
		}

		if diff := cmp.Diff(want, logs); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

func TestServer_Start_Error(t *testing.T) {
	t.Run("start server", func(t *testing.T) {
		logs := []string{}
		srv := &server{
			httpServer: &http.Server{
				Addr: "localhost:8080",
			},
			router:  &router{ServeMux: http.NewServeMux()},
			secrets: &stubSecretService{},
			log: &stubLogger{
				logs: &logs,
			},
			stopCh: make(chan os.Signal),
			errCh:  make(chan error),
		}

		httpServer := &http.Server{
			Addr: "localhost:8080",
		}

		go func() {
			go func() {
				time.Sleep(time.Millisecond * 100)
				httpServer.Shutdown(context.Background())
			}()
			httpServer.ListenAndServe()
		}()

		time.Sleep(time.Millisecond * 10)
		gotErr := srv.Start()
		if gotErr == nil {
			t.Errorf("Start() = nil; want error")
		}

		wantErr := errors.New("listen tcp 127.0.0.1:8080: bind: address already in use")
		if diff := cmp.Diff(wantErr.Error(), gotErr.Error()); diff != "" {
			t.Errorf("Start() = unexpected result (-want +got):\n%s\n", diff)
		}
	})
}

type stubLogger struct {
	logs *[]string
}

func (l *stubLogger) Info(msg string, args ...any) {
	if l.logs == nil {
		l.logs = &[]string{}
	}

	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *stubLogger) Error(msg string, args ...any) {
	if l.logs == nil {
		l.logs = &[]string{}
	}

	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *stubLogger) Debug(msg string, args ...any) {}

func (l *stubLogger) Warn(msg string, args ...any) {}
