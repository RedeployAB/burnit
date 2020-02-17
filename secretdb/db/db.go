package db

import (
	"context"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/secretdb/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection represents a connection to a database.
type Connection interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
}

// Client is a wrapper around mongo.Client.
type Client struct {
	*mongo.Client
}

// Collection is a wrapper around mongo.Collection.
type Collection struct {
	*mongo.Collection
}

// Connect is used to connect to database with options
// specified in the passed in options argument.
func Connect(opts config.Database) (*Client, error) {
	uri := opts.URI
	if len(uri) == 0 {
		uri = toURI(opts)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

// Close disconnects the connection to the database.
func Close(c Connection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := c.Disconnect(ctx); err != nil && err != mongo.ErrClientDisconnected {
		return err
	}
	return nil
}

func toURI(opts config.Database) string {
	var b strings.Builder

	b.WriteString("mongodb://")
	if opts.Username != "" {
		b.WriteString(opts.Username)
		if opts.Password != "" {
			b.WriteString(":" + opts.Password)
		}
		b.WriteString("@")
	}
	b.WriteString(opts.Address)
	if opts.SSL != false {
		b.WriteString("/?ssl=true")
	}

	return b.String()
}
