package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrNilClient is returned when the client is nil.
	ErrNilClient = errors.New("client is nil")
	// ErrDatabaseNotSet is returned when the database is not set.
	ErrDatabaseNotSet = errors.New("database not set")
	// ErrCollectionNotSet is returned when the collection is not set.
	ErrCollectionNotSet = errors.New("collection not set")
	// ErrDocumentNotUpdated is returned when a document is not updated.
	ErrDocumentNotUpdated = errors.New("document not updated")
	// ErrDocumentNotDeleted is returned when a document is not deleted.
	ErrDocumentNotDeleted = errors.New("document not deleted")
	// ErrDocumentsNotDeleted is returned when documents are not deleted.
	ErrDocumentsNotDeleted = errors.New("documents not deleted")
)

var (
	// ErrNoDocuments is returned when no documents are found.
	ErrNoDocuments = mongo.ErrNoDocuments
)
