package city

import (
	"net/http"

	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/gin-gonic/gin"
)

func SuggestCities(c *gin.Context) {
	ES := c.MustGet("ES").(*elasticsearch.ElasticStore)
	name := c.Query("name") // shortcut for c.Request.URL.Query().Get("name")

	response, err := ES.SuggestDocuments("city_suggest", name, "suggest", "city_id")

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	c.JSON(http.StatusOK, response)

}
