package db

// mongo.go consists of abstrations and wrappers
// around the mongodb drivers to make the repository
// parts of these libraries testable with the help
// of interfacces located in db.go.

import (
	"context"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect is used to connect to database with options
// specified in the passed in options argument.
func Connect(opts config.Database) (Client, error) {
	client, err := mongoConnect(opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Close disconnects a connection to a database.
func Close(c Connection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := c.Disconnect(ctx); err != nil && err != mongo.ErrClientDisconnected {
		return err
	}
	return nil
}

// mongoClient wraps around mongo.Client to assist with
// implementations of calls related to mongo clients.
type mongoClient struct {
	client *mongo.Client
}

// Database wraps around mongoClients (mongo.Client)
// Database method.
func (c *mongoClient) Database(name string) Database {
	return &mongoDatabase{database: c.client.Database(name)}
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

// mongoDatabase wraps around mongo.Database to assist with
// implementations of calls related to mongo databases.
type mongoDatabase struct {
	database *mongo.Database
}

// Collection wraps around mongoDatabases (mongo.Datbase)
// Collection method.
func (db *mongoDatabase) Collection(name string) Collection {
	return &mongoCollection{collection: db.database.Collection(name)}
}

// mongoCollection wraps around mongo.Collection to assist with
// implementations of calls related to mongo collections.
type mongoCollection struct {
	collection *mongo.Collection
}

// FindOne implements and calls the method FindOne from
// mongo.Collection.
func (c *mongoCollection) FindOne(id string) (*models.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	bsonQ := bson.D{{Key: "_id", Value: oid}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s models.Secret
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
func (c *mongoCollection) InsertOne(s *models.Secret) (*models.Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.collection.InsertOne(ctx, s)
	if err != nil {
		return nil, err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return &models.Secret{
		ID:         oid,
		Passphrase: s.Passphrase,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
	}, nil
}

// DeleteOne implements and calls the method InsertOne from
// mongo.Collection.
func (c *mongoCollection) DeleteOne(id string) (int64, error) {
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
func (c *mongoCollection) DeleteMany() (int64, error) {
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

	return &mongoClient{client: client}, nil
}

// toURI returns a mongodb connection URI from
// the provided config.Database options.
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
