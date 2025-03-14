package a

import "context"

type MyApp struct {
	Address  string
	Threads  int
	Callback func(context.Context) error
}

type MyAppRequest struct {
	Caller string
	Params []any
	Ctx    context.Context // want "context must not be saved as a struct field"
}
