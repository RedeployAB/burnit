package config

import (
	"fmt"

	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
)

// SetupUI sets up the UI.
func SetupUI(config UI) (ui.UI, session.Store, error) {
	var templatesDir, staticDir string
	var runtimeRender bool

	if config.RuntimeRender != nil && *config.RuntimeRender {
		templatesDir = defaultUIRuntimeRenderTemplateDir
		staticDir = defaultUIRuntimeRenderStaticDir
		runtimeRender = true
	}

	u, err := ui.New(func(o *ui.Options) {
		o.RuntimeRender = runtimeRender
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup UI: %w", err)
	}

	return u, session.NewInMemoryStore(), nil
}
