package dal

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	user_insert_sql = `INSERT INTO "user"("id", "name", "age", "sex", "money", "info", "create_by", "create_time", "modify_by", "modify_time", "version") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	user_update_Sql = `UPDATE "user" SET "name" = $1, "age" = $2, "sex" = $3, "money" = $4, "info" = $5, "create_by" = $6, "create_time" = $7, "modify_by" = $8, "modify_time" = $9, "version" = "version" + 1 WHERE "id" = $10 AND "version" = $11`
	user_delete_Sql = `DELETE FROM "user" WHERE "id" = $1 AND "version" = $2`
	user_getone_sql = `SELECT "id", "name", "age", "sex", "money","info", "create_by", "create_time", "modify_by", "modify_time", "version" FROM "user" WHERE "id" = $1`
	//
	user_ext_sql_list = `SELECT "id", "name", "age", "sex", "money", "info", "create_by", "create_time", "modify_by", "modify_time", "version" FROM "user" LIMIT $1 OFFSET $2`
)

func NewUserSex(v string) UserSex {
	v = strings.ToUpper(v)
	ok := false
	switch v {
	case "MALE":
		ok = true
	case "FEMALE":
		ok = true
	}
	if !ok {
		panic(fmt.Errorf("dal-> new user sex failed, value is invalid"))
	}
	return UserSex{v, true}
}

type UserSex struct {
	Sex string
	Valid bool
}

func (n *UserSex) Scan(value interface{}) error {
	if value == nil {
		n.Sex, n.Valid = "unknown", false
		return nil
	}

	vv, ok := value.(bool)
	if !ok {
		return fmt.Errorf("dal-> scan value failed, value type is not bool")
	}
	switch vv {
	case true:
		n.Sex = "MALE"
	case false:
		n.Sex = "FEMALE"
	default:
		return fmt.Errorf("dal-> scan value failed, value is out of range")
	}
	n.Valid = true

	return nil
}

func (n UserSex) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	switch n.Sex {
	case "MALE":
		return true, nil
	case "FEMALE":
		return false, nil
	}
	return nil, fmt.Errorf("dal-> %s is not sex value", n.Sex)
}

type UserInfo struct {
	NullJson  `json:"-"`
	Id string `json:"id"`
	Age int64 `json:"age"`
}

func (n *UserInfo) Scan(value interface{}) error {
	if value == nil {
		n.NullJson.Bytes, n.NullJson.Valid = nil, false
		return nil
	}
	switch value.(type) {
	case []byte:
		n.NullJson.Bytes = value.([]byte)
	case string:
		n.NullJson.Bytes = []byte(value.(string))
	default:
		return errors.New("UserInfo scan value failed, value type is invalid")
	}
	if err := json.Unmarshal(n.NullJson.Bytes, n); err != nil {
		return nil
	}
	n.NullJson.Valid = true
	return nil
}

func (n UserInfo) Value() (driver.Value, error) {
	if !n.NullJson.Valid {
		return nil, nil
	}
	p, err := json.Marshal(&n)
	if err != nil {
		return nil, err
	}
	return p, nil
}


func NewUserRow(id string, name string, age int64, sex string, money float64, createBy string) *UserRow {
	now := nowTime()
	return &UserRow{
		Id: sql.NullString{id, true},
		Name: sql.NullString{name, true},
		Age: sql.NullInt64{age, true},
		Sex: NewUserSex(sex),
		Money: sql.NullFloat64{money, true},
		CreateBy: sql.NullString{createBy, true},
		CreateTime:now,
		ModifyBy:sql.NullString{createBy, true},
		ModifyTime:now,
		Version:sql.NullInt64{0, true},
	}
}

type UserRow struct {
	Id sql.NullString
	Name sql.NullString
	Age sql.NullInt64
	Sex UserSex
	Money sql.NullFloat64
	Info UserInfo
	CreateBy sql.NullString
	CreateTime NullTime
	ModifyBy sql.NullString
	ModifyTime NullTime
	Version sql.NullInt64
}

func (row *UserRow) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "(id: %v, name: %v, age: %v, sex: %v, money: %v, info: %v, create_by: %v, create_time: %v, modify_by: %v, modify_time: %v, version: %v)",
				row.Id, row.Name, row.Age, row.Sex, row.Money, row.Info, row.CreateBy, row.CreateTime, row.ModifyBy, row.ModifyTime, row.Version)
		default:
			fmt.Fprintf(s, "&{%v, %v, %v, %v, %v, %v, %v, %v, %v, %v, %v}",
				row.Id, row.Name, row.Age, row.Sex, row.Money, row.Info, row.CreateBy, row.CreateTime, row.ModifyBy, row.ModifyTime, row.Version)
		}
	}
}

func scanUserRow(sa scanner) (userRow *UserRow, err error) {
	userRow = &UserRow{}
	scanErr := sa.Scan(
			&userRow.Id,
			&userRow.Name,
			&userRow.Age,
			&userRow.Sex,
			&userRow.Money,
			&userRow.Info,
			&userRow.CreateBy,
			&userRow.CreateTime,
			&userRow.ModifyBy,
			&userRow.ModifyTime,
			&userRow.Version,
		)
	if scanErr != nil {
		err = fmt.Errorf("dal-> scan failed. reason: %v", scanErr)
		return
	}
	return
}

type UserRowRangeFn func(ctx context.Context, row *UserRow, err error) error

func InsertUser(ctx context.Context, rows ...*UserRow) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal-> insert user failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal-> insert user failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, user_insert_sql)
	if prepareErr != nil {
		err = fmt.Errorf("dal-> insert user failed, prepared statment failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal-> insert user failed, close prepare statment failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, row.Id, row.Name, row.Age, row.Sex, row.Money, row.Info, row.CreateBy, row.CreateTime, row.ModifyBy, row.ModifyTime, row.Version)
		if execErr != nil {
			err = fmt.Errorf("dal-> insert user failed, execute statment failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal-> insert user failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal-> insert user failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		//id, getIdErr := result.LastInsertId()
		//if getIdErr != nil {
		//	err = fmt.Errorf("dal-> insert user failed, get last insert id failed. reason: %v", getIdErr)
		//	return
		//}
		//if id < 0 {
		//	err = errors.New("dal-> insert user failed, get last insert id failed. id is invalid")
		//	return
		//}
		//row.Id = id
		if hasLog() {
			logf("dal-> insert user success, affected : %d, sql : %s, row : %+v\n", affectedRows, user_insert_sql, row)
		}
	}
	return
}

func UpdateUser(ctx context.Context, rows ...*UserRow) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal-> update user failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal-> update user failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, user_update_Sql)
	if prepareErr != nil {
		err = fmt.Errorf("dal-> update user failed, prepared statement failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal-> update user failed, close prepare statement failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, row.Name, row.Age, row.Sex, row.Money, row.Info, row.CreateBy, row.CreateTime, row.ModifyBy, row.ModifyTime, row.Id, row.Version)
		if execErr != nil {
			err = fmt.Errorf("dal-> update user failed, execute statement failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal-> update user failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal-> update user failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		if hasLog() {
			logf("dal-> update user success, affected : %d, sql : %s, row : %+v\n", affectedRows, user_update_Sql, row)
		}
		row.Version.Int64 ++
	}
	return
}

func DeleteUser(ctx context.Context, rows ...*UserRow) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal-> delete user failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal-> delete user failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, user_delete_Sql)
	if prepareErr != nil {
		err = fmt.Errorf("dal-> delete user failed, prepared statment failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal-> delete user failed, close prepare statment failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, row.Id, row.Version)
		if execErr != nil {
			err = fmt.Errorf("dal-> delete user failed, execute statment failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal-> delete user failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal-> delete user failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		if hasLog() {
			logf("dal-> delete user success, affected : %d, sql : %s, row : %+v\n", affectedRows, user_delete_Sql, row)
		}
	}
	return
}

func LoadUserRow(ctx context.Context, id string) (userRow *UserRow, err error) {
	if ctx == nil {
		err = errors.New("dal-> load user failed, context is empty")
		return
	}
	if id == "" {
		err = errors.New("dal-> load user failed, id is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, user_getone_sql)
	if prepareErr != nil {
		err = fmt.Errorf("dal-> load user failed, prepared statment failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal-> load user failed, close prepare statment failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	row := stmt.QueryRowContext(ctx, &id)
	userRow, err = scanUserRow(row)
	if hasLog() {
		logf("dal-> load user success, sql : %s, id : %v, row : %+v\n", user_getone_sql, id, userRow)
	}
	return
}

func ListUserRows(ctx context.Context, limit int64, offset int64, rangeFn UserRowRangeFn) (err error) {
	if ctx == nil {
		err = errors.New("dal-> list user failed, context is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, user_ext_sql_list)
	if prepareErr != nil {
		err = fmt.Errorf("dal-> list user failed, prepared statment failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal-> list user failed, close prepare statment failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	rows, queryErr := stmt.QueryContext(ctx, &limit, &offset)
	if queryErr != nil {
		err = fmt.Errorf("dal-> list user failed, query failed. reason: %v", queryErr)
		return
	}
	defer func() {
		closeErr := rows.Close()
		if closeErr != nil {
			err = fmt.Errorf("dal-> list user failed, rows is invalid. reason: %v", closeErr)
			return
		}
	}()
	for rows.Next() {
		userRow, scanErr := scanUserRow(rows)
		rangeFnErr := rangeFn(ctx, userRow, scanErr)
		if rangeFnErr != nil {
			err = rangeFnErr
			return
		}
	}
	if hasLog() {
		logf("dal-> load user success, sql : %s, limit : %v, length : %v\n", user_ext_sql_list, limit, offset)
	}
	rowsErr := rows.Err()
	if rowsErr != nil {
		err = fmt.Errorf("dal-> list user failed, rows is invalid. reason: %v", rowsErr)
		return
	}
	return
}