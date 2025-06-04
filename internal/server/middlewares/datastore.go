package middlewares

import (
	"context"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"github.com/gin-gonic/gin"
)

var DatastoreKey = "datastore"

func DatastoreMiddleware(dt *pg.Datastore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), DatastoreKey, dt)
		c.Request = c.Request.WithContext(ctx)
		c.Set(DatastoreKey, dt)
		c.Next()
	}
}
