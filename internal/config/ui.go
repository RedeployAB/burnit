package config

import (
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
	var runtimeParse bool

	if config.RuntimeParse != nil && *config.RuntimeParse {
		templatesDir = defaultRuntimeParseTemplateDir
		staticDir = defaultRuntimeParseStaticDir
		runtimeParse = true
	}

	u, err := ui.New(func(o *ui.Options) {
		o.RuntimeParse = runtimeParse
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup UI: %w", err)
	}

	sessionStore, err := setupSessionStore(&config.Services.Session.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	sessionSvc, err := session.NewService(
		sessionStore,
		session.WithTimeout(config.Services.Session.Timeout),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to setup session service: %w", err)
	}

	return u, sessionSvc, nil
}

// setupSessionStore sets up the session store.
func setupSessionStore(config *SessionDatabase) (db.SessionStore, error) {
	client, err := setupDBClient(sessionDatabaseToDatabase(config))
	if err != nil {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	var store db.SessionStore
	switch {
	case client != nil && client.mongo != nil:
		store, err = mongo.NewSessionStore(client.mongo, func(o *mongo.SessionStoreOptions) {
			o.Timeout = config.Timeout
		})
	case client != nil && client.sql != nil:
		store, err = sql.NewSessionStore(client.sql, func(o *sql.SessionStoreOptions) {
			o.Timeout = config.Timeout
		})
	case client != nil && client.redis != nil:
		store, err = redis.NewSessionStore(client.redis)
	default:
		store = inmem.NewSessionStore()
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	return store, nil
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
