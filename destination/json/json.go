package json

import (
	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/convert"
	"github.com/jakseer/any2struct/destination"
)

// the struct with json tag

const tagType = "json"

var _ destination.Destination = &Destination{}

type Destination struct{}

func New() *Destination {
	return &Destination{}
}

func (d Destination) Convert(s *convert.Struct) *convert.Struct {
	for i, _ := range s.Fields {
		tagContent := strcase.ToSnake(s.Fields[i].Key)
		s.Fields[i].Tags = append(s.Fields[i].Tags, convert.StructFieldTag{
			Typ:     tagType,
			Content: tagContent,
		})
	}

	return s
}