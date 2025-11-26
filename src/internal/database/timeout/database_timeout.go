package timeout

import (
	"context"
	"fmt"
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

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("database operation panic: %v", r)
				select {
				case resultChan <- result{err: err}:
				case <-timeoutCtx.Done():
				}
			}
		}()
		data, err := operation(timeoutCtx)
		select {
		case resultChan <- result{data: data, err: err}:
		case <-timeoutCtx.Done():
		}
	}()
	select {
	case res := <-resultChan:
		return res.data, res.err
	case <-timeoutCtx.Done():
		return zero, context.DeadlineExceeded
	}
}
