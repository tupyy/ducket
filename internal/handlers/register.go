package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func RegisterHandlers(r *gin.RouterGroup) {
	validate = validator.New()

	transactionHandlers(r)
	tagHandlers(r)
	rulesHandlers(r)
}
