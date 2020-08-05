package mysql

import (
	"fmt"
	"reflect"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseQueryInsert(query *entry.Query, stmt *sqlparser.Insert) (err error) {

	query.TableList = append(query.TableList, &entry.QueryTable{
		Schema: stmt.Table.Qualifier.CompliantName(),
		Table:  stmt.Table.Name.CompliantName(),
		NameAs: "",
	})

	for _, column := range stmt.Columns {
		query.SelectExprList.ExprList = append(query.SelectExprList.ExprList, &entry.QueryExpr{
			Table: entry.QueryTable{
				Schema: stmt.Table.Qualifier.CompliantName(),
				Table:  stmt.Table.Name.CompliantName(),
				NameAs: "",
			},
			ColumnQualifierName: "",
			ColumnName:          column.CompliantName(),
			FuncName:            "",
			Name:                "",
			GoType:              nil,
		})
	}

	switch stmt.Rows.(type) {
	case *sqlparser.Select:
		err = fmt.Errorf("parse query insert failed, %s is not support, in \n%s", reflect.TypeOf(stmt.Rows), query.Sql)
	case *sqlparser.Union:
		err = fmt.Errorf("parse query insert failed, %s is not support, in \n%s", reflect.TypeOf(stmt.Rows), query.Sql)
	case sqlparser.Values:
		values := stmt.Rows.(sqlparser.Values)
		for _, value := range values {
			err = value.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
				switch node.(type) {
				case sqlparser.Exprs:
					idx := -1
					exprs := node.(sqlparser.Exprs)
					err = exprs.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
						switch node.(type) {
						case *sqlparser.SQLVal:
							idx++
							val := node.(*sqlparser.SQLVal)
							if val.Type == sqlparser.ValArg {
								query.CondExprList.ExprList = append(query.CondExprList.ExprList, &entry.CondExpr{
									Table: entry.QueryTable{
										Schema: stmt.Table.Qualifier.CompliantName(),
										Table:  stmt.Table.Name.CompliantName(),
										NameAs: "",
									},
									ColumnQualifierName: "",
									ColumnName:          query.SelectExprList.ExprList[idx].ColumnName,
									Name:                "",
									GoType:              nil,
								})
							}
						}
						return
					})
				default:
					err = fmt.Errorf("parse query insert failed, %s is not support, in \n%s", reflect.TypeOf(node), query.Sql)
				}
				return
			})
		}
	case *sqlparser.ParenSelect:
		err = fmt.Errorf("parse query insert failed, %s is not support, in \n%s", reflect.TypeOf(stmt.Rows), query.Sql)
	}

	if err != nil {
		return
	}

	if len(query.CondExprList.ExprList) <= 0 {
		err = fmt.Errorf("parse query insert failed, found zero insert rows, in \n%s", query.Sql)
		return
	}
	if len(query.CondExprList.ExprList) > len(query.SelectExprList.ExprList) {
		err = fmt.Errorf("parse query insert failed, found insert rows over columns, in \n%s", query.Sql)
		return
	}

	return
}
