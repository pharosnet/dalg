package mysql

import (
	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseQueryLimit(query *entry.Query, node *sqlparser.Limit) (err error) {
	if node == nil {
		return
	}
	switch node.Offset.(type) {
	case *sqlparser.SQLVal:
		val := node.Offset.(*sqlparser.SQLVal)
		if val.Type == sqlparser.ValArg {
			query.CondExprList.ExprList = append(query.CondExprList.ExprList, &entry.CondExpr{
				Table:               entry.QueryTable{},
				ColumnQualifierName: "",
				ColumnName:          "",
				PlaceHolder:         "",
				Args:                nil,
				Name:                "Offset",
				GoType:              entry.NewGoType("int"),
				IsArg:               true,
			})
		}
	}
	switch node.Rowcount.(type) {
	case *sqlparser.SQLVal:
		val := node.Rowcount.(*sqlparser.SQLVal)
		if val.Type == sqlparser.ValArg {
			query.CondExprList.ExprList = append(query.CondExprList.ExprList, &entry.CondExpr{
				Table:               entry.QueryTable{},
				ColumnQualifierName: "",
				ColumnName:          "",
				PlaceHolder:         "",
				Args:                nil,
				Name:                "Limit",
				GoType:              entry.NewGoType("int"),
				IsArg:               true,
			})
		}
	}
	return
}
