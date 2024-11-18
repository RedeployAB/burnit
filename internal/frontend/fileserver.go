package frontend

import (
	"io/fs"
	"net/http"
	_path "path"
	"strings"
)

// gzipFile is a file that has been compressed using gzip.
type gzipFile struct {
	path        string
	contentType string
}

// FileServer creates a file server handler.
func FileServer(fsys fs.FS) http.Handler {
	gzipped := map[string]gzipFile{}
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".gz") {
			p := strings.TrimSuffix(path, ".gz")
			gzipped[p] = gzipFile{
				path:        path,
				contentType: getContentType(p),
			}
		}
		return nil
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			if file, ok := gzipped[r.URL.Path]; ok {
				w.Header().Set("Content-Type", file.contentType)
				w.Header().Set("Content-Encoding", "gzip")
				r.URL.Path = file.path
			}
		}

		http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
	})
}

// getContentType returns the content type for a given path to a file.
func getContentType(path string) string {
	switch _path.Ext(path) {
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	}
	return ""
}
