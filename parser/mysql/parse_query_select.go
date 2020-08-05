package mysql

import (
	"fmt"
	"reflect"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseQuerySelect(query *entry.Query, stmt sqlparser.SelectStatement) (err error) {

	err = stmt.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node.(type) {
		case sqlparser.Comments:

		case sqlparser.SelectExprs:
			err = parseSelectExprs(query, node.(sqlparser.SelectExprs))
		case sqlparser.TableExprs:
			err = parseTableExprs(query, node.(sqlparser.TableExprs))
		case *sqlparser.Where:
			err = parseQueryWhere(query, node.(*sqlparser.Where))
		case *sqlparser.Limit:
			err = parseQueryLimit(query, node.(*sqlparser.Limit))
		case sqlparser.GroupBy:

		case sqlparser.OrderBy:

		default:
			err = fmt.Errorf("parse query select failed, %s is not support, in \n%s", reflect.TypeOf(node), query.Sql)
		}
		return
	})

	return
}
