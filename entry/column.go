package entry

import (
	"strings"
)

type ColumnType string

type Column struct {
	Name          string
	Type          ColumnType
	Null          bool
	GoName        string
	GoType        *GoType
	DefaultValue  string
	AutoIncrement bool
}

func (c *Column) MappingGoType() (ok bool) {
	for _, mapping := range ColumnTypeMappings {
		if string(mapping.ColumnType) == string(c.Type) && mapping.NullAble == c.Null {
			c.GoType = mapping.GoType
			ok = true
			return
		}
	}
	return
}

type ColumnTypeMapping struct {
	ColumnType ColumnType
	GoType     *GoType
	NullAble   bool
}

var ColumnTypeMappings []*ColumnTypeMapping

func OverrideColumnTypeMappings(columnType string, nullable bool, goType string) {
	columnType = strings.TrimSpace(strings.ToUpper(columnType))
	override := false
	for _, mapping := range ColumnTypeMappings {
		if string(mapping.ColumnType) == columnType && mapping.NullAble == nullable {
			mapping.GoType = NewGoType(goType)
			override = true
			break
		}
	}
	if !override {
		ColumnTypeMappings = append(ColumnTypeMappings, &ColumnTypeMapping{
			ColumnType: ColumnType(columnType),
			GoType:     NewGoType(goType),
			NullAble:   nullable,
		})
	}
}

func init() {
	ColumnTypeMappings = make([]*ColumnTypeMapping, 0, 1)
	ColumnTypeMappings = append(ColumnTypeMappings,
		&ColumnTypeMapping{
			ColumnType: TINYINT,
			GoType:     NewGoType("database/sql.NullInt32"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TINYINT,
			GoType:     NewGoType("int32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: SMALLINT,
			GoType:     NewGoType("database/sql.NullInt32"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: SMALLINT,
			GoType:     NewGoType("int32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: MEDIUMINT,
			GoType:     NewGoType("database/sql.NullInt32"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: MEDIUMINT,
			GoType:     NewGoType("int32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: INT,
			GoType:     NewGoType("database/sql.NullInt32"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: INT,
			GoType:     NewGoType("int32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: INTEGER,
			GoType:     NewGoType("database/sql.NullInt32"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: INTEGER,
			GoType:     NewGoType("int32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: BIGINT,
			GoType:     NewGoType("database/sql.NullInt64"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: BIGINT,
			GoType:     NewGoType("int64"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: FLOAT,
			GoType:     NewGoType("database/sql.NullFloat64"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: FLOAT,
			GoType:     NewGoType("float32"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: DOUBLE,
			GoType:     NewGoType("database/sql.NullFloat64"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: DOUBLE,
			GoType:     NewGoType("float64"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: DECIMAL,
			GoType:     NewGoType("database/sql.NullFloat64"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: DECIMAL,
			GoType:     NewGoType("float64"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: BOOLEAN,
			GoType:     NewGoType("database/sql.NullBool"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: BOOLEAN,
			GoType:     NewGoType("bool"),
			NullAble:   false,
		},
		// DATE, TIME, YEAR, DATETIME, c
		&ColumnTypeMapping{
			ColumnType: DATE,
			GoType:     NewGoType("database/sql.NullTime"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: DATE,
			GoType:     NewGoType("time.Time"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: TIME,
			GoType:     NewGoType("database/sql.NullTime"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TIME,
			GoType:     NewGoType("time.Time"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: YEAR,
			GoType:     NewGoType("database/sql.NullTime"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: YEAR,
			GoType:     NewGoType("time.Time"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: DATETIME,
			GoType:     NewGoType("database/sql.NullTime"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: DATETIME,
			GoType:     NewGoType("time.Time"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: TIMESTAMP,
			GoType:     NewGoType("database/sql.NullTime"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TIMESTAMP,
			GoType:     NewGoType("time.Time"),
			NullAble:   false,
		},
		// CHAR, NCHAR, VARCHAR, NVARCHAR
		&ColumnTypeMapping{
			ColumnType: CHAR,
			GoType:     NewGoType("database/sql.NullString"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: CHAR,
			GoType:     NewGoType("string"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: NCHAR,
			GoType:     NewGoType("database/sql.NullString"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: NCHAR,
			GoType:     NewGoType("string"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: VARCHAR,
			GoType:     NewGoType("database/sql.NullString"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: VARCHAR,
			GoType:     NewGoType("string"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: NVARCHAR,
			GoType:     NewGoType("database/sql.NullString"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: NVARCHAR,
			GoType:     NewGoType("string"),
			NullAble:   false,
		},
		// TINYBLOB, TINYTEXT, BLOB, TEXT, MEDIUMBLOB, MEDIUMTEXT, LONGBLOB, LONGTEXT
		&ColumnTypeMapping{
			ColumnType: TINYBLOB,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TINYBLOB,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: TINYTEXT,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TINYTEXT,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: BLOB,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: BLOB,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: TEXT,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TEXT,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: MEDIUMBLOB,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: TEXT,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: MEDIUMTEXT,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: MEDIUMTEXT,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: LONGBLOB,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: LONGBLOB,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: LONGTEXT,
			GoType:     NewGoType("github.com/pharosnet/dalc.NullBytes"),
			NullAble:   true,
		},
		&ColumnTypeMapping{
			ColumnType: LONGTEXT,
			GoType:     NewGoType("[]byte"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: JSON,
			GoType:     NewGoType("json.RawMessage"),
			NullAble:   false,
		},
		&ColumnTypeMapping{
			ColumnType: "",
			GoType:     NewGoType("github.com/pharosnet/dalc.NullJson"),
			NullAble:   true,
		},
	)
}
