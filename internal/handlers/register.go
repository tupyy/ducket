package handlers

import (
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
)

func RegisterApiV1Handlers(r *gin.RouterGroup) {
	v1.TransactionHandlers(r)
	v1.TagHandlers(r)
	v1.RulesHandlers(r)
	v1.SummaryHandlers(r)
}
