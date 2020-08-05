package entry

import (
	"fmt"
	"strings"
)

const (
	TINYINT   = ColumnType("TINYINT")
	SMALLINT  = ColumnType("SMALLINT")
	MEDIUMINT = ColumnType("MEDIUMINT")
	INT       = ColumnType("INT")
	INTEGER   = ColumnType("INTEGER")
	BIGINT    = ColumnType("BIGINT")
	FLOAT     = ColumnType("FLOAT")
	DOUBLE    = ColumnType("DOUBLE")
	DECIMAL   = ColumnType("DECIMAL")
	BOOLEAN   = ColumnType("BOOLEAN")

	DATE      = ColumnType("DATE")
	TIME      = ColumnType("TIME")
	YEAR      = ColumnType("YEAR")
	DATETIME  = ColumnType("DATETIME")
	TIMESTAMP = ColumnType("TIMESTAMP")

	CHAR       = ColumnType("CHAR")
	NCHAR      = ColumnType("NCHAR")
	VARCHAR    = ColumnType("VARCHAR")
	NVARCHAR   = ColumnType("NVARCHAR")
	TINYBLOB   = ColumnType("TINYBLOB")
	TINYTEXT   = ColumnType("TINYTEXT")
	BLOB       = ColumnType("BLOB")
	TEXT       = ColumnType("TEXT")
	MEDIUMBLOB = ColumnType("MEDIUMBLOB")
	MEDIUMTEXT = ColumnType("MEDIUMTEXT")
	LONGBLOB   = ColumnType("LONGBLOB")
	LONGTEXT   = ColumnType("LONGTEXT")
	JSON       = ColumnType("JSON")
)

var columnTypes map[string]ColumnType

func NewColumnType(v string) (columnType ColumnType, err error) {
	v = strings.ToUpper(strings.TrimSpace(v))
	leftBracketIndex := strings.Index(v, "(")
	if leftBracketIndex >= 0 {
		v = v[:leftBracketIndex]
	}
	has := false
	columnType, has = columnTypes[v]
	if !has {
		err = fmt.Errorf("unknown column type %s", v)
	}
	return
}

func init() {
	columnTypes = make(map[string]ColumnType)
	columnTypes["TINYINT"] = TINYINT
	columnTypes["SMALLINT"] = SMALLINT
	columnTypes["MEDIUMINT"] = MEDIUMINT
	columnTypes["INT"] = INT
	columnTypes["INTEGER"] = INTEGER
	columnTypes["BIGINT"] = BIGINT
	columnTypes["FLOAT"] = FLOAT
	columnTypes["DOUBLE"] = DOUBLE
	columnTypes["DECIMAL"] = DECIMAL
	columnTypes["BOOLEAN"] = BOOLEAN

	columnTypes["DATE"] = DATE
	columnTypes["TIME"] = TIME
	columnTypes["YEAR"] = YEAR
	columnTypes["DATETIME"] = DATETIME
	columnTypes["TIMESTAMP"] = TIMESTAMP

	columnTypes["CHAR"] = CHAR
	columnTypes["NCHAR"] = NCHAR
	columnTypes["VARCHAR"] = VARCHAR
	columnTypes["NVARCHAR"] = NVARCHAR
	columnTypes["TINYBLOB"] = TINYBLOB
	columnTypes["TINYTEXT"] = TINYTEXT
	columnTypes["BLOB"] = BLOB
	columnTypes["TEXT"] = TEXT
	columnTypes["MEDIUMBLOB"] = MEDIUMBLOB
	columnTypes["MEDIUMTEXT"] = MEDIUMTEXT
	columnTypes["LONGBLOB"] = LONGBLOB
	columnTypes["LONGTEXT"] = LONGTEXT
	columnTypes["JSON"] = JSON
}
