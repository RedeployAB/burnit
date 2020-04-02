package db

// db.go provides interfaces for collection based
// databases.

import (
	"context"

	"github.com/RedeployAB/burnit/burnitdb/internal/models"
)

// Connection provides methods Connect, Disconnect and Database.
type Connection interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Database(name string) Database
}

// Database provides method Collection.
type Database interface {
	Collection(name string) Collection
}

// Collection provides methods FindOne, InsertOne,
// DeleteOne and DeleteMany.
type Collection interface {
	FindOne(id string) (*models.Secret, error)
	InsertOne(s *models.Secret) (*models.Secret, error)
	DeleteOne(id string) (int64, error)
	DeleteMany() (int64, error)
}
