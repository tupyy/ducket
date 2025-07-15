package handlers

import (
	v1 "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
)

// RegisterApiV1Handlers registers all version 1 API handlers with the provided router group.
// This includes handlers for transactions, rules, tags, summary, and import endpoints.
func RegisterApiV1Handlers(r *gin.RouterGroup) {
	v1.TransactionHandlers(r)
	v1.LabelHandlers(r)
	v1.RulesHandlers(r)
	v1.SummaryHandlers(r)
	v1.ImportHandlers(r)
}
