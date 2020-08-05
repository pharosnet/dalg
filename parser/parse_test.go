package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharosnet/dalg/parser"
)

func TestParseMySql(t *testing.T) {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		t.Error(pwdErr)
		return
	}

	dialect := "mysql"
	schemaPath := filepath.Join(pwd, "mysql/schema.sql")
	queryPath := filepath.Join(pwd, "mysql/query_select.sql")

	tables, queries, err := parser.Parse(dialect, schemaPath, queryPath)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("------")
	for _, table := range tables {
		t.Log("table", table.FullName, table.Schema, table.Name, table.GoName)
		t.Log("\t", "pk", table.PKs)
		for _, column := range table.Columns {
			t.Log("\t", "column", column.GoName, column.Name, column.Type, column.DefaultValue, column.GoType, column.Null, column.AutoIncrement)
		}
	}
	t.Log("------")
	for _, query := range queries {
		t.Log(query.Kind, query.Name)
		t.Log("\t", "selects")
		for _, expr := range query.SelectExprList.ExprList {

			t.Log("\t\t",
				"schema:", expr.Table.Schema,
				"table:", expr.Table.Table,
				"table as:", expr.Table.NameAs,
				"column from:", expr.ColumnQualifierName,
				"column name:", expr.ColumnName,
				"name:", expr.Name,
				"go type:", expr.GoType,
				"func:", expr.FuncName,
			)

		}
		t.Log("\t", "tables")
		for _, table := range query.TableList {
			t.Log("\t\t", table.Schema, table.Table, table.NameAs, table.Ref)
		}
		t.Log("\t", "conds")
		for _, expr := range query.CondExprList.ExprList {
			t.Log("\t\t",
				"schema:", expr.Table.Schema,
				"table:", expr.Table.Table,
				"table as:", expr.Table.NameAs,
				"column from:", expr.ColumnQualifierName,
				"column name:", expr.ColumnName,
				"name:", expr.Name,
				"PL:", expr.PlaceHolder,
				"go type:", expr.GoType,
			)
		}
	}
}
