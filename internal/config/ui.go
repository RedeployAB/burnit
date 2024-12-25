package config

import (
	"fmt"

	"github.com/RedeployAB/burnit/internal/db"
	"github.com/RedeployAB/burnit/internal/db/inmem"
	"github.com/RedeployAB/burnit/internal/db/sql"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
)

// SetupUI sets up the UI.
func SetupUI(config UI, databaseConfig Database) (ui.UI, session.Service, error) {
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

	sessionSvc, err := setupSessionService(&databaseConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	return u, sessionSvc, nil
}

// setupSessionService sets up the session service.
func setupSessionService(config *Database) (session.Service, error) {
	client, err := setupDBClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	var store db.SessionStore
	switch {
	case client.sql != nil:
		store, err = sql.NewSessionStore(client.sql)
	default:
		store = inmem.NewSessionStore()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	svc, err := session.NewService(store)
	if err != nil {
		return nil, fmt.Errorf("failed to setup session service: %w", err)
	}

	return svc, nil
}
