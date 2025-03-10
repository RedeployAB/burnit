package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// defaultDatabase is the default database for the Redis client.
	defaultDatabase = 0
	// defaultConnectTimeout is the default timeout for connecting to the Redis client.
	defaultConnectTimeout = 10 * time.Second
)

// Client is the interface for the Redis client. Contains methods
// for interacting with the database.
type Client interface {
	Get(ctx context.Context, key string) ([]byte, error)
	HGet(ctx context.Context, key string) (map[string]string, error)
	Set(ctx context.Context, key string, value []byte, exp time.Duration) error
	HSet(ctx context.Context, key string, value map[string]any) error
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, exp time.Duration) error
	WithTransaction(ctx context.Context, fn TxFunc) (TxResult, error)
	WithTransactions(ctx context.Context, fns ...TxFunc) (TxResult, error)
	Close() error
}

// client wraps the Redis client.
type client struct {
	rdb *redis.Client
}

// ClientOptions contains options for the client.
type ClientOptions struct {
	URI                   string
	Address               string
	Database              int
	Username              string
	Password              string
	ConnectTimeout        time.Duration
	DialTimeout           time.Duration
	MaxRetries            int
	MinRetryBackoff       time.Duration
	MaxRetryBackoff       time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
	EnableTLS             bool
}

// ClientOption is a function that sets options for the client.
type ClientOption func(o *ClientOptions)

// NewClient creates and configures a new client.
func NewClient(options ...ClientOption) (*client, error) {
	opts := ClientOptions{
		Database:       defaultDatabase,
		ConnectTimeout: defaultConnectTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	rdbOpts, err := createClientOptions(&opts)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	rdb := redis.NewClient(rdbOpts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &client{
		rdb: rdb,
	}, nil
}

// createClientOptions creates a new client options for the underlying Redis client.
func createClientOptions(options *ClientOptions) (*redis.Options, error) {
	opts := &redis.Options{}
	if options == nil {
		return opts, nil
	}

	if len(options.URI) > 0 {
		var err error
		opts, err = redis.ParseURL(options.URI)
		if err != nil {
			return nil, err
		}
		return opts, nil
	}
	if options.Database > 0 {
		opts.DB = options.Database
	}
	if len(options.Address) > 0 {
		opts.Addr = options.Address
	}
	if len(options.Username) > 0 && len(options.Password) > 0 {
		opts.Username = options.Username
		opts.Password = options.Password
	}
	if options.DialTimeout > 0 {
		opts.DialTimeout = options.DialTimeout
	}
	if options.MaxRetries > 0 {
		opts.MaxRetries = options.MaxRetries
	}
	if options.MinRetryBackoff > 0 {
		opts.MinRetryBackoff = options.MinRetryBackoff
	}
	if options.MaxRetryBackoff > 0 {
		opts.MaxRetryBackoff = options.MaxRetryBackoff
	}
	if options.MaxOpenConnections > 0 {
		opts.MaxActiveConns = options.MaxOpenConnections
	}
	if options.MaxIdleConnections > 0 {
		opts.MaxIdleConns = options.MaxIdleConnections
	}
	if options.MaxConnectionLifetime > 0 {
		opts.ConnMaxLifetime = options.MaxConnectionLifetime
	}
	if options.EnableTLS {
		opts.TLSConfig = &tls.Config{}
	}

	return opts, nil
}

// Get returns the value for the key.
func (c client) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	return b, nil
}

// HGet returns the structured data for the key as a map of strings.
func (c client) HGet(ctx context.Context, key string) (map[string]string, error) {
	res, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrKeyNotFound
	}
	return res, nil
}

// Set the value for the key with an expiration time. If the expiration
// time is zero, the key will not expire.
func (c client) Set(ctx context.Context, key string, value []byte, exp time.Duration) error {
	return c.rdb.Set(ctx, key, value, exp).Err()
}

// HSet sets the value (structured data).
func (c client) HSet(ctx context.Context, key string, value map[string]any) error {
	return c.rdb.HSet(ctx, key, value).Err()
}

// Delete the key.
func (c client) Delete(ctx context.Context, key string) error {
	res, err := c.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrKeyNotFound
	}
	return nil
}

// Expire sets expire time for the key.
func (c client) Expire(ctx context.Context, key string, exp time.Duration) error {
	return c.rdb.Expire(ctx, key, exp).Err()
}

// WithTransaction runs the function as a transaction.
func (c *client) WithTransaction(ctx context.Context, fn TxFunc) (TxResult, error) {
	pipe := c.rdb.TxPipeline()
	tx := &tx{pipe: pipe}

	fn(tx)

	return execCommands(ctx, tx)
}

// WithTransactions runs the functions as transactions.
func (c *client) WithTransactions(ctx context.Context, fns ...TxFunc) (TxResult, error) {
	pipe := c.rdb.TxPipeline()
	tx := &tx{pipe: pipe}

	for _, fn := range fns {
		fn(tx)
	}

	return execCommands(ctx, tx)
}

// Close the client and its underlying connections.
func (c client) Close() error {
	err := c.rdb.Close()
	if err == nil || errors.Is(err, redis.ErrClosed) {
		return nil
	}
	return err
}
