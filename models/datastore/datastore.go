package datastore

import (
	"github.com/bapjiws/timezones_mc/models/document"
	"github.com/bapjiws/timezones_mc/models/suggest"
)

type Datastore interface {
	AddDocument(doc document.Document) error
	FindDocumentById(id string) (document.Document, error)
	SuggestDocuments(s suggest.Suggest) ([]document.Document, error)
}
