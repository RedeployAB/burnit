package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// defaultSessionStoreDatabase is the default database for the SessionStore.
	defaultSessionStoreDatabase = "burnit"
	// defaultSessionStoreCollection is the default collection for the SessionStore.
	defaultSessionStoreCollection = "sessions"
	// defaultSessionStoreTimeout is the default timeout for the SessionStore.
	defaultSessionStoreTimeout = 10 * time.Second
)

// sessionStore is a MongoDB implementation of a SessionStore.
type sessionStore struct {
	client        Client
	collection    string
	upsertSession upsertSessionFunc
	timeout       time.Duration
}

// SessionStoreOptions is the options for the SessionStore.
type SessionStoreOptions struct {
	Database   string
	Collection string
	Timeout    time.Duration
}

// SessionStoreOption is a function that sets options for the SessionStore.
type SessionStoreOption func(o *SessionStoreOptions)

// NewSessionStore creates and configures a new SessionStore.
func NewSessionStore(client Client, options ...SessionStoreOption) (*sessionStore, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SessionStoreOptions{
		Database:   defaultSessionStoreDatabase,
		Collection: defaultSessionStoreCollection,
		Timeout:    defaultSessionStoreTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	if len(opts.Database) == 0 {
		return nil, ErrDatabaseNotSet
	}

	store := &sessionStore{
		client:     client.Database(opts.Database),
		collection: opts.Collection,
		timeout:    opts.Timeout,
	}

	if client.ReplicaSetEnabled() {
		setUpsertSessionWithTransaction(store)
	} else {
		setUpsertSession(store)
	}

	return store, nil
}

// Get a session by its ID.
func (s sessionStore) Get(ctx context.Context, id string) (db.Session, error) {
	res, err := s.client.Collection(s.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		if errors.Is(err, ErrNoDocuments) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}

	var session db.Session
	if err := res.Decode(&session); err != nil {
		return db.Session{}, err
	}
	return session, nil
}

// GetByCSRFToken gets a session by its CSRF token.
func (s sessionStore) GetByCSRFToken(ctx context.Context, token string) (db.Session, error) {
	res, err := s.client.Collection(s.collection).FindOne(ctx, bson.D{{Key: "csrf.token", Value: token}})
	if err != nil {
		if errors.Is(err, ErrNoDocuments) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}

	var session db.Session
	if err := res.Decode(&session); err != nil {
		return db.Session{}, err
	}
	return session, nil
}

// Upsert a session. Create the session if it does not exist, otherwise update it.
func (s sessionStore) Upsert(ctx context.Context, session db.Session) (db.Session, error) {
	return s.upsertSession(ctx, session)
}

// Delete a session by its ID.
func (s sessionStore) Delete(ctx context.Context, id string) error {
	if err := s.client.Collection(s.collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		if errors.Is(err, ErrDocumentNotDeleted) {
			return dberrors.ErrSessionNotFound
		}
		return err
	}
	return nil
}

// DeleteByCSRFToken deletes a session by its CSRF token.
func (s sessionStore) DeleteByCSRFToken(ctx context.Context, token string) error {
	if err := s.client.Collection(s.collection).DeleteOne(ctx, bson.D{{Key: "csrf.token", Value: token}}); err != nil {
		if errors.Is(err, ErrDocumentNotDeleted) {
			return dberrors.ErrSessionNotFound
		}
		return err
	}
	return nil
}

// DeleteExpired deletes all expired sessions.
func (s sessionStore) DeleteExpired(ctx context.Context) error {
	filter := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: now()}}}}
	err := s.client.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		if errors.Is(err, ErrDocumentsNotDeleted) {
			return dberrors.ErrSessionsNotDeleted
		}
		return err
	}
	return nil
}

// Close the store and its underlying connections.
func (s sessionStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.Disconnect(ctx)
}

// upsertSessionFunc is a function that upserts a session.
type upsertSessionFunc func(ctx context.Context, session db.Session) (db.Session, error)

// setUpsertSession sets an upsertSessionFunc to the store.
func setUpsertSession(store *sessionStore) {
	store.upsertSession = func(ctx context.Context, session db.Session) (db.Session, error) {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "expiresAt", Value: session.ExpiresAt},
				{Key: "csrf.token", Value: session.CSRF.Token},
				{Key: "csrf.expiresAt", Value: session.CSRF.ExpiresAt},
			}},
		}

		id, err := store.client.Collection(store.collection).UpsertOne(ctx, bson.D{{Key: "_id", Value: session.ID}}, update)
		if err != nil {
			return db.Session{}, err
		}
		return store.Get(ctx, id)
	}
}

// setUpsertSessionWithTransaction sets an upsertSecretFunc for use with transactions
// to the store.
func setUpsertSessionWithTransaction(store *sessionStore) {
	store.upsertSession = func(ctx context.Context, session db.Session) (db.Session, error) {
		result, err := store.client.WithTransaction(ctx, func(ctx context.Context, client Client) (any, error) {
			update := bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "expiresAt", Value: session.ExpiresAt},
					{Key: "csrf.token", Value: session.CSRF.Token},
					{Key: "csrf.expiresAt", Value: session.CSRF.ExpiresAt},
				}},
			}

			id, err := store.client.Collection(store.collection).UpsertOne(ctx, bson.D{{Key: "_id", Value: session.ID}}, update)
			if err != nil {
				return db.Session{}, err
			}

			res, err := store.client.Collection(store.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
			if err != nil {
				if errors.Is(err, ErrNoDocuments) {
					return db.Session{}, dberrors.ErrSessionNotFound
				}
				return db.Session{}, err
			}

			var session db.Session
			if err := res.Decode(&session); err != nil {
				return db.Session{}, err
			}
			return session, nil
		})
		if err != nil {
			return db.Session{}, nil
		}

		session, ok := result.(db.Session)
		if !ok {
			return db.Session{}, errors.New("invalid document for session")
		}
		return session, nil
	}
}
