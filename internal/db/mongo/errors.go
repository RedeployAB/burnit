package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrDocumentNotUpserted is returned when a document is not upserted.
	ErrDocumentNotUpserted = errors.New("document not upserted")
	// ErrDocumentNotDeleted is returned when a document is not deleted.
	ErrDocumentNotDeleted = errors.New("document not deleted")
	// ErrDocumentsNotDeleted is returned when documents are not deleted.
	ErrDocumentsNotDeleted = errors.New("documents not deleted")
)

var (
	// ErrNoDocuments is returned when no documents are found.
	ErrNoDocuments = mongo.ErrNoDocuments
)
