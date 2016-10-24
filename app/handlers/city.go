package handlers

import (
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/models/suggest"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SuggestCities(c *gin.Context) {
	ES := c.MustGet("Datastore").(*elasticsearch.ElasticStore)

	response, err := ES.SuggestDocuments(suggest.Suggest{
		SuggesterName: "city_suggest",
		Text:          c.Query("name"), // shortcut for c.Request.URL.Query().Get("name")
		Field:         "suggest",
		PayloadKeys: map[string]string{
			"city_id": "city_id",
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	c.JSON(http.StatusOK, response)

}

func FindCityById(c *gin.Context) {
	ES := c.MustGet("Datastore").(*elasticsearch.ElasticStore)
	id := c.Param("id") // shortcut for c.Request.URL.Query().Get("name")

	response, err := ES.FindDocumentById(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	// TODO: turn generic response into a City model here (e.g, to get rid of "suggest")?
	c.JSON(http.StatusOK, response)
}
