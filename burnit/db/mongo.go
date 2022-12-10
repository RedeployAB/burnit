package db

import (
	"context"
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoClient wraps around mongo.Client to assist with
// implementations of calls related to mongo clients.
type mongoClient struct {
	client     *mongo.Client
	collection *mongo.Collection
	address    string
}

// MongoClientOptions contains options for mongoClient.
type MongoClientOptions struct {
	URI           string
	Address       string
	Database      string
	Collection    string
	Username      string
	Password      string
	SSL           bool
	DirectConnect bool
}

// NewMongoClient creates and returns a *mongoClient.
func NewMongoClient(opts *MongoClientOptions) (*mongoClient, error) {
	if opts == nil {
		opts = &MongoClientOptions{}
	}

	clientOpts := options.Client()
	if len(opts.URI) == 0 && len(opts.Address) != 0 && len(opts.Database) != 0 {
		opts.URI = connectionURI(opts.Address, opts.Database)
		clientOpts = setClientOptions(clientOpts, opts)
	} else {
		// Custom error? Update error message.
		return nil, errors.New("could not create URI for connection")
	}

	client, err := mongo.NewClient(clientOpts.ApplyURI(opts.URI))
	if err != nil {
		return nil, err
	}

	return &mongoClient{
		client:     client,
		collection: client.Database(opts.Database).Collection(opts.Collection),
		address:    clientOpts.Hosts[0],
	}, nil
}

// Connect wraps around mongoClients (mongo.Client)
// Connect method.
func (c *mongoClient) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

// Disconnect wraps around mongoClients (mongo.Client)
// Disconnect method.
func (c *mongoClient) Disconnect(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil && err != mongo.ErrClientDisconnected {
		return err
	}
	return nil
}

// Find implements and calls the method FindOne from
// mongo.Collection.
func (c *mongoClient) Find(ctx context.Context, id string) (*Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	var s Secret
	if err = c.collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&s); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

// Insert implents and calls the the method InsertOne from
// mongo.Collection.
func (c *mongoClient) Insert(ctx context.Context, s *Secret) (*Secret, error) {
	res, err := c.collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}

	return &Secret{
		ID:        res.InsertedID.(primitive.ObjectID).Hex(),
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// Delete implements and calls the method DeleteOne from
// mongo.Collection.
func (c *mongoClient) Delete(ctx context.Context, id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, nil
	}

	res, err := c.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteMany implements and calls the method DeleteMany from
// mongo.Collection.
func (c *mongoClient) DeleteMany(ctx context.Context) (int64, error) {
	bsonQ := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: time.Now()}}}}
	res, err := c.collection.DeleteMany(ctx, bsonQ)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// toURI returns a mongodb connection URI from the provided
// config.Database options.
func toURI(opts *MongoClientOptions) string {
	var b strings.Builder
	b.WriteString("mongodb://")
	if len(opts.Username) > 0 {
		b.WriteString(opts.Username)
		if len(opts.Password) > 0 {
			b.WriteString(":" + opts.Password)
		}
		b.WriteString("@")
	}

	b.WriteString(opts.Address)
	if len(opts.Database) > 0 {
		b.WriteString("/" + opts.Database)
	}

	if opts.SSL {
		b.WriteString("?ssl=true")
	}

	return b.String()
}

// setClientOptions set options on the provided ClientOptions.
func setClientOptions(clientOpts *options.ClientOptions, opts *MongoClientOptions) *options.ClientOptions {
	if len(opts.Username) > 0 {
		credential := options.Credential{Username: opts.Username}
		if len(opts.Password) > 0 {
			credential.Password = opts.Password
			credential.PasswordSet = true
		}
		clientOpts.SetAuth(credential)
	}
	if opts.SSL {
		clientOpts.SetTLSConfig(&tls.Config{})
	}
	if opts.DirectConnect {
		clientOpts.SetDirect(opts.DirectConnect)
	}

	return clientOpts
}

// connectionURI creates a connection URI from the provided
// address and database.
func connectionURI(address, database string) string {
	var b strings.Builder
	b.WriteString("mongodb://")
	b.WriteString(address)
	if len(database) > 0 {
		b.WriteString("/")
		b.WriteString(database)
	}
	return b.String()
}
