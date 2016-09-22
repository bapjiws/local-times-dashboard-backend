package main

import (
	"github.com/bapjiws/timezones_mc/app/handlers"
	"github.com/bapjiws/timezones_mc/app/middleware"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	"github.com/gin-gonic/gin"
)

const (
	API_BASE = "api"
	PORT     = ":8888"
)

var ES *elasticsearch.ElasticStore

// TODO: create an object/map for context and initialize it here
func init() {
	ES = elasticsearch.NewElasticStore(configs.CityStoreConfig)
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.Use(middleware.SetContext(ES))

	cityRouter := router.Group(API_BASE)
	cityRouter.GET("/city", handlers.SuggestCities)
	cityRouter.GET("/city/:id", handlers.FindCityById)

	router.Run(PORT)
}
