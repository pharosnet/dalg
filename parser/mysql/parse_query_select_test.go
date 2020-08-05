package mysql_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pharosnet/dalg/parser/mysql"
)

func TestParseMySQLQuerySelect(t *testing.T) {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		t.Error(pwdErr)
		return
	}

	schemaPath := filepath.Join(pwd, "query_select.sql")

	p, pErr := ioutil.ReadFile(schemaPath)
	if pErr != nil {
		t.Error(pErr)
		return
	}

	queries, parseErr := mysql.ParseMySQLQuery(string(p))
	if parseErr != nil {
		t.Error(parseErr)
		return
	}

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
			t.Log("\t\t", table.Schema, table.Table, table.NameAs)
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
			)
		}
	}

}
