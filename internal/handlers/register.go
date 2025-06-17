package handlers

import (
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
)

func RegisterApiV1Handlers(r *gin.RouterGroup) {
	validate = validator.New()

	v1.TransactionHandlers(r, validate)
	v1.TagHandlers(r, validate)
	v1.RulesHandlers(r, validate)
	v1.SummaryHandlers(r, validate)
}
