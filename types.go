package conc

type idPair struct {
	Id  int
	Val Result
}

type Result struct {
	Val any
	Err error
}

type ResultT[T any] struct {
	Val T
	Err error
}

type AllResults []Result

func (a AllResults) AllValues() []any {
	vals := make([]any, len(a))
	for i, rst := range a {
		vals[i] = rst.Val
	}
	return vals
}

func (a AllResults) AllErrors() []error {
	errors := make([]error, len(a))
	for i, rst := range a {
		errors[i] = rst.Err
	}
	return errors
}

func (a AllResults) AnyError() bool {
	for _, rst := range a {
		if rst.Err != nil {
			return true
		}
	}
	return false
}

func (a AllResults) AllSuccess() bool {
	return !a.AnyError()
}
