package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/internal/models"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

// redisClient wraps around redis.Client to assist with
// implementations of calls related to redis clients.
type redisClient struct {
	client *redis.Client
}

// Connect connects and tests connection to redis.
func (c *redisClient) Connect(ctx context.Context) error {
	// Add retry logic.
	_, err := c.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

// Disconnect closes connection to redis.
func (c *redisClient) Disconnect(ctx context.Context) error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return nil
}

// FindOne implements and calls the method Get from
// redis.Client.
func (c *redisClient) FindOne(id string) (*models.Secret, error) {
	res, err := c.client.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var s models.Secret
	if err := json.Unmarshal([]byte(res), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

// InsertOne implents and calls the the method Set from
// redis.Client.
func (c *redisClient) InsertOne(s *models.Secret) (*models.Secret, error) {
	s.ID = uuid.New().String()
	secret, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	if err := c.client.Set(s.ID, secret, s.ExpiresAt.Sub(time.Now())).Err(); err != nil {
		return nil, err
	}

	return &models.Secret{
		ID:         s.ID,
		Passphrase: s.Passphrase,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
	}, nil
}

// DeleteOne implements and calls the method Del from
// redis.Client.
func (c *redisClient) DeleteOne(id string) (int64, error) {
	res, err := c.client.Del(id).Result()
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
	client := redis.NewClient(&redis.Options{
		Addr:     opts.URI,
		Password: opts.Password,
		DB:       0,
	})

	return &redisClient{client: client}
}

// redisConnect implements redis.Client connection methods,
// helpers and connections checks.
func redisConnect(opts config.Database) (*redisClient, error) {
	client := newRedisClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
