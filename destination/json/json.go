package json

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/convert"
	"github.com/jakseer/any2struct/destination"
)

// the struct with json tag

var _ destination.Destination = &Destination{}

type Destination struct{}

func New() *Destination {
	return &Destination{}
}

func (d Destination) Convert(s *convert.GoStruct) string {

	ret := fmt.Sprintf("type %s struct {\n", s.Name)
	for _, field := range s.Fields {
		s := fmt.Sprintf("%s\t%s\t\t`json:\"%s\"`\t// %s",
			strcase.ToCamel(field.Name), field.Type, field.Name, field.Comment)
		ret = ret + "\t" + s + "\n"
	}
	ret += "}"

	return ret
}
