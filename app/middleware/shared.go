package middleware

import (
	"github.com/bapjiws/local_times_dashboard_backend/models/datastore"
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

func AllowCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
	}
}
