package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	//"github.com/bapjiws/timezones_mc/app/handlers/city"
	"time"
	"log"
	"net/http"
)

var ES *elasticsearch.ElasticStore
//var c *gin.Context

func InitElasticsearch() {
	ES = elasticsearch.NewElasticStore(configs.CityStoreConfig)
}

// Move InitElasticsearch's stuff in here and remove the function itself?
func init() {
	InitElasticsearch()
}

//func main() {
//	router := gin.Default()
//
//
//	//c.Set("ES", "test")
//
//	router.GET("/", city.Hello)
//	router.GET("/city", city.SuggestCities)
//
//	router.Run()
//}


func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Set("ES", ES)

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		status := c.Writer.Status()
		log.Printf("Latency: %s Status: %d\n",  latency, status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/city", func(c *gin.Context) {
		ES := c.MustGet("ES").(*elasticsearch.ElasticStore)
		name := c.Query("name") // shortcut for c.Request.URL.Query().Get("name")

		response, err :=  ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")

		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		c.JSON(http.StatusOK, response)

	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}