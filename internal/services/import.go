package services

import (
	"context"
	"fmt"
	"io"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/pkg/reader"
	"go.uber.org/zap"
)

type ImportService struct {
	dt *pg.Datastore
}

// NewImportService creates a new instance of ImportService with the provided datastore.
func NewImportService(dt *pg.Datastore) *ImportService {
	return &ImportService{dt: dt}
}

// ImportResult represents the result of a file import operation
type ImportResult struct {
	Filename      string   `json:"filename"`
	TotalRows     int      `json:"total_rows"`
	ProcessedRows int      `json:"processed_rows"`
	CreatedCount  int      `json:"created_count"`
	UpdatedCount  int      `json:"ignored_count"`
	ErrorCount    int      `json:"error_count"`
	Errors        []string `json:"errors,omitempty"`
}

// FileUpload represents an uploaded file with its content
type FileUpload struct {
	Filename string
	Content  io.Reader
}

// ImportFiles processes multiple uploaded files and imports their transaction data
func (s *ImportService) ImportFiles(ctx context.Context, files []FileUpload) ([]ImportResult, error) {
	results := make([]ImportResult, 0, len(files))

	for _, file := range files {
		result := s.importSingleFile(ctx, file)
		results = append(results, result)
	}

	return results, nil
}

// importSingleFile processes a single file and returns the import result
func (s *ImportService) importSingleFile(ctx context.Context, file FileUpload) ImportResult {
	result := ImportResult{
		Filename: file.Filename,
		Errors:   make([]string, 0),
	}

	// Read transactions from file
	transactions, err := reader.ReadTransactionsFromFile(file.Filename, file.Content)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to read file: %v", err))
		result.ErrorCount = 1
		return result
	}

	result.TotalRows = len(transactions)

	// Apply rules to all transactions first
	ruleApplierService := NewRuleApplierService(s.dt)
	transactionsWithRules, err := ruleApplierService.ApplyRulesToTransactions(ctx, transactions)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Error applying rules to transactions: %v", err))
		result.ErrorCount++
		return result
	}

	// Process each transaction
	transactionService := NewTransactionService(s.dt)

	for _, transaction := range transactionsWithRules {
		result.ProcessedRows++

		// Check if transaction already exists
		existingTransaction, err := transactionService.GetTransaction(ctx, transaction.Hash)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error checking transaction %s: %v", transaction.Hash, err))
			result.ErrorCount++
			continue
		}

		// Create or update transaction (now with applied rule tags)
		_, err = transactionService.CreateOrUpdate(ctx, transaction)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error saving transaction %s: %v", transaction.Hash, err))
			result.ErrorCount++
			continue
		}

		if existingTransaction != nil {
			result.UpdatedCount++
		} else {
			result.CreatedCount++
		}

		zap.S().Debugw("Processed transaction with rules",
			"hash", transaction.Hash,
			"amount", transaction.Amount,
			"date", transaction.Date,
			"content", transaction.RawContent,
			"tags", transaction.Tags,
		)
	}

	zap.S().Infow("File import completed",
		"filename", result.Filename,
		"total_rows", result.TotalRows,
		"processed_rows", result.ProcessedRows,
		"created_count", result.CreatedCount,
		"ignored_count", result.UpdatedCount,
		"error_count", result.ErrorCount,
	)

	return result
}
