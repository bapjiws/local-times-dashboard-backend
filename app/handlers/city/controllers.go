package city

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK,  map[string]interface{}{"msg": "Wazzup!"})
}

func SuggestCities(c *gin.Context) {
	//firstname := c.DefaultQuery("firstname", "Guest")
	//lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

	//c.String(http.StatusOK, "Hello %s %s", firstname, lastname)

	//ES := c.Value("ES") // TODO: ok interface type assertion

	name := c.Query("name") // shortcut for c.Request.URL.Query().Get("name")

	//response, err :=  ES.(string).SuggestDocuments("city_suggest", name, "suggest", "city_id")
	//
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	//}

	c.JSON(http.StatusOK, map[string]interface{}{"error": name})

}
