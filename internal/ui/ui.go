package ui

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	_path "path"
	"path/filepath"
	"strings"
)

var (
	//go:embed templates/*
	templateFS embed.FS

	//go:embed static/js/main.min.js* static/css/main.min.css* static/icons/* static/images/*
	staticFS embed.FS
)

const (
	// defaultTemplateDir is the default directory for templates.
	defaultTemplateDir = "templates"
	// defaultInternalTemplateDir is the default directory for internal templates.
	defaultInternalTemplateDir = "internal/ui/templates"
	// defaultStaticDir is the default directory for static files.
	defaultStaticDir = "static"
)

// UI is an interface for rendering templates.
type UI interface {
	Render(w http.ResponseWriter, statusCode int, tmpl string, data any, options ...RenderOption)
	Static() fs.FS
	RuntimeRender() bool
}

// ui is a user interface handler.
type ui struct {
	templates     map[string]*template.Template
	templateDir   string
	staticFS      fs.FS
	runtimeRender bool
}

// Options for the UI.
type Options struct {
	TemplateDir   string
	StaticDir     string
	RuntimeRender bool
}

// Option is a function that configures the UI.
type Option func(o *Options)

// New returns a new UI.
func New(options ...Option) (*ui, error) {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	ui := &ui{
		templates:     make(map[string]*template.Template),
		templateDir:   opts.TemplateDir,
		runtimeRender: opts.RuntimeRender,
	}

	if err := ui.parseTemplates(templateFS, defaultTemplateDir, true); err != nil {
		return nil, err
	}

	if len(ui.templateDir) > 0 {
		if err := ui.parseTemplates(os.DirFS(ui.templateDir), ui.templateDir, false); err != nil {
			return nil, err
		}
	}

	if len(opts.StaticDir) == 0 {
		fsys, err := fs.Sub(staticFS, defaultStaticDir)
		if err != nil {
			return nil, err
		}
		ui.staticFS = fsys
	} else {
		ui.staticFS = os.DirFS(opts.StaticDir)
	}

	return ui, nil
}

// RenderOptions is a configuration for rendering a template.
type RenderOptions struct {
	Partial bool
}

// RenderOption is a function that configures the rendering of a template.
type RenderOption func(o *RenderOptions)

// WithPartial returns a RenderOption that renders a partial template.
func WithPartial() RenderOption {
	return func(o *RenderOptions) {
		o.Partial = true
	}
}

// Render a template with the given data.
func (u ui) Render(w http.ResponseWriter, statusCode int, tmpl string, data any, options ...RenderOption) {
	opts := RenderOptions{}
	for _, option := range options {
		option(&opts)
	}

	execTemplate := "base.html"
	if opts.Partial {
		execTemplate = tmpl
	}

	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", "text/html")
	}

	if u.runtimeRender {
		dir := u.templateDir
		if len(dir) == 0 {
			dir = defaultInternalTemplateDir
		}

		templates := []string{dir + "/" + tmpl + ".html"}
		if partial, err := os.Stat(dir + "/" + "partial-" + tmpl + ".html"); err == nil {
			templates = append(templates, dir+"/"+partial.Name())
		}

		if !opts.Partial {
			templates = append([]string{dir + "/base.html"}, templates...)
		}

		t, err := template.ParseFiles(templates...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		if err := t.ExecuteTemplate(w, execTemplate, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(statusCode)
	if err := u.templates[tmpl].ExecuteTemplate(w, execTemplate, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Static returns the static file system.
func (u ui) Static() fs.FS {
	return u.staticFS
}

// RuntimeRender returns true if the UI should render templates at runtime.
func (u ui) RuntimeRender() bool {
	return u.runtimeRender
}

// trimExtension trims the extension from a file name.
func trimExtension(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// parseTemplates parses and adds templates to the UI.
func (u *ui) parseTemplates(fsys fs.FS, path string, embedded bool) error {
	var prefix string
	if embedded {
		prefix = path + "/"
	} else {
		path = "."
	}

	templates := map[string][]string{}
	err := fs.WalkDir(fsys, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.Contains(d.Name(), ".html") || d.Name() == "base.html" {
			return nil
		}

		tmplName := _path.Base(trimExtension(d.Name()))
		if _, ok := templates[tmplName]; !ok {
			templates[tmplName] = []string{prefix + "base.html"}
		}
		templates[tmplName] = append(templates[tmplName], path)
		return nil
	})
	if err != nil {
		return err
	}

	for tmpl, files := range templates {
		t, err := template.ParseFS(fsys, files...)
		if err != nil {
			return err
		}
		u.templates[tmpl] = t
	}
	return nil
}
