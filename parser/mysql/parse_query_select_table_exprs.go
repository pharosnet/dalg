package mysql

import (
	"fmt"
	"reflect"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseTableExprs(query *entry.Query, exprs sqlparser.TableExprs) (err error) {

	err = exprs.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		err = parseTableExpr(query, node.(sqlparser.TableExpr))
		return
	})

	return
}

func parseTableExpr(query *entry.Query, expr sqlparser.TableExpr) (err error) {
	switch expr.(type) {
	case *sqlparser.AliasedTableExpr:
		t := expr.(*sqlparser.AliasedTableExpr)
		switch t.Expr.(type) {
		case *sqlparser.TableName:
			table := t.Expr.(*sqlparser.TableName)
			queryTable := &entry.QueryTable{}
			queryTable.Table = table.Name.CompliantName()
			if !table.Qualifier.IsEmpty() {
				queryTable.Schema = table.Qualifier.CompliantName()
			}
			if !t.As.IsEmpty() {
				queryTable.NameAs = t.As.CompliantName()
			}
			query.TableList = append(query.TableList, queryTable)
		case *sqlparser.Subquery:
			err = fmt.Errorf("parse table exprs failed, %v is not support, \n%s", reflect.TypeOf(t), query.Sql)
		}
	case *sqlparser.JoinTableExpr:
		t := expr.(*sqlparser.JoinTableExpr)
		err = parseTableExpr(query, t.LeftExpr)
		if err != nil {
			return
		}
		err = parseTableExpr(query, t.RightExpr)
		if err != nil {
			return
		}
	default:
		err = fmt.Errorf("parse table exprs failed, %v is not support, \n%s", reflect.TypeOf(expr), query.Sql)
	}
	return
}
