package conc

import (
	"sync"
)

type Promise[T any] struct {
	wg  sync.WaitGroup
	rst ResultT[T]
}

func NewPromise[T any](f func() (T, error)) *Promise[T] {
	p := &Promise[T]{}
	p.wg.Add(1)
	go func() {
		p.rst.Val, p.rst.Err = f()
		p.wg.Done()
	}()
	return p
}

func (p *Promise[T]) Wait() ResultT[T] {
	p.wg.Wait()
	return p.rst
}
