package json

import (
	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/destination"
	template2 "github.com/jakseer/any2struct/template"
)

// the struct with json tag

const tagType = "json"

var _ destination.Destination = &Destination{}

type Destination struct{}

func New() *Destination {
	return &Destination{}
}

func (d Destination) Convert(s *template2.Struct) *template2.Struct {
	for i := range s.Fields {
		tagContent := strcase.ToSnake(s.Fields[i].Key)
		s.Fields[i].Tags = append(s.Fields[i].Tags, template2.StructFieldTag{
			Typ:     tagType,
			Content: tagContent,
		})
	}

	return s
}
