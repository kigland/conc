package conc

import (
	"sort"
	"sync"
)

type SortedWaitGroup struct {
	lck     sync.Mutex
	id      int
	wg      sync.WaitGroup
	sem     chan struct{}
	rst     []idPair
	maxConc int
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

	sem := w.sem
	maxConc := w.maxConc

	w.wg.Add(1)
	go func(id int, sem chan struct{}, maxConc int) {
		defer w.wg.Done()
		if maxConc > 0 {
			sem <- struct{}{}
			defer func() { <-sem }()
		}
		rst, err := f()

		w.lck.Lock()
		w.rst = append(w.rst, idPair{id, Result{Val: rst, Err: err}})
		w.lck.Unlock()
	}(id, sem, maxConc)

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
	w.initConcSem()
}

func (w *SortedWaitGroup) initConcSem() {
	if w.sem != nil {
		close(w.sem)
	}
	if w.maxConc > 0 {
		w.sem = make(chan struct{}, w.maxConc)
	} else {
		w.sem = nil
	}
}

func NewSortedWaitGroup() *SortedWaitGroup {
	w := &SortedWaitGroup{}
	w.Clear()
	return w
}

func (w *SortedWaitGroup) WithNewMaxConc(maxConc int) *SortedWaitGroup {
	w.maxConc = maxConc
	w.initConcSem()
	return w
}

func NewSortedWaitGroupMaxConc(maxConc int) *SortedWaitGroup {
	w := &SortedWaitGroup{
		maxConc: maxConc,
	}
	w.Clear()
	return w
}
