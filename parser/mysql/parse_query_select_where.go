package mysql

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseQueryWhere(query *entry.Query, node *sqlparser.Where) (err error) {
	if node == nil {
		return
	}
	if node.Type == sqlparser.HavingStr {
		err = fmt.Errorf("parse where failed, having is not support")
		return
	}

	conds := make([]*entry.CondExpr, 0, 1)
	err = node.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type) {
		case *sqlparser.ExistsExpr:
			// todo fill
		case *sqlparser.ColName:
			col := node.(*sqlparser.ColName)
			cond := &entry.CondExpr{
				ColumnName: col.Name.CompliantName(),
			}
			if col.Qualifier != nil {
				cond.ColumnQualifierName = col.Qualifier.Name.CompliantName()
			}
			conds = append(conds, cond)
		case *sqlparser.SQLVal:
			val := node.(*sqlparser.SQLVal)
			cond := conds[len(conds)-1]
			valBytes := bytes.TrimSpace(val.Val)

			if val.Type == sqlparser.ValArg {
				cond.IsArg = true
			} else if bytes.IndexByte(valBytes, '#') == 0 && bytes.LastIndexByte(valBytes, '#') == len(val.Val)-1 {
				cond.IsArg = true
				cond.PlaceHolder = string(valBytes)
			}
		case *sqlparser.Subquery:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.FuncExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.CaseExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.ValuesFuncExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.ConvertExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.MatchExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		case *sqlparser.GroupConcatExpr:
			err = fmt.Errorf("parse where failed, %v is not support", reflect.TypeOf(node))
		}
		kontinue = true
		return
	})

	if len(conds) == 0 {
		return
	}
	for _, cond := range conds {
		if cond.IsArg {
			query.CondExprList.ExprList = append(query.CondExprList.ExprList, cond)
		}
	}

	return
}
