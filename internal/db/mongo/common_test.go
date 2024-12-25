package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type stubMongoClient struct {
	err     error
	secrets []db.Secret
}

func (c *stubMongoClient) Database(database string) Client {
	return c
}

func (c *stubMongoClient) Collection(collection string) Client {
	return c
}

func (c stubMongoClient) FindOne(ctx context.Context, filter any) (Result, error) {
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
					return stubResult{data: data}, nil
				}
			default:
				return nil, errors.New("invalid filter")
			}
		}
	}

	return nil, ErrNoDocuments
}

func (c *stubMongoClient) InsertOne(ctx context.Context, document any) (string, error) {
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

func (c *stubMongoClient) DeleteOne(ctx context.Context, filter any) error {
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
	return ErrDocumentNotDeleted
}

func (c *stubMongoClient) DeleteMany(ctx context.Context, filter any) error {
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

func (c stubMongoClient) Disconnect(ctx context.Context) error {
	return nil
}

type stubResult struct {
	err  error
	data []byte
}

func (r stubResult) Decode(v any) error {
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
