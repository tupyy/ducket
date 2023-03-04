package models

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	"github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


Table: transactions_labels
[ 0] transaction_id                                 INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] label_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []


JSON Sample
-------------------------------------
{    "transaction_id": 19,    "label_id": 3}



*/

// TransactionsLabels struct is a row record of the transactions_labels table in the finance database
type TransactionsLabels struct {
	//[ 0] transaction_id                                 INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	TransactionID int32 `gorm:"primary_key;column:transaction_id;type:INT4;"`
	//[ 1] label_id                                       INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	LabelID int32 `gorm:"primary_key;column:label_id;type:INT4;"`
}

var transactions_labelsTableInfo = &TableInfo{
	Name: "transactions_labels",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "transaction_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "TransactionID",
			GoFieldType:        "int32",
			JSONFieldName:      "transaction_id",
			ProtobufFieldName:  "transaction_id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "label_id",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "INT4",
			DatabaseTypePretty: "INT4",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "INT4",
			ColumnLength:       -1,
			GoFieldName:        "LabelID",
			GoFieldType:        "int32",
			JSONFieldName:      "label_id",
			ProtobufFieldName:  "label_id",
			ProtobufType:       "int32",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *TransactionsLabels) TableName() string {
	return "transactions_labels"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *TransactionsLabels) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *TransactionsLabels) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *TransactionsLabels) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *TransactionsLabels) TableInfo() *TableInfo {
	return transactions_labelsTableInfo
}
