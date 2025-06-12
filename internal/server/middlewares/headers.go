package middlewares

import "github.com/gin-gonic/gin"

func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}
