package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	"github.com/bapjiws/timezones_mc/app/handlers/city"
)

var ES *elasticsearch.ElasticStore
var ctx *gin.Context

func InitElasticsearch() {
	ES = elasticsearch.NewElasticStore(configs.CityStoreConfig)
}

// Move InitElasticsearch's stuff in here and remove the function itself?
func init() {
	InitElasticsearch()
}

func main() {
	router := gin.Default()


	//ctx.Set("ES", ES)

	router.GET("/", city.Hello)
	router.GET("/city", city.SuggestCities)

	router.Run()
}
