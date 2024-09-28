package config

import (
	"fmt"

	"github.com/RedeployAB/burnit/db"
	"github.com/RedeployAB/burnit/db/mongo"
	"github.com/RedeployAB/burnit/secret"
)

// services contains the configured and setup services.
type services struct {
	Secrets secret.Service
}

// SetupServices configures and sets up the services.
func SetupServices(config Services) (*services, error) {
	repo, err := setupSecretRepository(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret repository: %w", err)
	}

	secrets, err := secret.NewService(
		repo,
		secret.WithEncryptionKey(config.Secret.EncryptionKey),
		secret.WithTimeout(config.Secret.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret service: %w", err)
	}

	return &services{
		Secrets: secrets,
	}, nil
}

// setupSecretRepository sets up the secret repository.
func setupSecretRepository(config *Database) (db.SecretRepository, error) {
	var enableTLS bool
	if config.EnableTLS != nil {
		enableTLS = *config.EnableTLS
	}

	var repo db.SecretRepository
	switch databaseDriver(config) {
	case databaseDriverMongo:
		client, err := mongo.NewClient(func(o *mongo.ClientOptions) {
			o.ConnectionString = config.ConnectionString
			o.Hosts = []string{config.Address}
			o.Username = config.Username
			o.Password = config.Password
			o.ConnectTimeout = config.ConnectTimeout
			o.EnableTLS = enableTLS
		})
		if err != nil {
			return nil, err
		}

		repo, err = mongo.NewSecretRepository(client, func(o *mongo.SecretRepositoryOptions) {
			o.Database = config.Database
			o.Timeout = config.Timeout
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported database driver")
	}

	return repo, nil
}

// databaseDriver returns the database driver.
func databaseDriver(_ *Database) string {
	return databaseDriverMongo
}

var (
	databaseDriverMongo = "mongo"
)
