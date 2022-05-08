package destination

import (
	template2 "github.com/jakseer/any2struct/template"
)

// convert GoStruct to specified struct string

type Destination interface {
	// Convert generate struct tag according convert.Struct.Fields
	Convert(*template2.Struct) *template2.Struct
}
