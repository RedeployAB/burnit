package config

import "github.com/RedeployAB/burnit/internal/frontend"

// SetupUI sets up the frontend UI.
func SetupUI(config Frontend) (frontend.UI, error) {
	var templatesDir, staticDir string
	if config.RuntimeRender {
		templatesDir = defaultFrontendRuntimeRenderTemplateDir
		staticDir = defaultFrontendRuntimeRenderStaticDir
	}

	return frontend.NewUI(func(o *frontend.UIOptions) {
		o.RuntimeRender = config.RuntimeRender
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
}
