package require

type TestingT interface {
	Errorf(format string, args ...interface{})
}

func Equal(t TestingT, a, b interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func Equalf(t TestingT, a, b interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}

type Assertions struct{}

func New(t TestingT) *Assertions {
	return nil
}

func (*Assertions) Equal(a, b interface{}, msgAndArgs ...interface{}) {
	panic("not implemented")
}

func (*Assertions) Equalf(a, b interface{}, msg string, args ...interface{}) {
	panic("not implemented")
}
