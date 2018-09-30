package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pharosnet/dalg/sample/postgres/dal"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags | log.Llongfile)
	db, dbErr := sql.Open("postgres", "host=101.132.179.213 port=15432 user=postgres password=freedom1581 dbname=postgres sslmode=disable")
	if dbErr != nil {
		logger.Println(dbErr)
		return
	}
	defer db.Close()
	dal.SetLog(logger)
	root := context.Background()
	tx, txErr := db.BeginTx(root, &sql.TxOptions{sql.LevelDefault, false})
	if txErr != nil {
		logger.Println(txErr)
		return
	}
	ctx := dal.WithPreparer(root, tx)
	row := dal.NewUserRow("4", "name3", 44, "MALE", 22.2, "1")
	row.Version.Int64 = 1
	n, err := dal.DeleteUser(ctx, row)
	if err != nil {
		logger.Println(err, tx.Rollback())
		return
	}
	logger.Println("cmt", tx.Commit(), n, *row)
}
