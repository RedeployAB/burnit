package config

import (
	"errors"
	"fmt"

	"github.com/RedeployAB/burnit/internal/db"
	"github.com/RedeployAB/burnit/internal/db/inmem"
	"github.com/RedeployAB/burnit/internal/db/mongo"
	"github.com/RedeployAB/burnit/internal/db/redis"
	"github.com/RedeployAB/burnit/internal/db/sql"
	"github.com/RedeployAB/burnit/internal/secret"
	"github.com/RedeployAB/burnit/internal/session"
	"github.com/RedeployAB/burnit/internal/ui"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

// services contains the configured services and UI.
type services struct {
	Secrets secret.Service
	UI      ui.UI
}

// Setup configures the services and UI and returns the configured components.
func Setup(config *Configuration) (*services, error) {
	secretSvc, err := setupSecretService(config.Services.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret service: %w", err)
	}

	var ui ui.UI
	if config.Server.BackendOnly == nil || !*config.Server.BackendOnly {
		ui, err = setupUI(config.UI)
		if err != nil {
			return nil, fmt.Errorf("failed to setup frontend services: %w", err)
		}
	}

	return &services{
		Secrets: secretSvc,
		UI:      ui,
	}, nil
}

// setupSecretService sets up the secret service.
func setupSecretService(config Secret) (secret.Service, error) {
	dbClient, err := setupDBClient(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	store, err := setupSecretStore(dbClient, &config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret store: %w", err)
	}

	return secret.NewService(
		store,
		secret.WithTimeout(config.Timeout),
	)
}

// setupSecretStore sets up the secret store.
func setupSecretStore(clients *dbClient, config *Database) (db.SecretStore, error) {
	var store db.SecretStore
	var err error
	switch {
	case clients.mongo != nil:
		store, err = mongo.NewSecretStore(clients.mongo, func(o *mongo.SecretStoreOptions) {
			o.Timeout = config.Timeout
		})
	case clients.sql != nil:
		store, err = sql.NewSecretStore(clients.sql, func(o *sql.SecretStoreOptions) {
			o.Timeout = config.Timeout
		})
	case clients.redis != nil:
		store, err = redis.NewSecretStore(clients.redis)
	default:
		return nil, errors.New("no database client configured")
	}

	if err != nil {
		return nil, err
	}

	return store, nil
}

// setupUI sets up the UI.
func setupUI(config UI) (ui.UI, error) {
	var templatesDir, staticDir string
	var runtimeParse bool

	if config.RuntimeParse != nil && *config.RuntimeParse {
		templatesDir = defaultRuntimeParseTemplateDir
		staticDir = defaultRuntimeParseStaticDir
		runtimeParse = true
	}

	sessionStore, err := setupSessionStore(&config.Services.Session.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup session store: %w", err)
	}

	sessionSvc, err := session.NewService(
		sessionStore,
		session.WithTimeout(config.Services.Session.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup session service: %w", err)
	}

	u, err := ui.New(sessionSvc, func(o *ui.Options) {
		o.RuntimeParse = runtimeParse
		o.TemplateDir = templatesDir
		o.StaticDir = staticDir
	})
	if err != nil {
		return nil, fmt.Errorf("failed to setup UI: %w", err)
	}

	return u, nil
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
