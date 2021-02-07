package db

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// redisClient wraps around redis.Client to assist with
// implementations of calls related to redis clients.
type redisClient struct {
	client *redis.Client
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

// GetAddress returns the address (host) of the client.
func (c *redisClient) GetAddress() string {
	return c.client.Options().Addr
}

// FindOne implements and calls the method Get from
// redis.Client.
func (c *redisClient) FindOne(id string) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

// InsertOne implents and calls the the method Set from
// redis.Client.
func (c *redisClient) InsertOne(s *Secret) (*Secret, error) {
	s.ID = uuid.New().String()
	secret, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.client.Set(ctx, s.ID, secret, s.ExpiresAt.Sub(time.Now())).Err(); err != nil {
		return nil, err
	}

	return &Secret{
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}, nil
}

// DeleteOne implements and calls the method Del from
// redis.Client.
func (c *redisClient) DeleteOne(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.client.Del(ctx, id).Result()
	if err != nil {
		return 0, err
	}
	return res, nil
}

// DeleteOne exists to fulfill interface Database.
func (c *redisClient) DeleteMany() (int64, error) {
	return 0, nil
}

// newRedisClient creates a new redisClient object.
func newRedisClient(opts config.Database) *redisClient {
	clientOpts := &redis.Options{
		Addr:            opts.Address,
		Password:        opts.Password,
		DB:              0,
		DialTimeout:     15 * time.Second,
		MaxRetries:      20,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 5 * time.Second,
		ReadTimeout:     time.Minute,
	}
	if opts.SSL {
		clientOpts.TLSConfig = &tls.Config{}
	}
	client := redis.NewClient(clientOpts)

	return &redisClient{client: client}
}

// redisConnect implements redis.Client connection methods,
// helpers and connections checks.
func redisConnect(opts config.Database) (*redisClient, error) {
	if len(opts.URI) > 0 {
		opts = fromURI(opts)
	}

	client := newRedisClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func fromURI(opts config.Database) config.Database {
	parts := strings.Split(opts.URI, ",")
	opts.Address = config.AddressFromRedisURI(parts[0])

	for i := 1; i < len(parts); i++ {
		subParts := strings.SplitN(parts[i], "=", 2)
		switch strings.ToLower(subParts[0]) {
		case "password":
			opts.Password = subParts[1]
		case "ssl":
			if strings.ToLower(subParts[1]) == "true" {
				opts.SSL = true
			} else {
				opts.SSL = false
			}
		}
	}

	return opts
}
