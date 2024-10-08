package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/db"
	"go.mongodb.org/mongo-driver/bson"
)

type mockMongoClient struct {
	err     error
	secrets []db.Secret
}

func (c *mockMongoClient) Database(database string) Client {
	return c
}

func (c *mockMongoClient) Collection(collection string) Client {
	return c
}

func (c mockMongoClient) FindOne(ctx context.Context, filter any) (Result, error) {
	if c.err != nil {
		return nil, c.err
	}
	if c.secrets != nil {
		for _, secret := range c.secrets {
			switch f := filter.(type) {
			case bson.D:
				if f[0].Key == "id" || f[0].Key == "_id" && f[0].Value == secret.ID {
					data, err := json.Marshal(secret)
					if err != nil {
						return nil, err
					}
					return mockResult{data: data}, nil
				}
			default:
				return nil, errors.New("invalid filter")
			}
		}
	}

	return nil, ErrNoDocuments
}

func (c *mockMongoClient) InsertOne(ctx context.Context, document any) (string, error) {
	if c.err != nil {
		return "", c.err
	}

	if c.secrets != nil {
		secret, ok := document.(db.Secret)
		if !ok {
			return "", errors.New("invalid document")
		}
		c.secrets = append(c.secrets, secret)
		return secret.ID, nil
	}

	return "", errors.New("could not determine database type")
}

func (c *mockMongoClient) UpdateOne(ctx context.Context, filter, update any) error {
	if c.err != nil {
		return c.err
	}
	return nil
}

func (c *mockMongoClient) DeleteOne(ctx context.Context, filter any) error {
	if c.err != nil {
		return c.err
	}
	for i, secret := range c.secrets {
		switch f := filter.(type) {
		case bson.D:
			if f[0].Key == "id" || f[0].Key == "_id" && f[0].Value == secret.ID {
				c.secrets = append(c.secrets[:i], c.secrets[i+1:]...)
				return nil
			}
		default:
			return errors.New("invalid filter")
		}
	}
	return ErrNoDocuments
}

func (c *mockMongoClient) DeleteMany(ctx context.Context, filter any) error {
	if c.err != nil {
		return c.err
	}

	secretsToKeep := []db.Secret{}
	for _, secret := range c.secrets {
		if !secret.ExpiresAt.Before(time.Now().UTC()) {
			secretsToKeep = append(secretsToKeep, secret)
		}
	}
	c.secrets = secretsToKeep
	return nil
}

func (c mockMongoClient) Disconnect(ctx context.Context) error {
	return nil
}

type mockResult struct {
	err  error
	data []byte
}

func (r mockResult) Decode(v any) error {
	if r.err != nil {
		return r.err
	}
	if err := json.Unmarshal(r.data, v); err != nil {
		return err
	}
	return nil
}

var (
	errFindOne    = errors.New("find one error")
	errInsertOne  = errors.New("insert one error")
	errDeleteOne  = errors.New("delete one error")
	errDeleteMany = errors.New("delete many error")
)
