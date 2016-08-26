package datastore

import "timezones_mc/revel_app/app/models"

type Datastore interface {
	AddDocument(doc models.Document) error
}