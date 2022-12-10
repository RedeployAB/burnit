package db

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// redisClient wraps around redis.Client to assist with
// implementations of calls related to redis clients.
type redisClient struct {
	client *redis.Client
}

// RedisClientOptions contains options for redisClient.
type RedisClientOptions struct {
	URI         string
	Address     string
	Password    string
	Database    string
	SSL         bool
	DialTimeout time.Duration
	Timeout     time.Duration
}

// NewRedisClient creates and returns a *redisClient.
func NewRedisClient(opts *RedisClientOptions) (*redisClient, error) {
	if opts == nil {
		opts = &RedisClientOptions{}
	}

	var clientOpts *redis.Options
	if len(opts.URI) != 0 {
		return &redisClient{
			client: redis.NewClient(fromURI(opts.URI)),
		}, nil
	}

	database, err := strconv.Atoi(opts.Database)
	if err != nil {
		database = 0
	}

	if opts.DialTimeout == 0 {
		opts.DialTimeout = time.Second * 30
	}

	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 5
	}

	clientOpts = &redis.Options{
		Addr:            opts.Address,
		Password:        opts.Password,
		DB:              database,
		DialTimeout:     opts.DialTimeout,
		ReadTimeout:     opts.Timeout,
		MaxRetries:      20,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 5 * time.Second,
	}
	if opts.SSL {
		clientOpts.TLSConfig = &tls.Config{}
	}

	return &redisClient{
		client: redis.NewClient(clientOpts),
	}, nil
}

// Connect connects and tests connection to redis.
func (c *redisClient) Connect(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

// Disconnect closes connection to redis.
func (c *redisClient) Disconnect(ctx context.Context) error {
	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

// Find implements and calls the method Get from
// redis.Client.
func (c *redisClient) Find(ctx context.Context, id string) (*Secret, error) {
	res, err := c.client.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var s Secret
	if err := json.Unmarshal([]byte(res), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

// Insert implents and calls the the method Set from
// redis.Client.
func (c *redisClient) Insert(ctx context.Context, s *Secret) (*Secret, error) {
	s.ID = uuid.New().String()
	secret, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	if err := c.client.Set(ctx, s.ID, secret, time.Until(s.ExpiresAt)).Err(); err != nil {
		return nil, err
	}

	return &Secret{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// Delete implements and calls the method Del from
// redis.Client.
func (c *redisClient) Delete(ctx context.Context, id string) (int64, error) {
	res, err := c.client.Del(ctx, id).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

// DeleteOne exists to fulfill interface Database.
func (c *redisClient) DeleteMany(ctx context.Context) (int64, error) {
	return 0, nil
}

// fromURI creates and returns *redis.Options from provided URI.
func fromURI(uri string) *redis.Options {
	opts := &redis.Options{}

	parts := strings.Split(uri, ",")
	opts.Addr = AddressFromRedisURI(parts[0])
	for i := 1; i < len(parts); i++ {
		subParts := strings.SplitN(parts[i], "=", 2)
		switch strings.ToLower(subParts[0]) {
		case "password":
			opts.Password = subParts[1]
		case "ssl":
			if strings.ToLower(subParts[1]) == "true" {
				opts.TLSConfig = &tls.Config{}
			}
		}
	}
	return opts
}

// AddressFromRedisURI returns the address (<host>:<port>) from
// a redis connection string.
func AddressFromRedisURI(uri string) string {
	reg := regexp.MustCompile("^redis://|^rediss://")
	res := reg.ReplaceAllString(uri, "${1}")
	return strings.Split(res, ",")[0]
}
