package conc

func FuncWrap[X any](f func(X) any, x X) func() any {
	return func() any {
		return f(x)
	}
}

func FuncWrapT[X any, Y any](f func(X) Y, x X) func() Y {
	return func() Y {
		return f(x)
	}
}
