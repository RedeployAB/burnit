package db

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/db/mongo"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/mongo-driver/bson"
)

func TestNewSecretRepository(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			client  mongo.Client
			options []MongoSecretRepositoryOption
		}
		want    *MongoSecretRepository
		wantErr error
	}{
		{
			name: "new secret repository",
			input: struct {
				client  mongo.Client
				options []MongoSecretRepositoryOption
			}{
				client: &mockMongoClient{},
			},
			want: &MongoSecretRepository{
				client:     &mockMongoClient{},
				collection: defaultMongoSecretRepositoryCollection,
			},
		},
		{
			name: "new secret repository - with options",
			input: struct {
				client  mongo.Client
				options []MongoSecretRepositoryOption
			}{
				client: &mockMongoClient{},
				options: []MongoSecretRepositoryOption{
					func(o *MongoSecretRepositoryOptions) {
						o.Database = "test"
						o.Collection = "test"
					},
				},
			},
			want: &MongoSecretRepository{
				client:     &mockMongoClient{},
				collection: "test",
			},
		},
		{
			name: "new secret repository - nil client",
			input: struct {
				client  mongo.Client
				options []MongoSecretRepositoryOption
			}{
				client: nil,
			},
			wantErr: ErrNilClient,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewMongoSecretRepository(test.input.client, test.input.options...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(MongoSecretRepository{}, mockMongoClient{})); diff != "" {
				t.Errorf("NewMongoSecretRepository() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("NewMongoSecretRepository() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestMongoSecretRepository_Get(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			id      string
			secrets []Secret
			err     error
		}
		want    Secret
		wantErr error
	}{
		{
			name: "get secret",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id: "1",
				secrets: []Secret{
					{
						ID:    "1",
						Value: "secret",
					},
				},
			},
			want: Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "get secret - not found",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id: "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "get secret - error",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id:  "1",
				err: errFindOne,
			},
			wantErr: errFindOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := MongoSecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			got, gotErr := repo.Get(context.Background(), test.input.id)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Get() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Get() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestMongoSecretRepository_Create(t *testing.T) {
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secret  Secret
			secrets []Secret
			err     error
		}
		want    Secret
		wantErr error
	}{
		{
			name: "create secret",
			input: struct {
				secret  Secret
				secrets []Secret
				err     error
			}{
				secret: Secret{
					ID:    "1",
					Value: "secret",
				},
			},
			want: Secret{
				ID:    "1",
				Value: "secret",
			},
		},
		{
			name: "create secret - error",
			input: struct {
				secret  Secret
				secrets []Secret
				err     error
			}{
				secret: Secret{
					ID:    "1",
					Value: "secret",
				},
				err: errInsertOne,
			},
			wantErr: errInsertOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := MongoSecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			got, gotErr := repo.Create(context.Background(), test.input.secret)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Create() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Create() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestMongoSecretRepository_Delete(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			id      string
			secrets []Secret
			err     error
		}
		wantErr error
	}{
		{
			name: "delete secret",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id: "1",
				secrets: []Secret{
					{
						ID:    "1",
						Value: "secret",
					},
				},
			},
		},
		{
			name: "delete secret - not found",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id: "1",
			},
			wantErr: ErrSecretNotFound,
		},
		{
			name: "delete secret - error",
			input: struct {
				id      string
				secrets []Secret
				err     error
			}{
				id:  "1",
				err: errDeleteOne,
			},
			wantErr: errDeleteOne,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := MongoSecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := repo.Delete(context.Background(), test.input.id)

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Delete() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestMongoSecretRepository_DeleteExpired(t *testing.T) {
	date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	now = func() time.Time {
		return date
	}

	var tests = []struct {
		name  string
		input struct {
			secrets []Secret
			err     error
		}
		wantErr error
	}{
		{
			name: "delete expired secrets",
			input: struct {
				secrets []Secret
				err     error
			}{
				secrets: []Secret{
					{
						ID:        "1",
						Value:     "secret",
						ExpiresAt: date.Add(-time.Hour * 2),
					},
					{
						ID:        "2",
						Value:     "secret",
						ExpiresAt: date.Add(-time.Hour * 2),
					},
				},
			},
		},
		{
			name: "delete expired secrets - error",
			input: struct {
				secrets []Secret
				err     error
			}{
				err: errDeleteMany,
			},
			wantErr: errDeleteMany,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := MongoSecretRepository{
				client: &mockMongoClient{
					secrets: test.input.secrets,
					err:     test.input.err,
				},
			}

			gotErr := repo.DeleteExpired(context.Background())

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DeleteExpired() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type mockMongoClient struct {
	err     error
	secrets []Secret
}

func (c *mockMongoClient) Database(database string) mongo.Client {
	return c
}

func (c *mockMongoClient) Collection(collection string) mongo.Client {
	return c
}

func (c mockMongoClient) FindOne(ctx context.Context, filter any) (mongo.Result, error) {
	if c.err != nil {
		return nil, c.err
	}
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
	return nil, mongo.ErrNoDocuments
}

func (c *mockMongoClient) InsertOne(ctx context.Context, document any) (string, error) {
	if c.err != nil {
		return "", c.err
	}

	secret, ok := document.(Secret)
	if !ok {
		return "", errors.New("invalid document")
	}
	c.secrets = append(c.secrets, secret)
	return secret.ID, nil
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
	return mongo.ErrNoDocuments
}

func (c *mockMongoClient) DeleteMany(ctx context.Context, filter any) error {
	if c.err != nil {
		return c.err
	}

	secretsToKeep := []Secret{}
	for _, secret := range c.secrets {
		if !secret.ExpiresAt.Before(time.Now()) {
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
