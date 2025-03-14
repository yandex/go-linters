package a

func okOneNamedValueFilled() (r int) {
	return 1
}

func okOneNamedValueFilled2() (r int) {
	r = 2
	return r
}

func okOneNamedValueFilled3() (r int) {
	return r
}

func okOneNamedValueFilled4(x int) (r int) {
	if x > 0 {
		r = x
	}
	return r
}

func okTwoNamedValuesFilled() (r1 int, r2 int) {
	r1 = 0
	return r1, r2
}
