package pg

import sq "github.com/Masterminds/squirrel"

const (
	rulesTable              = "rules"
	labelsTable             = "labels"
	transactionTable        = "transactions"
	transactionsLabelsTable = "transactions_labels"
	rulesLabelsTable        = "rules_labels"

	// columns
	// Common
	colID        = "id"
	colCreatedAt = "created_at"

	// transaction table
	colDate               = "date"
	colTransactionType    = "kind"
	colTransactionContent = "content"
	colTransactionAmount  = "amount"
	colTransactionAccount = "account"
	colTransactionID      = "transaction_id"
	colHash               = "hash"
	colLabelID            = "label_id"

	// Rule table
	colRuleName   = "name"    // unused
	colRulPattern = "pattern" // unused

	// Labels table
	colLabelValue   = "value"
	colLabelKey     = "key"
	colRuleID       = "rule_id"
	colFilterRuleID = "rules.id"

	errUnableToWriteLabel  = "unable to write label: %w"
	errUnableToDeleteLabel = "unable to delete label: %w"
	errUnableToReadLabel   = "unable to read label: %w"

	errUnableToReadRule   = "unable to read rule: %w"
	errUnableToDeleteRule = "unable to delete rule: %w"
	errUnableToWriteRule  = "unable to write rule: %w"

	errUnableToDeleteTransaction  = "unable to delete transaction: %w"
	errUnableToWriteTransaction   = "unable to write transaction: %w"
	errUnableToCountTransactions  = "unable to count transactions: %w"
	errUnableToWriteRelationship  = "unable to write relationship: %w"
	errUnableToDeleteRelationship = "unable to delete relationship: %w"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// selectRulesStmt retrieves all rules with their associated labels via left join.
	// Example SQL: SELECT rules.id as rule_id, rules.pattern, b.id as label_id, b.key, b.value
	//              FROM rules
	//              LEFT JOIN (SELECT * FROM labels JOIN rules_labels as a on a.label_id = labels.id) as b ON b.rule_id = rules.id
	selectRulesStmt = psql.Select(
		"rules.id as rule_id",
		"rules.pattern",
		"b.id as label_id",
		"b.key",
		"b.value").
		From(rulesTable).
		LeftJoin("(SELECT * FROM labels JOIN rules_labels as a on a.label_id = labels.id) as b ON b.rule_id = rules.id")

	// selectLabelsStmt retrieves all labels with their associated rule IDs via left join.
	// Example SQL: SELECT labels.id, key, value, rule_id, created_at
	//              FROM labels
	//              LEFT JOIN rules_labels on rules_labels.label_id = labels.id
	selectLabelsStmt = psql.Select(
		"labels.id",
		colLabelKey,
		colLabelValue,
		colRuleID,
		colCreatedAt,
	).
		From(labelsTable).
		LeftJoin("rules_labels on rules_labels.label_id = labels.id")

	// selectTransactionLabelsStmt retrieves all transaction-label associations.
	// Example SQL: SELECT * FROM transactions_labels
	selectTransactionLabelsStmt = psql.Select("*").From(transactionsLabelsTable)

	countAllTransactions        = psql.Select("COUNT(*)").From(transactionTable)
	countTransactionsWithFilter = psql.Select("COUNT(*)").From(transactionsLabelsTable)

	// insertTransaction inserts a new transaction record.
	// Example SQL: INSERT INTO transactions (date, account, hash, kind, content, amount) VALUES ($1, $2, $3, $4, $5, $6)
	insertTransaction = psql.Insert(transactionTable).
				Columns(
			colDate,
			colTransactionAccount,
			colHash,
			colTransactionType,
			colTransactionContent,
			colTransactionAmount,
		)

	// insertTransactionLabel creates an association between a transaction and a label.
	// Example SQL: INSERT INTO transactions_labels (transaction_id, label_id) VALUES ($1, $2)
	insertTransactionLabelRuleRelationship = psql.Insert(transactionsLabelsTable).Columns(colTransactionID, colLabelID, colRuleID)
	insertTransactionLabelRelationship     = psql.Insert(transactionsLabelsTable).Columns(colTransactionID, colLabelID)
	insertLabelRuleRelationship            = psql.Insert(rulesLabelsTable).Columns(colLabelID, colRuleID)

	deleteTransactionLabelRuleRelationship = psql.Delete(transactionsLabelsTable)
	deleteLabelRuleRelationship            = psql.Delete(rulesLabelsTable)

	// selectTransactionStmp retrieves transactions with their associated labels via left join.
	// Example SQL: SELECT id, date, account, kind, content, amount, hash, label_id, key, value
	//              FROM transactions
	//              LEFT JOIN (
	//					select a.transaction_id,labels.id as label_id,labels.key,labels.value from transactions_labels as a INNER JOIN labels on a.label_id = labels.id
	//				) as b on b.transaction_id = transactions.id
	selectTransactionStmp = psql.Select(
		colID,
		colDate,
		colTransactionAccount,
		colTransactionType,
		colTransactionContent,
		colTransactionAmount,
		colHash,
		colLabelID,
		colLabelKey,
		colLabelValue,
		colRuleID,
	).
		From(transactionTable).
		LeftJoin(`(select a.transaction_id,labels.id as label_id,labels.key,labels.value, a.rule_id from transactions_labels as a INNER JOIN labels on a.label_id = labels.id) as b on b.transaction_id = transactions.id`)

	// insertRule inserts a new rule record.
	// Example SQL: INSERT INTO rules (id, pattern) VALUES ($1, $2)
	insertRule = psql.Insert(rulesTable).Columns("id", "pattern")

	// updateRule provides base statement for updating rule records.
	// Example SQL: UPDATE rules SET pattern = $1 WHERE id = $2
	updateRule = psql.Update(rulesTable)

	// insertLabel inserts a new label record.
	// Example SQL: INSERT INTO labels (key, value) VALUES ($1, $2)
	insertLabel = psql.Insert("labels").Columns(colLabelKey, colLabelValue)
)
