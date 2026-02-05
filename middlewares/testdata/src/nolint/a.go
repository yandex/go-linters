package a

// Triggers is an example function that triggers linter
// If nolint test fails on function Triggers, the test should be fixed
func Triggers(v int) bool {
	p := &v
	if p != nil { // want `tautological condition: non-nil != nil` `if you believe this report is false positive, please silence it with //nolint:nilness comment`
		return true
	}
	return false
}

// NotTriggers is a copy of Triggers function with additional
// nolint comment for disabling linter
func NotTriggers() bool {
	var test []int
	//nolint:nilness
	if test == nil {
		return true
	}
	return false
}
