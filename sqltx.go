package sqltx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// TXFn is a transaction function that performs the steps of a transaction
type TXFn func(tx TXer) error

// TXBeginner begins a database transaction
type TXBeginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TXer, error)
}

// RetryableErrorTester tests if an error is to be retried in the transaction
type RetryableErrorTester interface {
	IsRetryable(err error) bool
}

// Sleeper interface for context aware sleep
type Sleeper interface {
	Sleep(ctx context.Context, d time.Duration)
}

// DefaultSleeper is the default sleeper
type DefaultSleeper struct {
}

// Options contains transaction options
type Options struct {
	Name      string
	Isolation sql.IsolationLevel
	ReadOnly  bool
	Retries   int
	Backoff   backoff.BackOff
	Sleeper   Sleeper
}

const (
	// DefaultRetries is the default number of retries to perform
	DefaultRetries = 20
	// DefaultBackoffInitialInterval is the default backoff initial interval
	DefaultBackoffInitialInterval = 1 * time.Millisecond
	// DefaultBackoffMaxInterval is the default backoff max interval
	DefaultBackoffMaxInterval = 5 * time.Second
)

var (
	maxRetries = DefaultRetries

	backoffObject backoff.BackOff = &backoff.ExponentialBackOff{
		InitialInterval:     DefaultBackoffInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      DefaultBackoffMaxInterval,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}

	testers []RetryableErrorTester

	defaultSleeper = &DefaultSleeper{}
)

type sqlDB struct {
	db *sql.DB
}

// SetDefaultMaxRetries sets the default value for max retries
// Only set this once during initialization
func SetDefaultMaxRetries(retries int) {
	maxRetries = retries
}

// GetDefaultMaxRetries returns the default max number of retries
func GetDefaultMaxRetries() int {
	return maxRetries
}

// SetDefaultBackoff sets the default value for the backoff
// Only set this once during initialization
func SetDefaultBackoff(b backoff.BackOff) {
	backoffObject = b
}

// GetDefaultBackoff returns the default backoff object
func GetDefaultBackoff() backoff.BackOff {
	return backoffObject
}

// RegisterRetryableErrorTester registers a test function that is used
// to determine if a database error should be retried or not
func RegisterRetryableErrorTester(t RetryableErrorTester) {
	testers = append(testers, t)
}

// Tx performs a database transaction
// If an error occurres the transaction will be automatically rolled back
// The transaction will be retried up to 20 times if the server reports a
// retryable error
func Tx(ctx context.Context, db *sql.DB, opts *Options, fn TXFn) (err error) {
	return TxHandler(ctx, &sqlDB{db}, opts, fn)
}

// TxHandler is the handler that handles transaction details
func TxHandler(ctx context.Context, db TXBeginner, opts *Options, fn TXFn) (err error) {
	if len(testers) == 0 {
		panic("no sqltx error testers registered, please register by importing the relevant sqltx DB implementation(s)")
	}

	var tx TXer
	defer func() {
		nilTX := tx == nil || (reflect.ValueOf(tx).Kind() == reflect.Ptr && reflect.ValueOf(tx).IsNil())
		if p := recover(); p != nil {
			// A panic occurred, rollback and re-panic
			if !nilTX {
				_ = tx.Rollback()
			}
			panic(p)
		} else if err != nil {
			// Something went wrong, rollback
			if !nilTX {
				_ = tx.Rollback()
			}
		}
	}()

	var b backoff.BackOff
	useBackoff := backoffObject
	retries := maxRetries
	var txOpts *sql.TxOptions
	var sleeper Sleeper = defaultSleeper
	var name = ""
	if opts != nil {
		txOpts = &sql.TxOptions{
			Isolation: opts.Isolation,
			ReadOnly:  opts.ReadOnly,
		}
		if opts.Retries > 0 {
			retries = opts.Retries
		}
		if opts.Backoff != nil {
			useBackoff = opts.Backoff
		}
		if opts.Sleeper != nil {
			sleeper = opts.Sleeper
		}
		if opts.Name != "" {
			name = " '" + opts.Name + "'"
		}
	}
	for i := 0; i < retries; i++ {
		tx, err = db.BeginTx(ctx, txOpts)
		if err != nil {
			return fmt.Errorf("failed to start transaction%s: %w", name, err)
		}
		err = fn(tx)
		if err == nil {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("transaction%s commit error: %w", name, err)
			}
		}
		if err != nil {
			retryable := false
			for _, tester := range testers {
				if tester.IsRetryable(err) {
					retryable = true
					break
				}
			}
			if !retryable {
				return err
			}
			// Retryable error, try again (soon)
			if b == nil {
				b = useBackoff
				b.Reset()
			}
			sleepTime := b.NextBackOff()
			if sleepTime == backoff.Stop {
				return fmt.Errorf("transaction%s backoff max time reached", name)
			}
			sleeper.Sleep(ctx, sleepTime)
			continue
		}
		return err
	}
	return fmt.Errorf("transaction%s max retry count (%d) exceeded. Last error: %w", name, retries, err)

}

// Sleep sleeps the duration of d
func (s DefaultSleeper) Sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func (d *sqlDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (TXer, error) {
	return d.db.BeginTx(ctx, opts)
}
