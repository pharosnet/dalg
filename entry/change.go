package entry

type ChangeKind string

/*
	RENAME 					source="" target=new_table_name
	ADD						source="" target=new_column_name type
	MODIFY(ALTER COLUMN)	source="old_col" target=new_col type
	CHANGE					source="old_col" target=new_col type
	DROP(DROP COLUMN)
*/
type Change struct {
	Kind      ChangeKind
	TableName string
	Source    string
	Target    string
}
