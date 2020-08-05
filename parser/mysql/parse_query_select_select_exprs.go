package mysql

import (
	"fmt"
	"reflect"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

func parseSelectExprs(query *entry.Query, exprs sqlparser.SelectExprs) (err error) {

	err = exprs.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {

		switch node.(type) {
		case *sqlparser.StarExpr:
			err = fmt.Errorf("parse query exprs failed, star(*) expr is not support, \n%s", query.Sql)
		case *sqlparser.NonStarExpr:
			err = parseSelectNonStarExpr(query, node.(*sqlparser.NonStarExpr))
		default:
			err = fmt.Errorf("parse query exprs failed, %v is not support, \n%s", reflect.TypeOf(node), query.Sql)
		}

		return
	})

	return
}

func parseSelectNonStarExpr(query *entry.Query, stmt *sqlparser.NonStarExpr) (err error) {
	//
	queryExpr := &entry.QueryExpr{
		Table:               entry.QueryTable{},
		ColumnQualifierName: "",
		ColumnName:          "",
		FuncName:            "",
		Name:                "",
		GoType:              nil,
	}

	// as -> name
	if !stmt.As.IsEmpty() {
		queryExpr.Name = stmt.As.CompliantName()
	}

	node := stmt.Expr
	switch node.(type) {
	case *sqlparser.ColName:
		if queryExpr.ColumnName == "" {
			x := node.(*sqlparser.ColName)
			queryExpr.ColumnName = x.Name.CompliantName()
			if x.Qualifier != nil {
				queryExpr.ColumnQualifierName = x.Qualifier.Name.CompliantName()
			}
		}
	case *sqlparser.BinaryExpr: // + - * / & | 
		x := node.(*sqlparser.BinaryExpr)
		if !(x.Operator == "+" || x.Operator == "-" || x.Operator == "*" || x.Operator == "/" || x.Operator == "%") {
			err = fmt.Errorf("parse non star expr failed, %s is not support", x.Operator)
			return
		}
		left, ok := x.Left.(*sqlparser.ColName)
		if !ok {
			err = fmt.Errorf("parse non star expr failed, found binary expr, but left is not a column, left is %v", x.Left)
			return
		}
		if queryExpr.ColumnName == "" {
			queryExpr.ColumnName = left.Name.CompliantName()
			if left.Qualifier != nil {
				queryExpr.ColumnQualifierName = left.Qualifier.Name.CompliantName()
			}
		}
	//case *sqlparser.UnaryExpr: //

	case *sqlparser.ParenExpr:
		x := node.(*sqlparser.ParenExpr)
		_, ok := x.Expr.(*sqlparser.ComparisonExpr)
		if !ok {
			err = fmt.Errorf("parse non star expr failed, found paren expr, but it has not a comparion expr, %v", reflect.TypeOf(x.Expr))
			return
		}
		queryExpr.GoType = entry.NewGoType("bool")
	case *sqlparser.ComparisonExpr: // > < !=
		x := node.(*sqlparser.ComparisonExpr)
		if !(x.Operator == ">" || x.Operator == ">=" || x.Operator == "<" || x.Operator == "<=" || x.Operator == "=" || x.Operator == "!=") {
			err = fmt.Errorf("parse non star expr failed, %s is not support", x.Operator)
			return
		}
		queryExpr.GoType = entry.NewGoType("bool")
		//switch x.Left.(type) {
		//case *sqlparser.ColName:
		//	left := x.Left.(*sqlparser.ColName)
		//	if queryExpr.ColumnName == "" {
		//		queryExpr.ColumnName = left.Name.CompliantName()
		//		if left.Qualifier != nil {
		//			queryExpr.ColumnQualifierName = left.Qualifier.Name.CompliantName()
		//		}
		//	}
		//case *sqlparser.FuncExpr:
		//	fn := x.Left.(*sqlparser.FuncExpr)
		//	if !fn.IsAggregate() {
		//		err = fmt.Errorf("parse non star expr failed, found comparison expr, but left is not a aggregate func, left is %v", x.Left)
		//		return
		//	}
		//	if !(fn.Name.Lowered() == "count" || fn.Name.Lowered() == "sum" || fn.Name.Lowered() == "max" || fn.Name.Lowered() == "min" ) {
		//		err = fmt.Errorf("parse non star expr failed, found comparison expr, but left is not a supported aggregate func, left is %v", x.Left)
		//		return
		//	}
		//	if fn.Name.Lowered() == "count" {
		//		queryExpr.GoType = entry.NewGoType("int64")
		//	}
		//	err = fn.Exprs.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		//		fnX, ok := node.(*sqlparser.NonStarExpr)
		//		if !ok {
		//			err = fmt.Errorf("parse non star expr failed, found comparison expr, left func, but func expr is not non star expr %v", reflect.TypeOf(node))
		//			return
		//		}
		//		fnCol, ok := fnX.Expr.(*sqlparser.ColName)
		//		if !ok {
		//			err = fmt.Errorf("parse non star expr failed, found comparison expr, left func, but func expr is column expr %v", reflect.TypeOf(fnX.Expr))
		//			return
		//		}
		//		if queryExpr.ColumnName == "" {
		//			queryExpr.ColumnName = fnCol.Name.CompliantName()
		//			if fnCol.Qualifier != nil {
		//				queryExpr.ColumnQualifierName = fnCol.Qualifier.Name.CompliantName()
		//			}
		//		}
		//		return
		//	})
		//	if err != nil {
		//		return
		//	}
		//default:
		//	err = fmt.Errorf("parse non star expr failed, found comparison expr, but left is not a column or func, left is %v", x.Left)
		//	return
		//}
		//left, ok := x.Left.(*sqlparser.ColName)
		//if !ok {
		//	err = fmt.Errorf("parse non star expr failed, found binary expr, but left is not a column, left is %v", x.Left)
		//	return
		//}
		//if queryExpr.ColumnName == "" {
		//	queryExpr.ColumnName = left.Name.CompliantName()
		//	if left.Qualifier != nil {
		//		queryExpr.ColumnQualifierName = left.Qualifier.Name.CompliantName()
		//	}
		//}
	case *sqlparser.FuncExpr:
		fn := node.(*sqlparser.FuncExpr)
		if !fn.IsAggregate() {
			err = fmt.Errorf("parse non star expr failed, found func expr, but it is not a aggregate func, it is %v", node)
			return
		}
		if !(fn.Name.Lowered() == "count" || fn.Name.Lowered() == "sum" || fn.Name.Lowered() == "max" || fn.Name.Lowered() == "min") {
			err = fmt.Errorf("parse non star expr failed, found func expr, but it is not a supported aggregate func, it is %v", node)
			return
		}
		if fn.Name.Lowered() == "count" {
			queryExpr.GoType = entry.NewGoType("int")
		}
		queryExpr.FuncName = fn.Name.Lowered()
		err = fn.Exprs.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
			fnX, ok := node.(*sqlparser.NonStarExpr)
			if !ok {
				err = fmt.Errorf("parse non star expr failed, found func expr, but func expr is not non star expr %v", reflect.TypeOf(node))
				return
			}
			fnCol, ok := fnX.Expr.(*sqlparser.ColName)
			if !ok {
				err = fmt.Errorf("parse non star expr failed, found func expr, but func expr is column expr %v", reflect.TypeOf(fnX.Expr))
				return
			}
			if queryExpr.ColumnName == "" {
				queryExpr.ColumnName = fnCol.Name.CompliantName()
				if fnCol.Qualifier != nil {
					queryExpr.ColumnQualifierName = fnCol.Qualifier.Name.CompliantName()
				}

			}
			return
		})
		if err != nil {
			return
		}
	case *sqlparser.Subquery:
		x := node.(*sqlparser.Subquery)
		subQuery := entry.NewQuery()
		parseSubQueryErr := parseQuerySelect(subQuery, x.Select)
		if parseSubQueryErr != nil {
			err = fmt.Errorf("parse non star expr failed, found sub query expr, %v", parseSubQueryErr)
			return
		}
		subQuery.Fill()
		if len(subQuery.SelectExprList.ExprList) != 1 {
			err = fmt.Errorf("parse non star expr failed, sub query must has one expr, but found %d", len(subQuery.SelectExprList.ExprList))
			return
		}
		subExpr := subQuery.SelectExprList.ExprList[0]
		queryExpr.Table = subExpr.Table
		queryExpr.ColumnName = subExpr.ColumnName
		queryExpr.FuncName = subExpr.FuncName
		queryExpr.GoType = subExpr.GoType
		entry.QueryMergeCond(query, subQuery)
	case *sqlparser.ExistsExpr:
		x := node.(*sqlparser.ExistsExpr)
		subQuery := entry.NewQuery()
		parseSubQueryErr := parseQuerySelect(subQuery, x.Subquery.Select)
		if parseSubQueryErr != nil {
			err = fmt.Errorf("parse non star expr failed, found exists expr, %v", parseSubQueryErr)
			return
		}
		subQuery.Fill()
		queryExpr.GoType = entry.NewGoType("bool")
		entry.QueryMergeCond(query, subQuery)
	default:
		err = fmt.Errorf("parse non star expr failed, %v is not support", reflect.TypeOf(node))
	}

	queryExpr.BuildName()
	query.SelectExprList.ExprList = append(query.SelectExprList.ExprList, queryExpr)

	return
}
