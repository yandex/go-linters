package phrasebook

var sName string   // want "identifier 'sName' contains non-ideomatic notation"
var iCount int     // want "identifier 'iCount' contains non-ideomatic notation"
var bIsReady bool  // want "identifier 'bIsReady' contains non-ideomatic notation"
var fPrice float64 // want "identifier 'fPrice' contains non-ideomatic notation"

const iMaxLimit int = 100 // want "identifier 'iMaxLimit' contains non-ideomatic notation"

type aUsers []string // want "identifier 'aUsers' contains non-ideomatic notation"

var name string
var count int
var isReady bool
var price float64

const maxLimit int = 100

type users []string

var iDx bool
var stringer any

type IReader interface { // want "identifier 'IReader' contains non-ideomatic notation"
	Read()
}

type TUser struct { // want "identifier 'TUser' contains non-ideomatic notation"
	Name string
}

type Item struct {
	ID int
}

type Image interface {
	Draw()
}

type Table struct {
	Rows int
}

type Transaction interface {
	Commit()
}

func DoSomething() {
	sGreeting := "I will not buy this record, it is scratched" // want "identifier 'sGreeting' contains non-ideomatic notation"
	iCounter := 42                                             // want "identifier 'iCounter' contains non-ideomatic notation"
	fRatio := 3.14                                             // want "identifier 'fRatio' contains non-ideomatic notation"

	message := "My hovercraft is full of eels"
	counter := 42
	ratio := 3.14

	_ = sGreeting
	_ = iCounter
	_ = fRatio
	_ = message
	_ = counter
	_ = ratio
}
