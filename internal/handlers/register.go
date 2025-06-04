package handlers

import "github.com/gin-gonic/gin"

func RegisterHandlers(r *gin.RouterGroup) {
	transactionHandlers(r)
	tagHandlers(r)
}
