package timeout

import (
	"context"
	"time"
)

type DatabaseOperation[T any] func(ctx context.Context) (T, error)

func WithTimeout[T any](ctx context.Context, timeout time.Duration, operation DatabaseOperation[T]) (T, error) {
	var zero T

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Channel to receive the result
	type result struct {
		data T
		err  error
	}

	resultChan := make(chan result, 1)

	// Execute operation in goroutine
	go func() {
		data, err := operation(timeoutCtx)
		resultChan <- result{data: data, err: err}
	}()

	// Wait for either completion or timeout
	select {
	case res := <-resultChan:
		return res.data, res.err
	case <-timeoutCtx.Done():
		return zero, context.DeadlineExceeded
	}
}
