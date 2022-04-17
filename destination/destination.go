package destination

import "github.com/jakseer/any2struct/convert"

// convert GoStruct to specified struct string

type Destination interface {
	Convert(*convert.GoStruct) string
}
