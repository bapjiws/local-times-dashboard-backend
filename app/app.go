package main

import (
	"github.com/bapjiws/timezones_mc/app/handlers/city"
	"github.com/bapjiws/timezones_mc/app/middleware/shared"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	"github.com/gin-gonic/gin"
)

var ES *elasticsearch.ElasticStore

// TODO: create an object for context and initialize it here
func init() {
	ES = elasticsearch.NewElasticStore(configs.CityStoreConfig)
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	r.Use(shared.SetContext(ES))

	r.GET("/city", city.SuggestCities)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
