package yaml

import (
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/jakseer/any2struct/convert"
	"github.com/jakseer/any2struct/source"
)

var (
	// ErrInvalidYamlType invalid yaml type
	ErrInvalidYamlType = errors.New("invalid yaml type")

	// ErrInvalidYamlNestedStructure invalid yaml nested structure
	ErrInvalidYamlNestedStructure = errors.New("invalid yaml nested structure")
)

var _ source.Source = &Source{}

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s Source) Convert(input string) (*convert.Struct, error) {
	m := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(input), m)
	if err != nil {
		return nil, err
	}

	ss, err := parseYamlStruct(m)
	if err != nil {
		return nil, err
	}
	ss.Name = "Untitled"

	return ss, nil
}

func parseYamlStruct(m map[string]interface{}) (*convert.Struct, error) {
	ret := &convert.Struct{
		Fields: nil,
	}

	for k, v := range m {
		if field, err := parseYamlField(k, v); err == nil {
			ret.Fields = append(ret.Fields, *field)
		}
	}

	return ret, nil
}

func parseYamlField(key string, val interface{}) (*convert.StructField, error) {
	var fieldType convert.FieldTyp

	switch val.(type) {
	case bool:
		fieldType = convert.FieldTyp{
			Ptr: nil,
			Typ: convert.Bool,
		}
	case float64:
		fieldType = convert.FieldTyp{
			Ptr: nil,
			Typ: convert.Float64,
		}
	case string:
		fieldType = convert.FieldTyp{
			Ptr: nil,
			Typ: convert.String,
		}
	case int64:
		fieldType = convert.FieldTyp{
			Ptr: nil,
			Typ: convert.Int64,
		}
	case []interface{}:
		fieldType = convert.FieldTyp{
			Ptr: nil,
			Typ: convert.Array,
		}
	case map[string]interface{}:
		s, err := parseYamlStruct(val.(map[string]interface{}))
		if err != nil {
			return nil, ErrInvalidYamlNestedStructure
		}

		s.Name = key
		fieldType = convert.FieldTyp{
			Ptr: s,
			Typ: convert.StructTyp,
		}
	default:
		return nil, ErrInvalidYamlType
	}

	return &convert.StructField{
		Key: key,
		Typ: fieldType,
	}, nil
}
