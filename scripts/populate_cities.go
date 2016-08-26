package main

import (
	"timezones_mc/datastore/elasticsearch"
	"timezones_mc/datastore/elasticsearch/configs"
	"timezones_mc/revel_app/app/models"
)

func panicOnError(e error) {
	if e != nil {
		panic(e.Error())
	}
}

func main() {
	esStore := elasticsearch.NewElasticStore(configs.CityStoreConfig)

	err := esStore.Reindex()
	panicOnError(err)

	tempCity := models.City{
		CountryCode: "us",
		Name:        "NYC",
		Latitude:    123.45,
		Longitude:   -123.45,
	}

	err = esStore.AddDocument(tempCity)
	panicOnError(err)
}
