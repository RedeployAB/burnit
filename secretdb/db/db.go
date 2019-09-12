package db

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
)

// Connect is used to connect to database with options
// specified in the passed in options argument.
func Connect(opts config.Database) (*mongo.Client, error) {
	uri := opts.URI
	if uri == "" {
		uri = toConnectionURI(opts)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
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
