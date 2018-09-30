package main

import (
	_ "github.com/lib/pq"
	"context"
	"database/sql"
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
	ctx := dal.WithPreparer(root, db)
	userRow, err := dal.LoadUserRow(ctx, "55")
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("%+v", userRow)
}
