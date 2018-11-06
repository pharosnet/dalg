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

func waveDDL(tables []def.Interface) error {
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

func wavePostgresDDL(w Writer, tables []def.Interface) {
	for _, table := range tables {
		table.Pks = make([]def.Column, 0, 1)
		schema := strings.TrimSpace(table.Schema)
		if schema == "" {
			schema = "public"
		}
		name := strings.ToUpper(strings.TrimSpace(table.Name))
		w.WriteString(fmt.Sprintf(`-- Table: %s."%s"`, schema, name))
		w.WriteString("\n\n")
		w.WriteString(fmt.Sprintf(`DROP TABLE IF EXISTS %s."%s";`, schema, name))
		w.WriteString("\n\n")
		w.WriteString(fmt.Sprintf(`CREATE TABLE %s."%s"`, schema, name))
		w.WriteString("\n")
		w.WriteString("(")
		w.WriteString("\n")
		for _, col := range table.Columns {
			if col.Pk {
				table.Pks = append(table.Pks, col)
			}
			w.WriteString(fmt.Sprintf(`	"%s" %s `, strings.ToUpper(col.Name), col.Type))
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
			idxName := fmt.Sprintf(`%s_IDX_%s`, name, strings.ToUpper(strings.TrimSpace(index.Name)))
			w.WriteString(fmt.Sprintf(`-- Index: %s`, idxName))
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

func waveMysqlDDL(w Writer, tables []def.Interface) {
	for _, table := range tables {
		table.Pks = make([]def.Column, 0, 1)
		schema := strings.TrimSpace(table.Schema)
		name := strings.ToUpper(strings.TrimSpace(table.Name))
		if schema != "" {
			w.WriteString(fmt.Sprintf("-- Table: `%s`.`%s`", schema, name))
			w.WriteString("\n\n")
			w.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`.`%s`;", schema, name))
			w.WriteString("\n\n")
			w.WriteString(fmt.Sprintf("CREATE TABLE `%s`.`%s` (", schema, name))
			w.WriteString("\n")
		} else {
			w.WriteString(fmt.Sprintf("-- Table: %s", name))
			w.WriteString("\n\n")
			w.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", name))
			w.WriteString("\n\n")
			w.WriteString(fmt.Sprintf("CREATE TABLE `%s` (", name))
			w.WriteString("\n")
		}
		for _, col := range table.Columns {
			if col.Pk {
				table.Pks = append(table.Pks, col)
			}
			w.WriteString(fmt.Sprintf("	`%s` %s ", strings.ToUpper(col.Name), col.Type))
			if col.NotNull {
				w.WriteString(` NOT NULL`)
			}
			if col.Default != "" {
				if col.MapType == "string" || col.MapType == "sql.NullString" {
					w.WriteString(fmt.Sprintf(` DEFAULT '%s'`, col.Default))
				} else {
					w.WriteString(fmt.Sprintf(` DEFAULT %s`, col.Default))
				}
			}
			w.WriteString(", ")
			w.WriteString("\n")
		}
		pks := ""
		for _, key := range table.Pks {
			pks = pks + "," + fmt.Sprintf("`%s`", strings.ToUpper(key.Name))
		}
		if len(pks) > 0 {
			pks = pks[1:]
			w.WriteString(fmt.Sprintf("PRIMARY KEY (%s)", pks))
		}
		if table.Indexes != nil && len(table.Indexes) > 0 {
			w.WriteString(",")
			w.WriteString("\n")
			for indexIndex, index := range table.Indexes {

				idxName := fmt.Sprintf(`%s_IDX_%s`, name, strings.ToUpper(strings.TrimSpace(index.Name)))
				unique := ""
				if index.Unique {
					unique = "UNIQUE"
				}
				indexType := ""
				if index.Type != "" {
					indexType = fmt.Sprintf("USING %s", index.Type)
				}
				columns := strings.Split(index.Columns, ",")
				colExp := ""
				for i, col := range columns {
					if i > 0 {
						colExp = colExp + ", "
					}
					colExp = colExp + fmt.Sprintf("`%s`", strings.ToUpper(strings.TrimSpace(col)))
				}
				w.WriteString(fmt.Sprintf("%s KEY `%s` (%s) %s", unique, idxName, colExp, indexType))
				if indexIndex + 1 < len(table.Indexes) {
					w.WriteString(",")
				}
				w.WriteString("\n")
			}
		} else {
			w.WriteString("\n")
		}
		w.WriteString(");")
		w.WriteString("\n\n")
	}
}

func waveOracleDDL(w Writer, tables []def.Interface) {
	w.WriteString("-- not support")
	// TODO 
}
