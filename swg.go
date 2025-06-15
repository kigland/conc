package conc

import (
	"sort"
	"sync"
)

type SortedWaitGroup struct {
	lck sync.Mutex
	id  int
	wg  sync.WaitGroup
	rst []idPair
}

func (w *SortedWaitGroup) AddTaskWithoutErr(f func() any) *SortedWaitGroup {
	return w.AddTask(func() (any, error) {
		return f(), nil
	})
}

func (w *SortedWaitGroup) AddTask(f func() (any, error)) *SortedWaitGroup {
	w.lck.Lock()
	id := w.id
	w.id++
	w.lck.Unlock()

	w.wg.Add(1)
	go func(id int) {
		defer w.wg.Done()
		rst, err := f()

		w.lck.Lock()
		w.rst = append(w.rst, idPair{id, Result{Val: rst, Err: err}})
		w.lck.Unlock()
	}(id)

	return w
}

func (w *SortedWaitGroup) Wait() AllResults {
	w.wg.Wait()
	rsts := w.rst

	sort.Slice(rsts, func(i, j int) bool {
		return rsts[i].Id < rsts[j].Id
	})

	ordered := make([]Result, len(rsts))
	for i, pair := range rsts {
		ordered[i] = pair.Val
	}

	return ordered
}

func (w *SortedWaitGroup) Clear() {
	w.lck = sync.Mutex{}
	w.id = 0
	w.wg = sync.WaitGroup{}
	w.rst = nil
}

func NewSortedWaitGroup() *SortedWaitGroup {
	w := &SortedWaitGroup{}
	w.Clear()
	return w
}
