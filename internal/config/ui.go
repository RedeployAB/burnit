package config

import (
	"github.com/RedeployAB/burnit/internal/ui"
)

// SetupUI sets up the UI.
func SetupUI(config UI) (ui.UI, error) {
	var templatesDir, staticDir string
	var runtimeRender bool

	if config.RuntimeRender != nil && *config.RuntimeRender {
		templatesDir = defaultUIRuntimeRenderTemplateDir
		staticDir = defaultUIRuntimeRenderStaticDir
		runtimeRender = true
	}

	return ui.New(func(o *ui.Options) {
		o.RuntimeRender = runtimeRender
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
}
