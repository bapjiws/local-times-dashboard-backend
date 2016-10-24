package middleware

import (
	"github.com/bapjiws/timezones_mc/models/datastore"
	"github.com/gin-gonic/gin"
)

type Context struct {
	DS datastore.Datastore
}

func SetContext(ctx Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Datastore", ctx.DS)
	}
}
