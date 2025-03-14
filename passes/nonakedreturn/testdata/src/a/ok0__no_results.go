package a

func okNoReturns() {
	print(0)
}

func okNoResultsOneReturn(x int) {
	if x < 0 {
		return
	}

	print(x)
}

func okNoResultsTwoReturns(x int) {
	if x < 0 {
		return
	}

	if x == 0 {
		print("zero")
		return
	}

	print(x)
}
