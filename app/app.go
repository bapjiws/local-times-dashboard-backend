package main

import (
	//"net/http"
	"github.com/bapjiws/timezones_mc/app/handlers"
	"github.com/bapjiws/timezones_mc/app/middleware"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch"
	"github.com/bapjiws/timezones_mc/datastore/elasticsearch/configs"
	"github.com/gin-gonic/gin"
	//"github.com/itsjamie/gin-cors"
	//"time"
)

const (
	APIBASE = "api" // TODO: rename into CITY_API_BASE
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
	//, cors.Middleware(cors.Config{
	//	Origins:        "*",
	//	Methods:        "GET, PUT, POST, DELETE",
	//	RequestHeaders: "Origin, Authorization, Content-Type, access-control-allow-origin",
	//	ExposedHeaders: "",
	//	MaxAge: 50 * time.Second,
	//	Credentials: true,
	//	ValidateHeaders: false,
	//})

	router.GET("/city", handlers.SuggestCities)
	router.GET("/city/:id", handlers.FindCityById)

	routerGroup := router.Group(APIBASE)
	//routerGroup.OPTIONS("/city", preflight)
	routerGroup.GET("/city", handlers.SuggestCities)

	router.Run(PORT)
}

//func preflight(c *gin.Context) {
//	c.Header("Access-Control-Allow-Origin", "*")
//	c.Header("Access-Control-Allow-Headers", "access-control-allow-origin") // , access-control-allow-headers
//	c.JSON(http.StatusOK, struct{}{})
//}
