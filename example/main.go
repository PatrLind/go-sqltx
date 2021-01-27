package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/patrlind/go-sqltx"

	_ "github.com/patrlind/go-sqltx/errtesters/pq"
)

func main() {
	ctx := context.TODO()
	db := getMockDB()
	defer db.Close()

	// The Tx function will call the supplied function if the queries needs
	// to be retried.
	err := sqltx.Tx(ctx, db, &sqltx.Options{Isolation: sql.LevelSerializable}, func(tx sqltx.TXer) error {
		var name, value string = "Name 1", "Value 1"
		_, err := tx.ExecContext(ctx, "INSERT INTO data (name, value) VALUES ($1, $2)", name, value)
		if err != nil {
			return fmt.Errorf("query error: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = sqltx.Tx(ctx, db, &sqltx.Options{ReadOnly: true}, func(tx sqltx.TXer) error {
		rows, err := tx.QueryContext(ctx, "SELECT name, value FROM data")
		if err != nil {
			return fmt.Errorf("query error: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var name, value string
			err := rows.Scan(&name, &value)
			if err != nil {
				return fmt.Errorf("scan error: %w", err)
			}
			fmt.Printf("%s = %s\n", name, value)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func getMockDB() *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO data").WithArgs("Name 1", "Value 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT name, value FROM data").WillReturnRows(
		sqlmock.NewRows([]string{"name", "value"}).
			AddRow("Name 1", "Value 1"))
	mock.ExpectCommit()
	return db
}
