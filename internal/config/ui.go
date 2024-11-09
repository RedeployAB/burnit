package config

import (
	"github.com/RedeployAB/burnit/internal/frontend"
)

// SetupUI sets up the frontend UI.
func SetupUI(config Frontend) (frontend.UI, error) {
	var templatesDir, staticDir string
	var runtimeRender bool

	if config.RuntimeRender != nil && *config.RuntimeRender {
		templatesDir = defaultFrontendRuntimeRenderTemplateDir
		staticDir = defaultFrontendRuntimeRenderStaticDir
		runtimeRender = true
	}

	return frontend.NewUI(func(o *frontend.UIOptions) {
		o.RuntimeRender = runtimeRender
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
}
