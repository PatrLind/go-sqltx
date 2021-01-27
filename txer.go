package sqltx

// TXer handles transactions
type TXer interface {
	Querier
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error
}
