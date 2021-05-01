package db

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   = "burnit"
	collection = "secrets"
)

// mongoClient wraps around mongo.Client to assist with
// implementations of calls related to mongo clients.
type mongoClient struct {
	client     *mongo.Client
	collection *mongo.Collection
	address    string
}

// Connect wraps around mongoClients (mongo.Client)
// Connect method.
func (c *mongoClient) Connect(ctx context.Context) error {
	return c.client.Connect(ctx)
}

// Disconnect wraps around mongoClients (mongo.Client)
// Disconnect method.
func (c *mongoClient) Disconnect(ctx context.Context) error {
	return c.client.Disconnect(ctx)
}

// GetAddress returns the address (host) of the client.
func (c *mongoClient) GetAddress() string {
	return c.address
}

// FindOne implements and calls the method FindOne from
// mongo.Collection.
func (c *mongoClient) FindOne(id string) (*Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	bsonQ := bson.D{{Key: "_id", Value: oid}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s Secret
	if err = c.collection.FindOne(ctx, bsonQ).Decode(&s); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

// InsertOne implents and calls the the method InsertOne from
// mongo.Collection.
func (c *mongoClient) InsertOne(s *Secret) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}
	oid := res.InsertedID.(primitive.ObjectID).Hex()

	return &Secret{
		ID:        oid,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// DeleteOne implements and calls the method InsertOne from
// mongo.Collection.
func (c *mongoClient) DeleteOne(id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteMany implements and calls the method DeleteMany from
// mongo.Collection.
func (c *mongoClient) DeleteMany() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bsonQ := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: time.Now()}}}}

	res, err := c.collection.DeleteMany(ctx, bsonQ)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// mongoConnect implements mongo.Clients connection methods,
// helpers and connections checks.
func mongoConnect(opts config.Database) (*mongoClient, error) {
	clientOpts := options.Client()
	uri := opts.URI
	if len(uri) == 0 {
		uri = connectionURI(opts.Address, opts.Database)
		clientOpts = setOptions(clientOpts, opts)
	}

	client, err := mongo.NewClient(clientOpts.ApplyURI(uri))
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

	return &mongoClient{
		client:     client,
		collection: client.Database(database).Collection(collection),
		address:    clientOpts.Hosts[0],
	}, nil
}

// toURI returns a mongodb connection URI from
// the provided config.Database options.
func toURI(opts config.Database) string {
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

func setOptions(clientOpts *options.ClientOptions, opts config.Database) *options.ClientOptions {
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
