package sqltx

import (
	"context"
	"database/sql"
)

// Querier interface can be used in helper functions that can take either a
// transaction or a database connection
type Querier interface {
	// PrepareContext creates a prepared statement for use within a transaction.
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	//
	// To use an existing prepared statement on this transaction, see Tx.Stmt.
	//
	// The provided context will be used for the preparation of the context, not
	// for the execution of the returned statement. The returned statement
	// will run in the transaction context.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Prepare creates a prepared statement for use within a transaction.
	//
	// The returned statement operates within the transaction and can no longer
	// be used once the transaction has been committed or rolled back.
	//
	// To use an existing prepared statement on this transaction, see Tx.Stmt.
	Prepare(query string) (*sql.Stmt, error)

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

	// ExecContext executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Exec executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// Query executes a query that returns rows, typically a SELECT.
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRow(query string, args ...interface{}) *sql.Row
}
