package a

func okOneValueFilled() int {
	return 1
}

func okOneValueFilled2() int {
	x := 2
	return x
}

func okTwoValuesFilled() (int, int) {
	x := 3
	return 1, x
}
