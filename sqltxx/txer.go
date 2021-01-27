package sqltxx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/patrlind/go-sqltx"
)

// TXer handles transactions
type TXer interface {
	sqltx.TXer

	// DriverName returns the driverName used by the DB which began this transaction.
	DriverName() string

	// Rebind a query within a transaction's bindvar type.
	Rebind(query string) string

	// Unsafe returns a version of Tx which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	Unsafe() *sqlx.Tx

	// BindNamed binds a query within a transaction's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	// NamedQuery within a transaction.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	// NamedExec a named query within a transaction.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)

	// Select within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error

	// Queryx within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryRowx within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	// Get within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error

	// MustExec runs MustExec within a transaction.
	// Any placeholder parameters are replaced with supplied args.
	MustExec(query string, args ...interface{}) sql.Result

	// Preparex  a statement within a transaction.
	Preparex(query string) (*sqlx.Stmt, error)

	// Stmtx returns a version of the prepared statement which runs within a transaction.  Provided
	// stmt can be either *sql.Stmt or *sqlx.Stmt.
	Stmtx(stmt interface{}) *sqlx.Stmt

	// NamedStmt returns a version of the prepared statement which runs within a transaction.
	NamedStmt(stmt *sqlx.NamedStmt) *sqlx.NamedStmt

	// PrepareNamed returns an sqlx.NamedStmt
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}
