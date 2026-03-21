package handlers

import (
	"fmt"
	"net/http"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	pkgErrors "git.tls.tupangiu.ro/cosmin/finante/pkg/errors"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/reader"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ImportTransactions(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expected multipart form with files"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	var account int64
	if vals, ok := form.Value["account"]; ok && len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &account)
	}

	var results []v1.ImportFileResult

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			results = append(results, v1.ImportFileResult{Filename: fh.Filename, Errors: 1})
			continue
		}

		transactions, err := reader.ReadTransactionsFromFile(fh.Filename, account, f)
		f.Close()
		if err != nil {
			zap.S().Errorw("failed to parse file", "filename", fh.Filename, "error", err)
			results = append(results, v1.ImportFileResult{Filename: fh.Filename, Errors: 1})
			continue
		}

		created, skipped, errCount := 0, 0, 0
		for _, txn := range transactions {
			_, err := h.txnSvc.Create(c.Request.Context(), txn)
			if err != nil {
				if pkgErrors.IsDuplicateResourceError(err) {
					skipped++
				} else {
					zap.S().Errorw("failed to create transaction", "error", err)
					errCount++
				}
				continue
			}
			created++
		}

		results = append(results, v1.ImportFileResult{
			Filename: fh.Filename,
			Created:  created,
			Skipped:  skipped,
			Errors:   errCount,
		})
	}

	totalCreated := 0
	for _, r := range results {
		totalCreated += r.Created
	}

	c.JSON(http.StatusOK, v1.ImportResponse{
		Files:   results,
		Message: fmt.Sprintf("imported %d transactions from %d file(s)", totalCreated, len(files)),
	})
}
