package sqltx

import (
	"context"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	RegisterRetryableErrorTester(mockErrorTester{})
}

func TestSqltx(t *testing.T) {
	sleeper := &dummySleeper{}
	ctx := context.TODO()

	db, mock, err := sqlmock.New()
	if !assert.NoError(t, err) {
		return
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = Tx(ctx, db, &Options{Name: "1", Sleeper: sleeper}, func(tx TXer) error {
		if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
			return err
		}
		if _, err = tx.Exec("INSERT INTO product_viewers (user_id, product_id) VALUES (?, ?)", 2, 3); err != nil {
			return err
		}
		return nil
	})
	if !assert.NoError(t, err) {
		return
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	callCount := 0
	err = Tx(ctx, db, &Options{Name: "2", Sleeper: sleeper}, func(tx TXer) error {
		defer func() { callCount++ }()
		if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
			return err
		}
		if callCount == 0 {
			return fmt.Errorf("try again")
		}
		return nil
	})
	if !assert.NoError(t, err) {
		return
	}

	b := &backoff.ExponentialBackOff{
		InitialInterval:     DefaultBackoffInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      math.MaxInt64,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	for i := 1; i < 25; i++ {
		fmt.Println("i", i)
		for j := 0; j < i; j++ {
			fmt.Println("MOCK", j)
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
		}
		err = Tx(ctx, db, &Options{Name: fmt.Sprintf("3-%d", i), Retries: i, Sleeper: sleeper, Backoff: b}, func(tx TXer) error {
			fmt.Println("TX FUNC")
			if _, err = tx.Exec("UPDATE products SET views = views + 1"); err != nil {
				return err
			}
			return fmt.Errorf("try again")
		})
		if !assert.EqualError(t, err, fmt.Sprintf("transaction '3-%d' max retry count (%d) exceeded. Last error: try again", i, i)) {
			return
		}
	}
}

type mockErrorTester struct {
}

func (e mockErrorTester) IsRetryable(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "try again")
}

type dummySleeper struct {
}

// Sleep sleeps the duration of d
func (s dummySleeper) Sleep(ctx context.Context, d time.Duration) {
}
