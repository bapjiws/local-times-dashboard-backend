package shared

import (
	"log"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/gin-gonic/gin"
)

// TODO: create and object for storing the context?
func SetContext(es *elasticsearch.ElasticStore) gin.HandlerFunc {
	log.Println("USE SetContext")

	return func(c *gin.Context) {
		c.Set("ES", es)
	}
}
