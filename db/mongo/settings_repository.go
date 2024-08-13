package mongo

import (
	"context"
	"errors"
	"time"

	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// defaultSettingsRepositoryDatabase is the default database for the settings repository.
	defaultSettingsRepositoryDatabase = "burnit"
	// defaultSettingsRepositoryCollection is the default collection for the settings repository.
	defaultSettingsRepositoryCollection = "settings"
	// defaultSettingsRepositoryTimeout is the default timeout for the settings repository.
	defaultSettingsRepositoryTimeout = 10 * time.Second
)

// SettingsRepository is a repository for the settings.
type SettingsRepository struct {
	client     Client
	collection string
	timeout    time.Duration
}

// SettingsRepositoryOptions is the options for the SettingsRepository.
type SettingsRepositoryOptions struct {
	Database   string
	Collection string
	Timeout    time.Duration
}

// SettingsRepositoryOption is a function that sets options for the SettingsRepository.
type SettingsRepositoryOption func(o *SettingsRepositoryOptions)

// NewSettingsRepository creates and configures a new SettingsRepository.
func NewSettingsRepository(client Client, options ...SettingsRepositoryOption) (*SettingsRepository, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SettingsRepositoryOptions{
		Database:   defaultSettingsRepositoryDatabase,
		Collection: defaultSettingsRepositoryCollection,
		Timeout:    defaultSettingsRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	return &SettingsRepository{
		client:     client.Database(opts.Database),
		collection: opts.Collection,
		timeout:    opts.Timeout,
	}, nil
}

// Get settings.
func (r SettingsRepository) Get(ctx context.Context) (models.Settings, error) {
	res, err := r.client.Collection(r.collection).FindOne(ctx, bson.D{{Key: "_id", Value: "security"}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Settings{}, dberrors.ErrSettingsNotFound
		}
		return models.Settings{}, err
	}

	var s models.Security
	if err := res.Decode(&s); err != nil {
		return models.Settings{}, err
	}

	return models.Settings{Security: s}, nil
}

// Create settings.
func (r SettingsRepository) Create(ctx context.Context, settings models.Settings) (models.Settings, error) {
	if len(settings.Security.ID) == 0 {
		settings.Security.ID = "security"
	}

	_, err := r.client.Collection(r.collection).InsertOne(ctx, settings.Security)
	if err != nil {
		return models.Settings{}, err
	}

	return r.Get(ctx)
}

// Update settings.
func (r SettingsRepository) Update(ctx context.Context, settings models.Settings) (models.Settings, error) {
	if err := r.client.Collection(r.collection).UpdateOne(ctx, bson.D{{Key: "_id", Value: "security"}}, settings.Security); err != nil {
		return models.Settings{}, err
	}
	return r.Get(ctx)
}
