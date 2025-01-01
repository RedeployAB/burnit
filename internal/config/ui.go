package config

import (
	"errors"
	"fmt"

	"github.com/RedeployAB/burnit/internal/db"
	"github.com/RedeployAB/burnit/internal/db/inmem"
	"github.com/RedeployAB/burnit/internal/db/mongo"
	"github.com/RedeployAB/burnit/internal/db/redis"
	"github.com/RedeployAB/burnit/internal/db/sql"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"
)

// SetupUI sets up the UI.
func SetupUI(config UI) (ui.UI, session.Service, error) {
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

	sessionSvc, err := setupSessionService(&config.Services.Session.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	return u, sessionSvc, nil
}

// setupSessionService sets up the session service.
func setupSessionService(config *SessionDatabase) (session.Service, error) {
	client, err := setupDBClient(sessionDatabaseToDatabase(config))
	if err != nil && !errors.Is(err, ErrCouldNotDetermineDatabaseDriver) {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	var store db.SessionStore
	switch {
	case client != nil && client.mongo != nil:
		store, err = mongo.NewSessionStore(client.mongo)
	case client != nil && client.sql != nil:
		store, err = sql.NewSessionStore(client.sql)
	case client != nil && client.redis != nil:
		store, err = redis.NewSessionStore(client.redis)
	default:
		store = inmem.NewSessionStore()
		err = nil
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

// sessionDatabaseToDatabase converts a session database to a database.
func sessionDatabaseToDatabase(db *SessionDatabase) *Database {
	return &Database{
		Driver:                db.Driver,
		URI:                   db.URI,
		Address:               db.Address,
		Database:              db.Database,
		Username:              db.Username,
		Password:              db.Password,
		Timeout:               db.Timeout,
		ConnectTimeout:        db.ConnectTimeout,
		MaxOpenConnections:    db.MaxOpenConnections,
		MaxIdleConnections:    db.MaxIdleConnections,
		MaxConnectionLifetime: db.MaxConnectionLifetime,
		Mongo:                 Mongo(db.Mongo),
		Postgres:              Postgres(db.Postgres),
		MSSQL:                 MSSQL(db.MSSQL),
		SQLite:                SQLite(db.SQLite),
		Redis:                 Redis(db.Redis),
	}
}
