package sqltxx

import (
	"github.com/patrlind/go-sqltx"
)

// TXer handles transactions
type TXer interface {
	Querier
	sqltx.TXer
}
