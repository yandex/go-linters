package a

type Msg struct {
	XXX_sizecache int32
}

func X() Msg { // want "X returns proto by value"
	return Msg{}
}

func Y(Msg) {} // want "Y passes proto by value"

func Z() {
	var m Msg
	_ = m // want "assignment copies proto value"
}

func NilPkg() {
	var err error
	_ = err
}
