package shared

import (
	"log"
	"time"

	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	log.Println("USE Logger")

	return func(c *gin.Context) {
		t := time.Now()

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		status := c.Writer.Status()
		log.Printf("Latency: %s Status: %d\n", latency, status)
	}
}

// TODO: create and object for storing the context?
func SetContext(es *elasticsearch.ElasticStore) gin.HandlerFunc {
	log.Println("USE SetContext")

	return func(c *gin.Context) {
		c.Set("ES", es)
	}
}
