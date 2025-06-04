package handlers

import "github.com/gin-gonic/gin"

func tagHandlers(r *gin.RouterGroup) {
	r.GET("/api/v1/tags", func(c *gin.Context) {})
}
