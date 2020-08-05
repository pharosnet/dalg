package mysql

import (
	"fmt"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseQueryUpdate(query *entry.Query, stmt *sqlparser.Update) (err error) {

	nameAs := ""
	if !stmt.Table.As.IsEmpty() {
		nameAs = stmt.Table.As.CompliantName()
	}
	schema := ""
	tableName := ""

	err = stmt.Table.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type) {
		case *sqlparser.TableName:
			table := node.(*sqlparser.TableName)
			tableName = table.Name.CompliantName()
			if !table.Qualifier.IsEmpty() {
				schema = table.Qualifier.CompliantName()
			}
		}
		return
	})

	if err != nil {
		return
	}

	query.TableList = append(query.TableList, &entry.QueryTable{
		Schema: schema,
		Table:  tableName,
		NameAs: nameAs,
	})

	for _, expr := range stmt.Exprs {
		columnName := expr.Name.Name.CompliantName()
		isArg := false
		argExpr := expr.Expr
		switch argExpr.(type) {
		case *sqlparser.SQLVal:
			val := argExpr.(*sqlparser.SQLVal)
			if val.Type == sqlparser.ValArg {
				isArg = true
			}
		}
		if isArg {
			query.SelectExprList.ExprList = append(query.SelectExprList.ExprList, &entry.QueryExpr{
				Table: entry.QueryTable{
					Schema: schema,
					Table:  tableName,
					NameAs: nameAs,
				},
				ColumnQualifierName: "",
				ColumnName:          columnName,
				FuncName:            "",
				Name:                "",
				GoType:              nil,
			})
		}
	}

	// where
	err = parseQueryWhere(query, stmt.Where)
	if err != nil {
		return
	}

	if len(query.CondExprList.ExprList) < 1 {
		err = fmt.Errorf("parse delete failed, found no condm in \n%s", query.Sql)
		return
	}

	return
}
