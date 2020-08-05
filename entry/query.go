package entry

import (
	"fmt"
	"strings"

	"github.com/pharosnet/dalg/parser/commons"
)

const (
	SelectQueryKind = QueryKind("SELECT")
	InsertQueryKind = QueryKind("INSERT")
	UpdateQueryKind = QueryKind("UPDATE")
	DeleteQueryKind = QueryKind("DELETE")
)

type QueryKind string

func NewQuery() *Query {
	return &Query{
		fill:      false,
		Sql:       "",
		Kind:      "",
		Name:      "",
		TableList: make([]*QueryTable, 0, 1),
		SelectExprList: &SelectExprList{
			ExprList: make([]*QueryExpr, 0, 1),
		},
		CondExprList: &CondExprList{

			ExprList: make([]*CondExpr, 0, 1),
		},
	}
}

// mer
func QueryMergeCond(root *Query, sub *Query) {
	if root == nil || sub == nil {
		return
	}

	sub.Fill()
	//root.ExprList = append(root.ExprList, sub.ExprList...)
	if len(sub.CondExprList.ExprList) > 0 {
		root.TableList = append(root.TableList, sub.TableList...)
		root.CondExprList.ExprList = append(root.CondExprList.ExprList, sub.CondExprList.ExprList...)
	}

	return
}

type Query struct {
	fill           bool
	RawName        string
	Sql            string
	Kind           QueryKind
	Name           string
	TableList      []*QueryTable
	SelectExprList *SelectExprList
	CondExprList   *CondExprList
}

func (q *Query) Fill() {
	if q.fill {
		return
	}
	for _, expr := range q.SelectExprList.ExprList {
		for _, table := range q.TableList {
			if expr.ColumnQualifierName == table.Table || expr.ColumnQualifierName == table.NameAs {
				expr.Table = *table
				break
			}
		}
	}
	for _, expr := range q.CondExprList.ExprList {
		for _, table := range q.TableList {
			if expr.ColumnQualifierName == table.Table || expr.ColumnQualifierName == table.NameAs {
				expr.Table = *table
				break
			}
		}
	}
	q.fill = true
	return
}

type SelectExprList struct {
	ExprList []*QueryExpr
}

type CondExprList struct {
	ExprList []*CondExpr
}

type QueryExpr struct {
	Table               QueryTable
	ColumnQualifierName string // table or table as
	ColumnName          string // column or name as
	FuncName            string
	Name                string
	GoType              *GoType
}

func (e *QueryExpr) BuildName() {
	if e.Name != "" {
		return
	}
	//if e.ColumnQualifierName != "" {
	//	e.Name = commons.SnakeToCamel(e.ColumnQualifierName)
	//	return
	//}
	if e.ColumnName != "" {
		x := strings.ToLower(e.ColumnName)
		if e.FuncName != "" {
			x = fmt.Sprintf("%s_%s", x, e.FuncName)
		}
		e.Name = commons.SnakeToCamel(x)
	}
}

type QueryTable struct {
	Schema string
	Table  string
	NameAs string
	Ref    *Table
}

type CondExpr struct {
	Table               QueryTable
	ColumnQualifierName string // table or table as
	ColumnName          string // column or name as
	PlaceHolder         string
	Args                []string
	Name                string
	GoType              *GoType
	IsArg               bool
}

func (e *CondExpr) BuildName() {
	if e.Name != "" {
		return
	}
	//if e.ColumnQualifierName != "" {
	//	e.Name = commons.SnakeToCamel(e.ColumnQualifierName)
	//	return
	//}
	if e.ColumnName != "" {
		x := strings.ToLower(e.ColumnName)
		e.Name = commons.SnakeToCamel(x)
	}
}
