package sqltxx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/patrlind/go-sqltx"
)

type sqlxDB struct {
	db *sqlx.DB
}

// TXFn is a transaction function that performs the steps of a transaction
type TXFn func(tx TXer) error

// Tx performs a database transaction
// If an error occurres the transaction will be automatically rolled back
// The transaction will be retried up to 20 times if the server reports a
// retryable error
func Tx(ctx context.Context, db *sqlx.DB, opts *sqltx.Options, fn TXFn) (err error) {
	return sqltx.TxHandler(ctx, &sqlxDB{db}, opts, func(tx sqltx.TXer) error {
		return fn(tx.(*sqlx.Tx))
	})
}

func (d *sqlxDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltx.TXer, error) {
	tx, err := d.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return tx.Tx, nil
}
