package mysql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitessio/vitess/go/vt/sqlparser"

	"github.com/pharosnet/dalg/entry"
	"github.com/pharosnet/dalg/parser/commons"
)

func ParseMySQLQuery(content string) (queries []*entry.Query, err error) {

	if !strings.Contains(content, ";") {
		err = fmt.Errorf("parse query failed, there is no ; ")
		return
	}

	queries = make([]*entry.Query, 0, 1)

	blocks := strings.Split(content, ";")

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if len(block) == 0 || len(commons.ReadWords([]byte(block))) == 0 {
			continue
		}
		name := ""
		querySQL := ""

		lines := commons.NewLines(block)
		for lines.HasNext() {
			line := lines.NextLine()
			upperLine := strings.ToUpper(line)
			words := commons.ReadWords([]byte(upperLine))
			if commons.WordsContainsAll(words, "--", "NAME:") {
				nameIdx := commons.WordsIndex(words, "NAME:")
				name = lines.CurrentLineWords()[nameIdx+1]
				querySQL = lines.Remain()
			}
			if name != "" {
				break
			}
			if commons.WordsContainsOne(words, "SELECT", "INSERT", "UPDATE", "DELETE") {
				break
			}
		}
		if name == "" || querySQL == "" {
			err = fmt.Errorf("parse query failed, can not read name: for query block, or sql in query block,\n%s", block)
			return
		}

		query, parseQueryErr := parseMySQLQuery0(name, querySQL)
		if parseQueryErr != nil {
			err = parseQueryErr
			return
		}
		queries = append(queries, query)
	}

	nameMap := make(map[string]string)
	for _, query := range queries {
		_, has := nameMap[query.Name]
		if has {
			err = fmt.Errorf("parse query failed, query name is repeat, %s", query.Name)
			return
		}
		nameMap[query.Name] = query.Name
	}

	return
}

func parseMySQLQuery0(name string, content string) (query *entry.Query, err error) {
	querySQL, _ := sqlparser.SplitTrailingComments(content)
	stmt, parseErr := sqlparser.Parse(strings.ToUpper(querySQL))
	if parseErr != nil {
		err = fmt.Errorf("parse query failed, %v, sql:\n%s", parseErr, querySQL)
		return
	}
	query = entry.NewQuery()
	query.RawName = strings.ToLower(name)
	query.Name = commons.SnakeToCamel(strings.ToLower(name))
	query.Sql = strings.ReplaceAll(querySQL, "\t", " ")

	switch stmt.(type) {
	case sqlparser.SelectStatement:
		query.Kind = entry.SelectQueryKind
		err = parseQuerySelect(query, stmt.(sqlparser.SelectStatement))
	case *sqlparser.Insert:
		query.Kind = entry.InsertQueryKind
		err = parseQueryInsert(query, stmt.(*sqlparser.Insert))
	case *sqlparser.Update:
		query.Kind = entry.UpdateQueryKind
		err = parseQueryUpdate(query, stmt.(*sqlparser.Update))
	case *sqlparser.Delete:
		query.Kind = entry.DeleteQueryKind
		err = parseQueryDelete(query, stmt.(*sqlparser.Delete))
	default:
		err = fmt.Errorf("parse query failed, %v is unsupported", reflect.TypeOf(stmt))
	}

	if err != nil {
		err = fmt.Errorf("parse query failed, %v, sql:\n%s", err, querySQL)
	}
	query.Fill()

	return
}
