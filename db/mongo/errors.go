package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrCollectionNotSet is returned when the collection is not set.
	ErrCollectionNotSet = errors.New("collection not set")
	// ErrDocumentNotDeleted is returned when a document is not deleted.
	ErrDocumentNotDeleted = errors.New("document not deleted")
	// ErrDocumentsNotDeleted is returned when documents are not deleted.
	ErrDocumentsNotDeleted = errors.New("documents not deleted")
)

var (
	// ErrNoDocuments is returned when no documents are found.
	ErrNoDocuments = mongo.ErrNoDocuments
)
