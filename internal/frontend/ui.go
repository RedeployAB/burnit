package frontend

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed templates/*
	templateFS embed.FS

	//go:embed static/*
	staticFS embed.FS
)

const (
	// defaultTemplateDir is the default directory for templates.
	defaultTemplateDir = "templates"
	// defaultInternalTemplateDir is the default directory for internal templates.
	defaultInternalTemplateDir = "internal/frontend/templates"
	// defaultStaticDir is the default directory for static files.
	defaultStaticDir = "static"
)

// UI is an interface for rendering templates.
type UI interface {
	Render(w http.ResponseWriter, statusCode int, tmpl string, data any)
	Static() fs.FS
}

// ui is a user interface handler.
type ui struct {
	templates     map[string]*template.Template
	templateDir   string
	staticFS      fs.FS
	runtimeRender bool
}

// UIOptions is a configuration for the UI.
type UIOptions struct {
	TemplateDir   string
	StaticDir     string
	RuntimeRender bool
}

// UIOption is a function that configures the UI.
type UIOption func(o *UIOptions)

// New returns a new UI.
func NewUI(options ...UIOption) (*ui, error) {
	opts := UIOptions{}
	for _, option := range options {
		option(&opts)
	}

	ui := &ui{
		templates:     make(map[string]*template.Template),
		templateDir:   opts.TemplateDir,
		runtimeRender: opts.RuntimeRender,
	}

	if err := ui.addTemplates(templateFS, defaultTemplateDir, true); err != nil {
		return nil, err
	}

	if len(ui.templateDir) > 0 {
		if err := ui.addTemplates(dirFS{os.DirFS(ui.templateDir)}, ui.templateDir, false); err != nil {
			return nil, err
		}
	}

	if len(opts.StaticDir) == 0 {
		fs, err := fs.Sub(staticFS, defaultStaticDir)
		if err != nil {
			return nil, err
		}
		ui.staticFS = fs
	} else {
		ui.staticFS = os.DirFS(opts.StaticDir)
	}

	return ui, nil
}

// Render a template with the given data.
func (u ui) Render(w http.ResponseWriter, statusCode int, tmpl string, data any) {
	if u.runtimeRender {
		dir := u.templateDir
		if len(dir) == 0 {
			dir = defaultInternalTemplateDir
		}
		t, err := template.ParseFiles(dir+"/base.html", dir+"/"+tmpl+".html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err := u.templates[tmpl].ExecuteTemplate(w, "base.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Static returns the static file system.
func (u ui) Static() fs.FS {
	return u.staticFS
}

// trimExtension trims the extension from a file name.
func trimExtension(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// addTemplates adds a template to the UI.
func (u *ui) addTemplates(fs fs.ReadDirFS, path string, embedded bool) error {
	files, err := fs.ReadDir(path)
	if err != nil {
		return err
	}

	var prefix string
	if embedded {
		prefix = path + "/"
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), ".html") && file.Name() != "base.html" {
			tmpl, err := template.ParseFS(fs, prefix+"base.html", prefix+file.Name())
			if err != nil {
				return err
			}
			u.templates[trimExtension(file.Name())] = tmpl
		}
	}
	return nil
}

// dirFS is a wraper for fitting os.DirFS into fs.ReadDirFS.
type dirFS struct {
	fs.FS
}

// ReadDir reads a directory.
func (d dirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}
