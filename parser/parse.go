package parser

import (
	"fmt"

	"github.com/pharosnet/dalg/entry"
	"github.com/pharosnet/dalg/files"
	"github.com/pharosnet/dalg/logs"
	"github.com/pharosnet/dalg/parser/mysql"
)

func Parse(dialect string, schemaPath string, queryPath string) (tables []*entry.Table, queries []*entry.Query, err error) {

	schemaFiles, schemaFilesErr := files.ReadFiles(schemaPath)
	if schemaFilesErr != nil {
		err = schemaFilesErr
		return
	}

	queryFiles, queryFilesErr := files.ReadFiles(queryPath)
	if queryFilesErr != nil {
		err = queryFilesErr
		return
	}

	tables = make([]*entry.Table, 0, 1)
	queries = make([]*entry.Query, 0, 1)

	switch dialect {
	case "mysql":
		schemaNameMap := make(map[string]string)
		for _, file := range schemaFiles {
			schema, parerErr := mysql.ParseMySQLSchema(string(file.Content))
			if parerErr != nil {
				err = fmt.Errorf("parse %s failed, \n%v", file.Name, parerErr)
				return
			}
			_, has := schemaNameMap[schema.Name]
			if has {
				err = fmt.Errorf("parse %s failed, schema %s has occurrd in other files", file.Name, schema.Name)
				return
			}
			schemaNameMap[schema.Name] = schema.Name
			tables = append(tables, schema.Tables...)
			logs.Log().Println("parse schema", file.Name, "succeed")
		}
		queryNameMap := make(map[string]string)
		for _, file := range queryFiles {
			queryList, parerErr := mysql.ParseMySQLQuery(string(file.Content))
			if parerErr != nil {
				err = fmt.Errorf("parse %s failed, \n%v", file.Name, parerErr)
				return
			}
			for _, query := range queryList {
				_, has := queryNameMap[query.Name]
				if has {
					err = fmt.Errorf("parse %s failed, query %s has occurrd in other files", file.Name, query.Name)
					return
				}
				queryNameMap[query.Name] = query.Name
				queries = append(queries, query)
			}
			logs.Log().Println("parse query", file.Name, "succeed")
		}
		var fillErr error
		queries, fillErr = entry.QueryFill(tables, queries)
		if fillErr != nil {
			err = fillErr
			return
		}
	case "postgres":
		err = fmt.Errorf("postgres is not support")
	}

	return
}
