package pg

import sq "github.com/Masterminds/squirrel"

const (
	rulesTable              = "rules"
	labelsTable             = "labels"
	transactionTable        = "transactions"
	transactionsLabelsTable = "transactions_labels"
	rulesLabelsTable        = "rules_labels"

	colID                 = "id"
	colDate               = "date"
	colTransactionType    = "kind"
	colTransactionContent = "content"
	colTransactionAmount  = "amount"
	colTransactionAccount = "account"
	colTransactionID      = "transaction_id"
	colLabelID            = "label_id"
	colRuleName           = "name"
	colRulPattern         = "pattern"
	colLabelValue         = "value"
	colLabelKey           = "key"
	colRuleID             = "rule_id"
	colFilterRuleID       = "rules.id"
	colCreatedAt          = "created_at"
	colHash               = "hash"

	errUnableToWriteLabel  = "unable to write label: %w"
	errUnableToDeleteLabel = "unable to delete label: %w"
	errUnableToReadLabel   = "unable to read label: %w"

	errUnableToReadRule   = "unable to read rule: %w"
	errUnableToDeleteRule = "unable to delete rule: %w"
	errUnableToWriteRule  = "unable to write rule: %w"

	errUnableToDeleteTransaction = "unable to delete transaction: %w"
	errUnableToWriteTransaction  = "unable to write transaction: %w"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	selectRulesStmt = psql.Select(
		"rules.id as rule_id",
		"rules.pattern",
		"b.id as label_id",
		"b.key",
		"b.value").
		From(rulesTable).
		LeftJoin("(SELECT * FROM labels JOIN rules_labels as a on a.label_id = labels.id) as b ON b.rule_id = rules.id")

	selectLabelsStmt = psql.Select(
		"labels.id",
		colLabelKey,
		colLabelValue,
		colRuleID,
		colCreatedAt,
	).
		From(labelsTable).
		LeftJoin("rules_labels on rules_labels.label_id = labels.id")

	selectTransactionLabelsStmt = psql.Select("*").From(transactionsLabelsTable)

	countTransactionsPerLabelPerRuleStmt = psql.
						Select("b.id as label_id", "b.rule_id", "COUNT(transaction_id)").
						FromSelect(selectLabelsStmt, "b").
						InnerJoin("transactions_labels on transactions_labels.label_id = b.id").
						GroupBy("b.id", "b.rule_id")

	insertTransaction = psql.Insert(transactionTable).
				Columns(
			colDate,
			colTransactionAccount,
			colHash,
			colTransactionType,
			colTransactionContent,
			colTransactionAmount,
		)

	insertTransactionLabel = psql.Insert(transactionsLabelsTable).Columns(colTransactionID, colLabelID)

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
	).
		From(transactionTable).
		LeftJoin(`(select a.transaction_id,labels.id as label_id,labels.key,labels.value from transactions_labels as a INNER JOIN labels on a.label_id = labels.id) as b on b.transaction_id = transactions.id`)

	insertRule  = psql.Insert(rulesTable).Columns("id", "pattern")
	updateRule  = psql.Update(rulesTable)
	insertLabel = psql.Insert("labels").Columns(colLabelKey, colLabelValue)
)
