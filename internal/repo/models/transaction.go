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


Table: transaction
[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
[ 1] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [now() AT TIME ZONE 'UTC']
[ 2] date                                           TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
[ 3] transaction_type                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 4] description                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 5] recipient                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 6] amount                                         NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []


JSON Sample
-------------------------------------
{    "id": 37,    "created_at": "2094-10-09T08:56:50.415676689+02:00",    "date": "2075-12-07T09:48:10.029999623+01:00",    "transaction_type": "yasZXsNiQpSJuHwWZxIVfGNpW",    "description": "EsmRFCWSNXJmKQhOckqyrrWMc",    "recipient": "WTXaTtqYNiKOopTITZfUhXkLw",    "amount": 0.7233102408856489}



*/

// Transaction struct is a row record of the transaction table in the finance database
type Transaction struct {
	//[ 0] id                                             INT4                 null: false  primary: true   isArray: false  auto: false  col: INT4            len: -1      default: []
	ID int32 `gorm:"primary_key;column:id;type:INT4;"`
	//[ 1] created_at                                     TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [now() AT TIME ZONE 'UTC']
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;"`
	//[ 2] date                                           TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: []
	Date time.Time `gorm:"column:date;type:TIMESTAMP;"`
	//[ 3] transaction_type                               TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	TransactionType string `gorm:"column:transaction_type;type:TEXT;"`
	//[ 4] description                                    TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Description string `gorm:"column:description;type:TEXT;"`
	//[ 5] recipient                                      TEXT                 null: false  primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
	Recipient string `gorm:"column:recipient;type:TEXT;"`
	//[ 6] amount                                         NUMERIC              null: false  primary: false  isArray: false  auto: false  col: NUMERIC         len: -1      default: []
	Amount float64 `gorm:"column:amount;type:NUMERIC;"`
}

var transactionTableInfo = &TableInfo{
	Name: "transaction",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "id",
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
			GoFieldName:        "ID",
			GoFieldType:        "int32",
			JSONFieldName:      "id",
			ProtobufFieldName:  "id",
			ProtobufType:       "int32",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "created_at",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "CreatedAt",
			GoFieldType:        "time.Time",
			JSONFieldName:      "created_at",
			ProtobufFieldName:  "created_at",
			ProtobufType:       "uint64",
			ProtobufPos:        2,
		},

		&ColumnInfo{
			Index:              2,
			Name:               "date",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TIMESTAMP",
			DatabaseTypePretty: "TIMESTAMP",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TIMESTAMP",
			ColumnLength:       -1,
			GoFieldName:        "Date",
			GoFieldType:        "time.Time",
			JSONFieldName:      "date",
			ProtobufFieldName:  "date",
			ProtobufType:       "uint64",
			ProtobufPos:        3,
		},

		&ColumnInfo{
			Index:              3,
			Name:               "transaction_type",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "TransactionType",
			GoFieldType:        "string",
			JSONFieldName:      "transaction_type",
			ProtobufFieldName:  "transaction_type",
			ProtobufType:       "string",
			ProtobufPos:        4,
		},

		&ColumnInfo{
			Index:              4,
			Name:               "description",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Description",
			GoFieldType:        "string",
			JSONFieldName:      "description",
			ProtobufFieldName:  "description",
			ProtobufType:       "string",
			ProtobufPos:        5,
		},

		&ColumnInfo{
			Index:              5,
			Name:               "recipient",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "TEXT",
			DatabaseTypePretty: "TEXT",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "TEXT",
			ColumnLength:       -1,
			GoFieldName:        "Recipient",
			GoFieldType:        "string",
			JSONFieldName:      "recipient",
			ProtobufFieldName:  "recipient",
			ProtobufType:       "string",
			ProtobufPos:        6,
		},

		&ColumnInfo{
			Index:              6,
			Name:               "amount",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "NUMERIC",
			DatabaseTypePretty: "NUMERIC",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "NUMERIC",
			ColumnLength:       -1,
			GoFieldName:        "Amount",
			GoFieldType:        "float64",
			JSONFieldName:      "amount",
			ProtobufFieldName:  "amount",
			ProtobufType:       "float",
			ProtobufPos:        7,
		},
	},
}

// TableName sets the insert table name for this struct type
func (t *Transaction) TableName() string {
	return "transaction"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (t *Transaction) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (t *Transaction) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (t *Transaction) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (t *Transaction) TableInfo() *TableInfo {
	return transactionTableInfo
}
