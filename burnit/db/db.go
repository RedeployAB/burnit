package db

// db.go provides interfaces for collection based
// databases.

import (
	"context"
)

// Client wraps around methods Connect, Disconnect, Find, Insert, Delete
// and DeleteMany.
type Client interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Find(ctx context.Context, id string) (*Secret, error)
	Insert(ctx context.Context, s *Secret) (*Secret, error)
	Delete(ctx context.Context, id string) (int64, error)
	DeleteMany(ctx context.Context) (int64, error)
}
