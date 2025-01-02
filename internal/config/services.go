package config

import (
	"errors"
	"fmt"

	"github.com/RedeployAB/burnit/internal/db"
	"github.com/RedeployAB/burnit/internal/db/mongo"
	"github.com/RedeployAB/burnit/internal/db/redis"
	"github.com/RedeployAB/burnit/internal/db/sql"
	"github.com/RedeployAB/burnit/internal/secret"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

// services contains the configured and setup services.
type services struct {
	Secrets secret.Service
}

// SetupServices configures and sets up the services.
func SetupServices(config Services) (*services, error) {
	dbClient, err := setupDBClient(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	store, err := setupSecretStore(dbClient, &config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret store: %w", err)
	}

	secrets, err := secret.NewService(
		store,
		secret.WithValueMaxCharacters(config.Secret.ValueMaxCharacters),
		secret.WithTimeout(config.Secret.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret service: %w", err)
	}

	return &services{
		Secrets: secrets,
	}, nil
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
