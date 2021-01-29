package sqltx

import (
	"context"
	"database/sql"
)

// TXer handles transactions
type TXer interface {
	Querier
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	// StmtContext returns a transaction-specific prepared statement from
	// an existing statement.
	//
	// Example:
	//  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
	//  ...
	//  tx, err := db.Begin()
	//  ...
	//  res, err := tx.StmtContext(ctx, updateMoney).Exec(123.45, 98293203)
	//
	// The provided context is used for the preparation of the statement, not for the
	// execution of the statement.
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt

	// Stmt returns a transaction-specific prepared statement from
	// an existing statement.
	//
	// Example:
	//  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
	//  ...
	//  tx, err := db.Begin()
	//  ...
	//  res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	Stmt(stmt *sql.Stmt) *sql.Stmt
}
