package nilness

func Triggers() bool {
	var test []int
	if test == nil { // want `tautological condition: nil == nil`
		return true
	}
	return false
}

func NotTriggers() bool {
	var test []int
	return test == nil
}
