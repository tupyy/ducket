package handlers

import "github.com/gin-gonic/gin"

func transactionHandlers(r *gin.RouterGroup) {
	r.GET("/transactions", func(c *gin.Context) {
		c.JSON(200, "")
	})
}
