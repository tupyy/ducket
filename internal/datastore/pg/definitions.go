package pg

import sq "github.com/Masterminds/squirrel"

const (
	rulesTable            = "rules"
	tagsTable             = "tags"
	transactionTable      = "transactions"
	transactionsTagsTable = "transactions_tags"
	rulesTagsTable        = "rules_tags"

	colID                 = "id"
	colDate               = "date"
	colTransactionType    = "kind"
	colTransactionContent = "content"
	colTransactionAmount  = "amount"
	colTransactionID      = "transaction_id"
	colTagID              = "tag_id"
	colRuleName           = "name"
	colRulPattern         = "pattern"
	colValue              = "value"
	colTag                = "tag"
	colRuleID             = "rule_id"
	colAmount             = "amount"
	colHash               = "hash"

	errUnableToWriteTag          = "unable to write tag: %w"
	errUnableToDeleteTag         = "unable to delete tag: %w"
	errUnableToDeleteRule        = "unable to delete rule: %w"
	errUnableToWriteRule         = "unable to write rule: %w"
	errUnableToDeleteTransaction = "unable to delete transaction: %w"
	errUnableToWriteTransaction  = "unable to write transaction: %w"
	errUnableToReadRule          = "unable to read rule: %w"
	errUnableToReadTag           = "unable to read tag: %w"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	selectRulesStmt = psql.Select("rules.*", "b.value").
			From(rulesTable).
			LeftJoin("(SELECT * FROM tags JOIN rules_tags as a on a.tag = tags.value) as b ON b.rule_id = rules.id")

	selectTagsStmt = psql.Select(colValue, colRuleID).From(tagsTable).
			InnerJoin("rules_tags on rules_tags.tag = tags.value")

	selectTransactionTagsStmt = psql.Select("*").From(transactionsTagsTable)

	insertTransaction = psql.Insert(transactionTable).
				Columns(
			colDate,
			colHash,
			colTransactionType,
			colTransactionContent,
			colTransactionAmount,
		)

	insertTransactionTag  = psql.Insert(transactionsTagsTable).Columns(colTransactionID, colTagID)
	selectTransactionStmp = psql.Select(colID, colDate, colTransactionType, colTransactionContent, colAmount, colTagID, colRuleID, colHash).
				From(transactionTable).
				LeftJoin("transactions_tags ON transactions_tags.transaction_id = transactions.id")

	insertRule = psql.Insert(rulesTable).Columns("id", "name", "pattern")
	updateRule = psql.Update(rulesTable)
	insertTag  = psql.Insert("tags").Columns("value")
)
