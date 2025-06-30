package v1

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"git.tls.tupangiu.ro/cosmin/finante/internal/services"
	dtContext "git.tls.tupangiu.ro/cosmin/finante/pkg/context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	maxFileSize = 10 << 20 // 10MB
)

// ImportHandlers registers the file import HTTP handlers with the provided router group.
// This includes the endpoint for uploading and processing transaction files.
func ImportHandlers(r *gin.RouterGroup) {
	r.POST("/import", func(c *gin.Context) {
		// Parse multipart form
		err := c.Request.ParseMultipartForm(maxFileSize)
		if err != nil {
			zap.S().Errorw("failed to parse multipart form", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
			return
		}

		// Get uploaded files
		form := c.Request.MultipartForm
		files := form.File["files"]

		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
			return
		}

		// Validate file types and prepare file uploads
		var fileUploads []services.FileUpload
		supportedExtensions := map[string]bool{
			".xlsx": true,
			".xls":  true,
			".csv":  true,
		}

		for _, fileHeader := range files {
			// Check file size
			if fileHeader.Size > maxFileSize {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("File '%s' exceeds maximum size of %d bytes", fileHeader.Filename, maxFileSize),
				})
				return
			}

			// Check file extension
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !supportedExtensions[ext] {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("File '%s' has unsupported extension '%s'. Supported: .xlsx, .xls, .csv", fileHeader.Filename, ext),
				})
				return
			}

			// Open file
			file, err := fileHeader.Open()
			if err != nil {
				zap.S().Errorw("failed to open uploaded file", "filename", fileHeader.Filename, "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to open file '%s'", fileHeader.Filename),
				})
				return
			}

			fileUploads = append(fileUploads, services.FileUpload{
				Filename: fileHeader.Filename,
				Content:  file,
			})
		}

		// Make sure to close all files after processing
		defer func() {
			for _, fileUpload := range fileUploads {
				if closer, ok := fileUpload.Content.(interface{ Close() error }); ok {
					closer.Close()
				}
			}
		}()

		// Process files using import service
		dt := dtContext.MustFromContext(c)
		importService := services.NewImportService(dt)

		results, err := importService.ImportFiles(c.Request.Context(), fileUploads)
		if err != nil {
			zap.S().Errorw("failed to import files", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process files"})
			return
		}

		// Calculate summary statistics
		summary := calculateSummary(results)

		zap.S().Infow("Import operation completed",
			"files_count", len(results),
			"total_processed", summary["total_processed"],
			"total_created", summary["total_created"],
			"total_ignored", summary["total_ignored"],
			"total_errors", summary["total_errors"],
		)

		// Return results
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Processed %d files", len(results)),
			"summary": summary,
			"results": results,
		})
	})
}

// calculateSummary computes aggregate statistics from import results
func calculateSummary(results []services.ImportResult) map[string]int {
	summary := map[string]int{
		"files_processed": len(results),
		"total_rows":      0,
		"total_processed": 0,
		"total_created":   0,
		"total_ignored":   0,
		"total_errors":    0,
	}

	for _, result := range results {
		summary["total_rows"] += result.TotalRows
		summary["total_processed"] += result.ProcessedRows
		summary["total_created"] += result.CreatedCount
		summary["total_ignored"] += result.IgnoredCount
		summary["total_errors"] += result.ErrorCount
	}

	return summary
}
