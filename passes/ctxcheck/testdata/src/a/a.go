package a

import "context"

func CtxAsFirstArg(ctx context.Context, a int, b string) {

}

func CtxAsSecondArg(a int, ctx context.Context, b string) { // want "context parameter must be supplied as first argument of function"

}
