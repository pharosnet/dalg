package mysql

import "github.com/pharosnet/dalg/entry"

const (
	TableDropChangeKind    = "DROP TABLE"
	TableRenameChangeKind  = "RENAME TABLE"
	ColumnAddChangeKind    = "ADD COLUMN"
	ColumnModifyChangeKind = "MODIFY COLUMN"
	ColumnChangeChangeKind = "CHANGE COLUMN"
	ColumnDropChangeKind   = "DROP COLUMN"
)

type ChangeKind string

/*
	DROP TABLE						source=""
	RENAME TABLE					source="" target=new_table_name
	ADD	COLUMN						source="" target=new_column_name type
	MODIFY COLUMN					source="old_col" target=type
	CHANGE COLUMN					source="old_col" target=new_col type
	DROP COLUMN						source="old_col"
*/
type Change struct {
	Kind    ChangeKind
	Schema  string
	Table   string
	Target  interface{}
	Content string
}

type ChangDropTable struct{}

type ChangeRenameTable struct {
	Schema string
	Name   string
}

type ChangeAddColumn struct {
	Name         string
	Type         entry.ColumnType
	DefaultValue string
	First        bool
	After        string
}

type ChangeModifyColumn struct {
	Source       string
	Type         entry.ColumnType
	DefaultValue string
}

type ChangeChangeColumn struct {
	Source       string
	Name         string
	Type         entry.ColumnType
	DefaultValue string
}

type ChangDropColumn struct {
	Name string
}
