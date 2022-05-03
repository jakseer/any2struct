package destination

import (
	"github.com/jakseer/any2struct/convert"
)

// convert GoStruct to specified struct string

type Destination interface {
	// Convert generate struct tag according convert.Struct.Fields
	Convert(*convert.Struct) *convert.Struct
}
