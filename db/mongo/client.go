package mongo

import (
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Result is the interface for MongoDB results.
type Result interface {
	Decode(v any) error
}

// ClientOptions contains options for the client.
type ClientOptions = options.ClientOptions

// Client is the interface for the MongoDB client. Contains
// methods for interacting with the database and collections.
type Client interface {
	Database(database string) Client
	Collection(collection string) Client
	FindOne(ctx context.Context, filter any) (Result, error)
	InsertOne(ctx context.Context, document any) (string, error)
	DeleteOne(ctx context.Context, filter any) error
	DeleteMany(ctx context.Context, filter any) error
}

// client wraps the MongoDB client.
type client struct {
	cl   *mongo.Client
	db   *mongo.Database
	coll *mongo.Collection
	mu   sync.Mutex
}

// NewClient creates and configures a new client.
func NewClient(ctx context.Context, options ...*ClientOptions) (*client, error) {
	cl, err := mongo.Connect(ctx, options...)
	if err != nil {
		return nil, err
	}
	if err := cl.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &client{
		cl: cl,
	}, nil
}

// Database sets the database and returns the client.
func (c *client) Database(database string) Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.db = c.cl.Database(database)
	return c
}

// Collection sets the collection and returns the client.
func (c *client) Collection(collection string) Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.coll = c.db.Collection(collection)
	return c
}

// Disconnect disconnects the client.
func (c *client) Disconnect(ctx context.Context) error {
	err := c.cl.Disconnect(ctx)
	if err == nil || err == mongo.ErrClientDisconnected {
		return nil
	}
	return err
}

// FindOne finds a document in the collection.
func (c *client) FindOne(ctx context.Context, filter any) (Result, error) {
	res := c.coll.FindOne(ctx, filter)
	return res, res.Err()
}

// InsertOne inserts a document into the collection.
func (c *client) InsertOne(ctx context.Context, document any) (string, error) {
	res, err := c.coll.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}
	return parseID(res.InsertedID)
}

// DeleteOne deletes a document from the collection.
func (c *client) DeleteOne(ctx context.Context, filter any) error {
	res, err := c.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrDocumentNotDeleted
	}
	return nil
}

// DeleteMany deletes documents from the collection.
func (c *client) DeleteMany(ctx context.Context, filter any) error {
	res, err := c.coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrDocumentsNotDeleted
	}
	return nil
}

// parseID parses the ID into a string.
func parseID(id any) (string, error) {
	switch id := id.(type) {
	case primitive.ObjectID:
		return id.Hex(), nil
	case string:
		return id, nil
	}
	return "", errors.New("invalid ID")
}
