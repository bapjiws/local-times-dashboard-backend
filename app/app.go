package main

import (
	"net/http"

	"github.com/bapjiws/timezones_mc/app/handlers"
	"github.com/bapjiws/timezones_mc/app/middleware"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	"github.com/gin-gonic/gin"
)

const (
	APIBASE = "api"
	PORT    = ":8888"
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

	router.GET("/city", handlers.SuggestCities)
	router.GET("/city/:id", handlers.FindCityById)

	routerGroup := router.Group(APIBASE)
	routerGroup.OPTIONS("/city", preflight)
	routerGroup.GET("/city", handlers.SuggestCities)

	// Listen and server on 0.0.0.0:8888
	router.Run(PORT)
}

func preflight(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "http://localhost")
	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
	c.JSON(http.StatusOK, struct{}{})
}
