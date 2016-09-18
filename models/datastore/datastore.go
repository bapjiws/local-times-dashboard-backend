package datastore

import (
	"github.com/bapjiws/timezones_mc/models/document"
)

type Datastore interface {
	AddDocument(doc document.Document) error
	FindDocumentById(id string) (document.Document, error)
	SuggestDocuments(suggesterName string, text string, field string, payloadKey string) ([]document.Document, error)
}
