package datastore

import "github.com/bapjiws/timezones_mc/models"

type Datastore interface {
	AddDocument(doc models.Document) error
	FindDocumentById(id string) (models.Document, error)
	SuggestDocuments(suggesterName string, text string, field string, payloadKey string) ([]models.Document, error)
}