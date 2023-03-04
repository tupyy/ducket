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


Table: label
[ 0] key                                            VARCHAR(30)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 30      default: []
[ 1] value                                          VARCHAR(50)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 50      default: []


JSON Sample
-------------------------------------
{    "key": "FDUJBHRtFpCwInIOvKHFporHy",    "value": "EBXaSKHTFnOxKMFoHPFnYwnPZ"}



*/

// Label struct is a row record of the label table in the finance database
type Label struct {
	//[ 0] key                                            VARCHAR(30)          null: false  primary: true   isArray: false  auto: false  col: VARCHAR         len: 30      default: []
	Key string `gorm:"primary_key;column:key;type:VARCHAR;size:30;"`
	//[ 1] value                                          VARCHAR(50)          null: false  primary: false  isArray: false  auto: false  col: VARCHAR         len: 50      default: []
	Value string `gorm:"column:value;type:VARCHAR;size:50;"`
}

var labelTableInfo = &TableInfo{
	Name: "label",
	Columns: []*ColumnInfo{

		&ColumnInfo{
			Index:              0,
			Name:               "key",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(30)",
			IsPrimaryKey:       true,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       30,
			GoFieldName:        "Key",
			GoFieldType:        "string",
			JSONFieldName:      "key",
			ProtobufFieldName:  "key",
			ProtobufType:       "string",
			ProtobufPos:        1,
		},

		&ColumnInfo{
			Index:              1,
			Name:               "value",
			Comment:            ``,
			Notes:              ``,
			Nullable:           false,
			DatabaseTypeName:   "VARCHAR",
			DatabaseTypePretty: "VARCHAR(50)",
			IsPrimaryKey:       false,
			IsAutoIncrement:    false,
			IsArray:            false,
			ColumnType:         "VARCHAR",
			ColumnLength:       50,
			GoFieldName:        "Value",
			GoFieldType:        "string",
			JSONFieldName:      "value",
			ProtobufFieldName:  "value",
			ProtobufType:       "string",
			ProtobufPos:        2,
		},
	},
}

// TableName sets the insert table name for this struct type
func (l *Label) TableName() string {
	return "label"
}

// BeforeSave invoked before saving, return an error if field is not populated.
func (l *Label) BeforeSave() error {
	return nil
}

// Prepare invoked before saving, can be used to populate fields etc.
func (l *Label) Prepare() {
}

// Validate invoked before performing action, return an error if field is not populated.
func (l *Label) Validate(action Action) error {
	return nil
}

// TableInfo return table meta data
func (l *Label) TableInfo() *TableInfo {
	return labelTableInfo
}
