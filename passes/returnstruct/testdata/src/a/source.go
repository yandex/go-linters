package a

import "io"

type Copier interface {
	Copy() any
}

var _ Copier = (*User)(nil)
var _ io.Closer = (*User)(nil)

type User struct{}

func (u *User) Copy() any {
	return *u
}

func (u *User) Close() error {
	return nil
}

func ReturnNothing() {}

func ReturnBaseType() int {
	return 42
}

func ReturnStruct() *User {
	return new(User)
}

func ReturnInterface() Copier { // want `function must return concrete type, not interface a.Copier`
	return new(User)
}

func ReturnExternalInterface() io.Closer { // want `function must return concrete type, not interface io.Closer`
	return new(User)
}
