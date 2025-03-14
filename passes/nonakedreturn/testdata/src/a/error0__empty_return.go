package a

func EEmptyReturn() (x int) {
	x = 5
	return
}

func EEmptyReturn2() (x int, y float32) {
	x = 5
	y = 0.0
	return
}

func ESomeReturnsEmpty(a int) (x int, y float32) {
	if a > 0 {
		return
	}

	x += a
	return x, y
}

func ESomeReturnsEmpty2(a int) (x int, y float32) {
	if a > 0 {
		return 0, 0
	}

	if a < 0 {
		return
	}

	x += a
	return x, y
}
