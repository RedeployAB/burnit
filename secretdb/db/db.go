package db

import (
	"context"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/secretdb/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection is a wrapper around *mongo.Client.
type Connection struct {
	*mongo.Client
}

// Connector is an interface that represents
// Connect, Disconnect and a database connection.
type Connector interface {
	Connect(context.Context) error
	Disconnect(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

// Connect is used to connect to database with options
// specified in the passed in options argument.
func Connect(opts config.Database) (Connector, error) {
	uri := opts.URI
	if uri == "" {
		uri = toConnectionURI(opts)
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

	return &Connection{client}, nil
}

// Close disconnects the connection to the database.
func Close(c Connector) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := c.Disconnect(ctx); err != nil && err != mongo.ErrClientDisconnected {
		return err
	}
	return nil
}

func toConnectionURI(opts config.Database) string {
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
