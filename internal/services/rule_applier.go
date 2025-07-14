package services

import (
	"context"
	"regexp"

	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	"git.tls.tupangiu.ro/cosmin/finante/internal/entity"
	"go.uber.org/zap"
)

type RuleApplierService struct {
	dt *pg.Datastore
}

// NewRuleApplierService creates a new instance of RuleApplierService
func NewRuleApplierService(dt *pg.Datastore) *RuleApplierService {
	return &RuleApplierService{dt: dt}
}

// ApplyAllRulesToAllTransactions applies all rules to all transactions in the database
func (ras *RuleApplierService) ApplyAllRulesToAllTransactions(ctx context.Context) error {
	// Get all rules and transactions
	ruleService := NewRuleService(ras.dt)
	transactionService := NewTransactionService(ras.dt)

	rules, err := ruleService.GetRules(ctx)
	if err != nil {
		return err
	}

	// Get all transactions (without filters)
	transactions, err := transactionService.GetTransactions(ctx, &TransactionFilter{})
	if err != nil {
		return err
	}

	zap.S().Infow("Starting rule application",
		"total_rules", len(rules),
		"total_transactions", len(transactions),
	)

	// Compile all rule patterns once for efficiency
	compiledRules := make(map[string]*regexp.Regexp)
	ruleMap := make(map[string]entity.Rule)

	for _, rule := range rules {
		regex, err := regexp.Compile("(?i)" + rule.Pattern)
		if err != nil {
			zap.S().Warnw("Failed to compile rule pattern",
				"rule", rule.Name,
				"pattern", rule.Pattern,
				"error", err,
			)
			continue
		}
		compiledRules[rule.Name] = regex
		ruleMap[rule.Name] = rule
	}

	// Apply rules to each transaction
	matchCount := 0
	for _, transaction := range transactions {
		hasMatches := false
		updatedTransaction := transaction

		// Clear existing labels to reapply all rules
		updatedTransaction.Labels = make(map[int]string)

		for ruleName, regex := range compiledRules {
			rule := ruleMap[ruleName]

			if regex.MatchString(transaction.RawContent) {
				// Apply all labels from this rule to the transaction
				for _, label := range rule.Labels {
					updatedTransaction.Labels[label.ID] = rule.Name
					zap.S().Debugw("Applied rule to transaction",
						"rule", rule.Name,
						"label_key", label.Key,
						"label_value", label.Value,
						"label_id", label.ID,
						"transaction_id", transaction.ID,
						"transaction_content", transaction.RawContent,
					)
				}
				hasMatches = true
				matchCount++
			}
		}

		// Update the transaction if there were matches or if we need to clear old tags
		if hasMatches || len(transaction.Labels) > 0 {
			_, err := transactionService.CreateOrUpdate(ctx, updatedTransaction)
			if err != nil {
				zap.S().Errorw("Failed to update transaction with rule tags",
					"transaction_id", transaction.ID,
					"error", err,
				)
				continue
			}
		}
	}

	zap.S().Infow("Rule application completed",
		"total_matches", matchCount,
		"processed_transactions", len(transactions),
	)

	return nil
}

// ApplyRulesToTransaction applies all rules to a single transaction and returns the updated transaction
func (ras *RuleApplierService) ApplyRulesToTransaction(ctx context.Context, transaction entity.Transaction) (entity.Transaction, error) {
	ruleService := NewRuleService(ras.dt)

	// Get all rules from the database
	rules, err := ruleService.GetRules(ctx)
	if err != nil {
		return transaction, err
	}

	// Apply each rule that matches the transaction content
	for _, rule := range rules {
		matched, err := ras.MatchesRule(transaction.RawContent, rule.Pattern)
		if err != nil {
			zap.S().Warnw("Failed to compile rule pattern",
				"rule", rule.Name,
				"pattern", rule.Pattern,
				"error", err,
			)
			continue
		}

		if !matched {
			continue
		}

		// Apply all labels from this rule to the transaction
		for _, label := range rule.Labels {
			transaction.Labels[label.ID] = rule.Name
			zap.S().Debugw("Applied rule to transaction",
				"rule", rule.Name,
				"label_key", label.Key,
				"label_value", label.Value,
				"label_id", label.ID,
				"transaction_content", transaction.RawContent,
			)
		}
	}

	return transaction, nil
}

// ApplyRulesToTransactions applies rules to multiple transactions in batch
func (ras *RuleApplierService) ApplyRulesToTransactions(ctx context.Context, transactions []entity.Transaction) ([]entity.Transaction, error) {
	ruleService := NewRuleService(ras.dt)

	// Get all rules once to avoid multiple database calls
	rules, err := ruleService.GetRules(ctx)
	if err != nil {
		return transactions, err
	}

	// Compile all rule patterns once for efficiency
	compiledRules := make(map[string]*regexp.Regexp)
	ruleMap := make(map[string]entity.Rule)

	for _, rule := range rules {
		regex, err := regexp.Compile("(?i)" + rule.Pattern)
		if err != nil {
			zap.S().Warnw("Failed to compile rule pattern",
				"rule", rule.Name,
				"pattern", rule.Pattern,
				"error", err,
			)
			continue
		}
		compiledRules[rule.Name] = regex
		ruleMap[rule.Name] = rule
	}

	// Apply rules to each transaction
	for i := range transactions {
		for ruleName, regex := range compiledRules {
			rule := ruleMap[ruleName]

			if regex.MatchString(transactions[i].RawContent) {
				// Apply all labels from this rule to the transaction
				for _, label := range rule.Labels {
					transactions[i].Labels[label.ID] = rule.Name
					zap.S().Debugw("Applied rule to transaction",
						"rule", rule.Name,
						"label_key", label.Key,
						"label_value", label.Value,
						"label_id", label.ID,
						"transaction_content", transactions[i].RawContent,
					)
				}
			}
		}
	}

	return transactions, nil
}

// MatchesRule checks if the transaction content matches the rule pattern
func (ras *RuleApplierService) MatchesRule(content, pattern string) (bool, error) {
	// Compile the regex pattern (case-insensitive)
	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		return false, err
	}

	// Check if the pattern matches the content
	return regex.MatchString(content), nil
}
