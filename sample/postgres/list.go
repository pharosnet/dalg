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
	rows := make([]*dal.UserRow, 0, 5)
	listErr := dal.ListUserRows(ctx, 10, 0, func(ctx context.Context, row *dal.UserRow, err error) error {
		if err != nil {
			logger.Println(err)
			return err
		}
		rows = append(rows, row)
		return nil
	})
	if listErr != nil {
		logger.Println(listErr)
		return
	}
	for i, row := range rows {
		logger.Printf("%d -> %+v", i, row)
	}

}
