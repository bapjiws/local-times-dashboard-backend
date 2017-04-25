package datastore

import (
	"github.com/bapjiws/local_times_dashboard_backend/models/document"
	"github.com/bapjiws/local_times_dashboard_backend/models/suggest"
)

type Datastore interface {
	AddDocument(doc document.Document) error
	FindDocumentById(id string) (document.Document, error)
	SuggestDocuments(s suggest.Suggest) ([]document.Document, error)
}
