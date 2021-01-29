package sqltxx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/patrlind/go-sqltx"
)

// Querier interface can be used in helper functions that can take either a
// transaction or a database connection
type Querier interface {
	sqltx.Querier

	// DriverName returns the driverName passed to the Open function for this DB.
	DriverName() string

	// Rebind transforms a query from QUESTION to the DB driver's bindvar type.
	Rebind(query string) string

	// BindNamed binds a query using the DB driver's bindvar type.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	// NamedQuery using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)

	// NamedExec using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExec(query string, arg interface{}) (sql.Result, error)

	// NamedExecContext using this DB.
	// Any named placeholder parameters are replaced with fields from arg.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	// Select using this DB.
	// Any placeholder parameters are replaced with supplied args.
	Select(dest interface{}, query string, args ...interface{}) error

	// SelectContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	Get(dest interface{}, query string, args ...interface{}) error

	// GetContext using this DB.
	// Any placeholder parameters are replaced with supplied args.
	// An error is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Queryx queries the database and returns an *sqlx.Rows.
	// Any placeholder parameters are replaced with supplied args.
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryxContext queries the database and returns an *sqlx.Rows.
	// Any placeholder parameters are replaced with supplied args.
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// QueryRowx queries the database and returns an *sqlx.Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowx(query string, args ...interface{}) *sqlx.Row

	// QueryRowxContext queries the database and returns an *sqlx.Row.
	// Any placeholder parameters are replaced with supplied args.
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// MustExec (panic) runs MustExec using this database.
	// Any placeholder parameters are replaced with supplied args.
	MustExec(query string, args ...interface{}) sql.Result

	// MustExecContext (panic) runs MustExec using this database.
	// Any placeholder parameters are replaced with supplied args.
	MustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result

	// Preparex returns an sqlx.Stmt instead of a sql.Stmt
	Preparex(query string) (*sqlx.Stmt, error)

	// PreparexContext returns an sqlx.Stmt instead of a sql.Stmt.
	//
	// The provided context is used for the preparation of the statement, not for
	// the execution of the statement.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// PrepareNamed returns an sqlx.NamedStmt
	PrepareNamed(query string) (*sqlx.NamedStmt, error)

	// PrepareNamedContext returns an sqlx.NamedStmt
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}
