package middleware

import (
	"github.com/bapjiws/timezones_mc/models/datastore"
	"github.com/gin-gonic/gin"
)

// TODO: create and object for storing the context?
func SetContext(ds datastore.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Datastore", ds)
	}
}
