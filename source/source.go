package source

import "github.com/jakseer/any2struct/convert"

// parse input and generate GoStruct

type Source interface {
	Convert(string) (*convert.Struct,error)
}
