package db

// db.go provides interfaces for collection based
// databases.

import (
	"context"
	"time"

	"github.com/RedeployAB/burnit/burnit/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// Client provides methods Connect, Disconnect and Database.
type Client interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	GetAddress() string
}

// Database represents a connection/client/collection
// in context of a database.
type Database interface {
	FindOne(id string) (*Secret, error)
	InsertOne(s *Secret) (*Secret, error)
	DeleteOne(id string) (int64, error)
	DeleteMany() (int64, error)
}

// Connect is used to connect to database with options
// specified in the passed in options argument.
func Connect(opts config.Database) (Client, error) {
	var client Client
	var err error

	switch opts.Driver {
	case "mongo":
		client, err = mongoConnect(opts)
	case "redis":
		client, err = redisConnect(opts)
	}

	if err != nil {
		return nil, err
	}
	return client, nil
}

// Close disconnects a connection to a database.
func Close(c Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := c.Disconnect(ctx); err != nil && err != mongo.ErrClientDisconnected {
		return err
	}
	return nil
}
