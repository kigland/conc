package conc

func Pipeline(x any, fs ...func(any) (any, error)) (any, error) {
	rst := x
	var err error
	for _, f := range fs {
		rst, err = f(rst)
		if err != nil {
			return rst, err
		}
	}
	return rst, nil
}

func PipelineArr(x any, fs []func(any) (any, error)) (any, error) {
	rst := x
	var err error
	for _, f := range fs {
		rst, err = f(rst)
		if err != nil {
			return rst, err
		}
	}
	return rst, nil
}
