package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mgoopts "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const (
	// defaultConnectTimeout is the default timeout for connecting to the MongoDB client.
	defaultConnectTimeout = 10 * time.Second
)

// Result is the interface for MongoDB results.
type Result interface {
	Decode(v any) error
}

// ClientOptions contains options for the client.
type ClientOptions struct {
	URI                string
	Hosts              []string
	Username           string
	Password           string
	ConnectTimeout     time.Duration
	MaxOpenConnections int
	EnableTLS          bool
	ReplicaSet         string
}

// ClientOption is a function that sets options for the client.
type ClientOption func(o *ClientOptions)

// TxFunc is a function that performs a transaction.
type TxFunc func(ctx context.Context, client Client) (any, error)

// Client is the interface for the MongoDB client. Contains
// methods for interacting with the database and collections.
type Client interface {
	Database(database string) Client
	Collection(collection string) Client
	FindOne(ctx context.Context, filter any) (Result, error)
	InsertOne(ctx context.Context, document any) (string, error)
	UpsertOne(ctx context.Context, filter, update any) (string, error)
	DeleteOne(ctx context.Context, filter any) error
	DeleteMany(ctx context.Context, filter any) error
	WithTransaction(ctx context.Context, fn TxFunc) (any, error)
	WithTransactions(ctx context.Context, fns ...TxFunc) ([]any, error)
	ReplicaSetEnabled() bool
	Disconnect(ctx context.Context) error
}

// client wraps the MongoDB client.
type client struct {
	cl         *mongo.Client
	db         *mongo.Database
	coll       *mongo.Collection
	mu         sync.Mutex
	replicaSet string
}

// NewClient creates and configures a new client.
func NewClient(options ...ClientOption) (*client, error) {
	opts := ClientOptions{
		ConnectTimeout: defaultConnectTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	clientOpts := createClientOptions(&opts)
	cl, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := cl.Ping(ctx, nil); err != nil {
		return nil, err
	}

	var replicaSet string
	if clientOpts.ReplicaSet != nil && len(*clientOpts.ReplicaSet) > 0 {
		replicaSet = *clientOpts.ReplicaSet
		ok, err := replicaSetEnabled(ctx, cl)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New("replica set not enabled")
		}
	}

	return &client{
		cl:         cl,
		replicaSet: replicaSet,
	}, nil
}

// createClientOptions creates a new client options for underlying MongoDB client.
func createClientOptions(options *ClientOptions) *mgoopts.ClientOptions {
	opts := mgoopts.Client()
	if options == nil {
		return opts
	}

	if len(options.URI) > 0 {
		return opts.ApplyURI(options.URI)
	}
	if len(options.Hosts) > 0 {
		opts.Hosts = options.Hosts
	}
	if len(options.Username) > 0 && len(options.Password) > 0 {
		opts.Auth = &mgoopts.Credential{
			Username: options.Username,
			Password: options.Password,
		}
	}
	if options.MaxOpenConnections > 0 {
		opts.SetMaxPoolSize(uint64(options.MaxOpenConnections))
	}
	if options.EnableTLS {
		opts.TLSConfig = &tls.Config{}
	}
	if len(options.ReplicaSet) > 0 {
		opts.SetReplicaSet(options.ReplicaSet)
	}

	return opts
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

// UpsertOne upserts a document into the collection.
func (c *client) UpsertOne(ctx context.Context, filter, update any) (string, error) {
	res, err := c.coll.UpdateOne(ctx, filter, update, mgoopts.Update().SetUpsert(true))
	if err != nil {
		return "", err
	}

	if res.UpsertedCount > 0 {
		return parseID(res.UpsertedID)
	}
	if res.ModifiedCount > 0 {
		filter, ok := filter.(bson.D)
		if !ok {
			return "", errors.New("could not parse ID")
		}
		id, ok := filter[0].Value.(string)
		if !ok {
			return "", errors.New("could not parse ID")
		}
		return id, nil
	}

	return "", ErrDocumentNotUpserted
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

// WithTransaction runs the function as a transaction. Transactions
// can only be run if replica set is configured.
func (c *client) WithTransaction(ctx context.Context, fn TxFunc) (any, error) {
	session, err := c.cl.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	return execTransaction(ctx, c, session, fn)
}

// WithTransactions runs the functions as transactions. Transactions
// can only be run if replica set is configured.
func (c *client) WithTransactions(ctx context.Context, fns ...TxFunc) ([]any, error) {
	session, err := c.cl.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	var results []any
	for _, fn := range fns {
		result, err := execTransaction(ctx, c, session, fn)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// execTransaction runs the provided TxFunc.
func execTransaction(ctx context.Context, client Client, session mongo.Session, fn TxFunc) (any, error) {
	result, err := session.WithTransaction(ctx, func(ctx mongo.SessionContext) (any, error) {
		return fn(ctx, client)
	}, mgoopts.Transaction().
		SetWriteConcern(writeconcern.Majority()).
		SetReadConcern(readconcern.Snapshot()))
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ReplicaSetEnabled returns true if the client is configured
// with a replica set.
func (c *client) ReplicaSetEnabled() bool {
	return len(c.replicaSet) > 0
}

// Disconnect disconnects the client.
func (c *client) Disconnect(ctx context.Context) error {
	err := c.cl.Disconnect(ctx)
	if err == nil || errors.Is(err, mongo.ErrClientDisconnected) {
		return nil
	}
	return err
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

// replicaSetEnabled checks if the server has replica set enabled.
func replicaSetEnabled(ctx context.Context, client *mongo.Client) (bool, error) {
	res := client.Database("admin").RunCommand(ctx, bson.M{"hello": 1})
	if res.Err() != nil {
		return false, res.Err()
	}

	var cmdRes bson.M
	if err := res.Decode(&cmdRes); err != nil {
		return false, err
	}

	if _, ok := cmdRes["setName"]; ok {
		return true, nil
	}
	return false, nil
}

// now is a function that returns the current time.
var now = func() time.Time {
	return time.Now().UTC()
}
