package require

type TestingT interface {
	Errorf(format string, args ...any)
}

func Equal(t TestingT, a, b any, msgAndArgs ...any) {
	panic("not implemented")
}

func Equalf(t TestingT, a, b any, msg string, args ...any) {
	panic("not implemented")
}

type Assertions struct{}

func New(t TestingT) *Assertions {
	return nil
}

func (*Assertions) Equal(a, b any, msgAndArgs ...any) {
	panic("not implemented")
}

func (*Assertions) Equalf(a, b any, msg string, args ...any) {
	panic("not implemented")
}
