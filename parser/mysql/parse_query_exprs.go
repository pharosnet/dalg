package mysql

import (
	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
)

/*
func (*AndExpr) iExpr()         {}
func (*OrExpr) iExpr()          {}
func (*NotExpr) iExpr()         {}
func (*ParenExpr) iExpr()       {}
func (*ComparisonExpr) iExpr()  {}
func (*RangeCond) iExpr()       {}
func (*IsExpr) iExpr()          {}
func (*ExistsExpr) iExpr()      {}
func (*SQLVal) iExpr()          {}
func (*NullVal) iExpr()         {}
func (BoolVal) iExpr()          {}
func (*ColName) iExpr()         {}
func (ValTuple) iExpr()         {}
func (*Subquery) iExpr()        {}
func (ListArg) iExpr()          {}
func (*BinaryExpr) iExpr()      {}
func (*UnaryExpr) iExpr()       {}
func (*IntervalExpr) iExpr()    {}
func (*CollateExpr) iExpr()     {}
func (*FuncExpr) iExpr()        {}
func (*CaseExpr) iExpr()        {}
func (*ValuesFuncExpr) iExpr()  {}
func (*ConvertExpr) iExpr()     {}
func (*MatchExpr) iExpr()       {}
func (*GroupConcatExpr) iExpr() {}
*/
func parseExpr(query *entry.Query, expr sqlparser.Expr) (queryExpr *entry.QueryExpr, err error) {
	//queryExpr = &entry.QueryExpr{
	//	Table:               nil,
	//	ColumnQualifierName: "",
	//	ColumnName:          "",
	//	FuncName:            "",
	//	Name:                "",
	//	GoType:              nil,
	//}
	//condExpr := &entry.QueryExpr{
	//	Table:               nil,
	//	ColumnQualifierName: "",
	//	ColumnName:          "",
	//	FuncName:            "",
	//	Name:                "",
	//	GoType:              nil,
	//}
	//isArgCond := false
	//switch expr.(type) {
	//case *sqlparser.AndExpr:
	//	x := expr.(*sqlparser.AndExpr)
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Left)
	//	_, err = parseExpr(query, x.Right)
	//case *sqlparser.OrExpr:
	//	x := expr.(*sqlparser.OrExpr)
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Left)
	//	_, err = parseExpr(query, x.Right)
	//case *sqlparser.NotExpr:
	//	x := expr.(*sqlparser.NotExpr)
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Expr)
	//case *sqlparser.ParenExpr:
	//	x := expr.(*sqlparser.ParenExpr)
	//	queryExpr, err = parseExpr(query, x.Expr)
	//case *sqlparser.ComparisonExpr:
	//	x := expr.(*sqlparser.ComparisonExpr)
	//	if x.Operator == sqlparser.JSONExtractOp || x.Operator == sqlparser.JSONUnquoteExtractOp {
	//		err = fmt.Errorf("parse expr failed, json op (%s) is not supported", x.Operator)
	//		return
	//	}
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Left)
	//	_, err = parseExpr(query, x.Right)
	//	_, err = parseExpr(query, x.Escape)
	//case *sqlparser.RangeCond:
	//	x := expr.(*sqlparser.RangeCond)
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Left)
	//	_, err = parseExpr(query, x.From)
	//	_, err = parseExpr(query, x.To)
	//case *sqlparser.ExistsExpr:
	//	x := expr.(*sqlparser.ExistsExpr)
	//	queryExpr.GoType = entry.NewGoType("bool")
	//	_, err = parseExpr(query, x.Subquery)
	//case *sqlparser.SQLVal:
	//	x := expr.(*sqlparser.SQLVal)
	//	if x.Type == sqlparser.ValArg {
	//		isArgCond = true
	//	}
	//case *sqlparser.NullVal:
	//	//x := expr.(*sqlparser.NullVal)
	//
	//case sqlparser.BoolVal:
	////
	//case *sqlparser.ColName:
	//	x := expr.(*sqlparser.ColName)
	//	condExpr.ColumnName = x.Name.CompliantName()
	//	if x.Qualifier != nil {
	//		// as
	//		condExpr.ColumnQualifierName = x.Qualifier.Name.CompliantName()
	//	}
	//case sqlparser.ValTuple:
	//	x := expr.(sqlparser.ValTuple)
	//	err = x.WalkSubtree(func(node sqlparser.SQLNode) (kontinue bool, err error) {
	//		xExpr := node.(sqlparser.Expr)
	//		_, err = parseExpr(condExprList, xExpr)
	//		return
	//	})
	//case *sqlparser.Subquery:
	//	x := expr.(*sqlparser.Subquery)
	//	subQuery := &entry.Query{}
	//	err = parseQuerySelect(subQuery, x.Select)
	//	if err != nil {
	//		condExprList.ExprList = append(condExprList.ExprList, subQuery.CondExprList.ExprList...)
	//	}
	//case sqlparser.ListArg:
	//
	//case *sqlparser.BinaryExpr:
	//
	//case *sqlparser.UnaryExpr:
	//
	//case *sqlparser.IntervalExpr:
	//
	//case *sqlparser.CollateExpr:
	//
	//case *sqlparser.FuncExpr:
	//
	//case *sqlparser.CaseExpr:
	//
	//case *sqlparser.ValuesFuncExpr:
	//
	//case *sqlparser.ConvertExpr:
	//
	//case *sqlparser.MatchExpr:
	//
	//case *sqlparser.GroupConcatExpr:
	//
	//default:
	//	err = fmt.Errorf("parse expr failed, %v is unsupported", reflect.TypeOf(expr))
	//}

	return
}
