package conc

import (
	"context"
	"errors"
	"time"
)

var ErrTimeout = errors.New("timeout")

func WaitOrTimeout[T any](fx func() (T, error), timeout time.Duration) (T, error) {
	ch := make(chan ResultT[T])
	go func() {
		rst, err := fx()
		ch <- ResultT[T]{Val: rst, Err: err}
	}()
	select {
	case rst := <-ch:
		return rst.Val, rst.Err
	case <-time.After(timeout):
		var zero T
		return zero, ErrTimeout
	}
}

func WaitOrTimeoutWithContext[T any](ctx context.Context, fx func() (T, error)) (T, error) {
	ch := make(chan ResultT[T])
	go func() {
		rst, err := fx()
		ch <- ResultT[T]{Val: rst, Err: err}
	}()
	select {
	case rst := <-ch:
		return rst.Val, rst.Err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}
