package main

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/RedeployAB/burnit/burnit/db"
	"github.com/RedeployAB/burnit/burnit/secret"
)

// SetupSecretService creates and returns a configured secret.Service.
func SetupSecretService(opts *config.Configuration) (secret.Service, error) {
	var client db.Client
	var err error
	switch opts.Database.Driver {
	case config.DatabaseDriverRedis:
		client, err = db.NewRedisClient(newRedisClientOptions(&opts.Database))
	case config.DatabaseDriverMongo:
		client, err = db.NewMongoClient(newMongoClientOptions(&opts.Database))
	default:
		return nil, errors.New("could not setup secret service")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	serviceOpts := &secret.ServiceOptions{
		Secrets:       db.NewSecretRepository(client, &db.SecretRepositoryOptions{}),
		EncryptionKey: opts.Server.Security.Encryption.Key,
	}

	return secret.NewService(serviceOpts), err
}

// newRedisClientOptions creates and returns new options for redis client.
func newRedisClientOptions(opts *config.Database) *db.RedisClientOptions {
	return &db.RedisClientOptions{
		URI:      opts.URI,
		Address:  opts.Address,
		Password: opts.Password,
		Database: opts.Database,
		SSL:      opts.SSL,
	}
}

// newMongoClientOptions creates and returns new options for mongo client.
func newMongoClientOptions(opts *config.Database) *db.MongoClientOptions {
	return &db.MongoClientOptions{
		URI:        opts.URI,
		Address:    opts.Address,
		Database:   opts.Database,
		Collection: "secrets",
		Username:   opts.Username,
		Password:   opts.Password,
		SSL:        opts.SSL,
	}
}
