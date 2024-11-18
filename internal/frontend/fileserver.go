package frontend

import (
	"io/fs"
	"net/http"
)

// FileServerOptions are the options for the file server handler.
type FileServerOptions struct {
	StripPrefix string
}

// FileServerOption is a function that sets a file server option.
type FileServerOption func(o *FileServerOptions)

// FileServer creates a file server handler.
func FileServer(fsys fs.FS, options ...FileServerOption) http.Handler {
	opts := FileServerOptions{}
	for _, option := range options {
		option(&opts)
	}

	fserver := http.FileServer(http.FS(fsys))
	if len(opts.StripPrefix) > 0 {
		fserver = http.StripPrefix(opts.StripPrefix, fserver)
	}

	return fserver
}

// WithFileServerStripPrefix returns a FileServerOption that sets the strip prefix for the file server.
func WithFileServerStripPrefix(prefix string) FileServerOption {
	return func(o *FileServerOptions) {
		o.StripPrefix = prefix
	}
}
