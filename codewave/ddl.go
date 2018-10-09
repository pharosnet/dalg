package codewave

import (
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func waveDDL(tables []*def.Interface) error {
	w := NewWriter()
	dialect := strings.ToLower(strings.TrimSpace(tables[0].Dialect))
	switch dialect {
	case "postgres":
		wavePostgresDDL(w, tables)
	case "mysql":
		waveMysqlDDL(w, tables)
	case "oracle":
		waveOracleDDL(w, tables)
	default:
		return errors.New("wave DDL failed, the dialect is not supported")
	}
	filename := filepath.Join(codeFileFolder, "DDL.sql")
	var err error
	f, openErr := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if openErr != nil {
		logger.Log().Println(openErr)
		return openErr
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
			logger.Log().Println(err)
		}
	}()
	n, wToErr := w.WriteTo(f)
	if wToErr == nil && n < int64(w.Len()) {
		err = io.ErrShortWrite
		logger.Log().Println(err)
		return err
	}
	return nil
}

func wavePostgresDDL(w Writer, tables []*def.Interface) {
	for _, table := range tables {
		schema := strings.TrimSpace(table.Schema)
		if schema == "" {
			schema = "public"
		}
		name := strings.ToUpper(strings.TrimSpace(table.Name))
		w.WriteString(fmt.Sprintf(`-- Table: %s."%s"`, schema, name))
		w.WriteString("\n\n")
		w.WriteString(fmt.Sprintf(`DROP TABLE %s."%s";`, schema, name))
		w.WriteString("\n\n")
		w.WriteString(fmt.Sprintf(`CREATE TABLE %s."%s"`, schema, name))
		w.WriteString("\n")
		w.WriteString("(")
		w.WriteString("\n")
		for _, col := range table.Columns {
			w.WriteString(fmt.Sprintf(`	"%s" %s" `, strings.ToUpper(col.Name), col.Type))
			if col.MapType == "string" || col.MapType == "sql.NullString" {
				w.WriteString(` COLLATE pg_catalog."default"`)
			}
			if col.NotNull {
				w.WriteString(` NOT NULL`)
			}
			if col.Default != "" {
				if col.MapType == "string" || col.MapType == "sql.NullString" {
					w.WriteString(fmt.Sprintf(` DEFAULT '%s'::character varying`, col.Default))
				} else {
					w.WriteString(fmt.Sprintf(` DEFAULT %s`, col.Default))
				}
			}
			w.WriteString(", ")
			w.WriteString("\n")
		}
		pks := ""
		for i, pk := range table.Pks {
			if i > 0 {
				pks = pks + ","
			}
			pks = pks + fmt.Sprintf(`"%s"`, strings.ToUpper(pk.Name))
		}
		w.WriteString(fmt.Sprintf(`	CONSTRAINT "%s_PK" PRIMARY KEY (%s)`, name, pks))
		w.WriteString("\n")
		w.WriteString(")")
		w.WriteString("\n\n")
		w.WriteString(`WITH (`)
		w.WriteString("\n")
		w.WriteString(`	OIDS = FALSE`)
		w.WriteString("\n")
		w.WriteString(")")
		w.WriteString("\n\n")
		tablespace := "pg_default"
		if table.Tablespace != "" {
			tablespace = table.Tablespace
		}
		w.WriteString(fmt.Sprintf(`TABLESPACE %s;`, tablespace))
		w.WriteString("\n\n")
		w.WriteString(fmt.Sprintf(`ALTER TABLE %s."%s" OWNER to %s;`, schema, name, table.Owner))
		w.WriteString("\n\n")
		// index
		for _, index := range table.Indexes {
			idxName := strings.ToUpper(strings.TrimSpace(index.Name))
			w.WriteString(fmt.Sprintf(`Index: %s`, idxName))
			w.WriteString("\n\n")
			columns := strings.Split(index.Columns, ",")
			colExp := ""
			for i, col := range columns {
				if i > 0 {
					colExp = colExp + ", "
				}
				col := strings.ToUpper(strings.TrimSpace(col))
				colType := ""
				for _, column := range table.Columns {
					if col == strings.ToUpper(strings.TrimSpace(column.Name)) {
						colType = column.MapType
						break
					}
				}
				if colType == "string" || colType == "sql.NullString" {
					colExp = colExp + fmt.Sprintf(`"%s" COLLATE pg_catalog."default" %s %s`, col, index.Ops, index.SortOrder)
				} else {
					colExp = colExp + fmt.Sprintf(`"%s" %s`, col, index.SortOrder)
				}
			}
			unique := ""
			if index.Unique {
				unique = "UNIQUE"
			}
			w.WriteString(fmt.Sprintf(`CREATE %s INDEX "%s" ON %s."%s" USING %s (%s) TABLESPACE %s;`, unique, idxName, schema, name, index.Type, colExp, table.Tablespace))
			w.WriteString("\n\n")
		}
	}
}

func waveMysqlDDL(w Writer, tables []*def.Interface) {

}

func waveOracleDDL(w Writer, tables []*def.Interface) {

}
