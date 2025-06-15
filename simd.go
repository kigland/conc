package conc

func SIMD[X, Y any](f func(X) Y, xs []X) []Y {
	swg := NewSortedWaitGroup()
	for _, x := range xs {
		swg.AddTaskWithoutErr(func() any {
			return f(x)
		})
	}
	rsts := swg.Wait()
	ys := make([]Y, len(rsts))
	for i, rst := range rsts {
		ys[i] = rst.Val.(Y)
	}
	return ys
}

func MISD[X any](fs []func(X) any, x X) []any {
	swg := NewSortedWaitGroup()
	for _, f := range fs {
		swg.AddTaskWithoutErr(func() any {
			return f(x)
		})
	}
	return swg.Wait().AllValues()
}
