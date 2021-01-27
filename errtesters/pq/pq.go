package pq

import (
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/patrlind/go-sqltx"
)

type errorTester struct {
}

func init() {
	sqltx.RegisterRetryableErrorTester(errorTester{})
}

func (e errorTester) IsRetryable(err error) bool {
	pqErr := &pq.Error{}
	if errors.As(err, &pqErr) {
		if pqErr.Code == "40001" || pqErr.Code == "40P01" || pqErr.Code == "25P02" {
			if pqErr.Code == "25P02" && !strings.Contains(strings.ToLower(pqErr.Message), "try again") {
				return false
			}
			return true
		}
	}
	return false
}
