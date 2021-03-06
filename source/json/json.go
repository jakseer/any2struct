package json

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/jakseer/any2struct/convert"
)

var (
	// ErrInvalidJSONType invalid json type
	ErrInvalidJSONType = errors.New("invalid json type")

	// ErrInvalidJSONNestedStructure invalid json nested structure
	ErrInvalidJSONNestedStructure = errors.New("invalid json nested structure")
)

const (
	// defaultStructName default struct name
	defaultStructName = "Untitled"
)

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s Source) Convert(str string) (*convert.Struct, error) {
	m := make(map[string]interface{})
	jsonDecoder := json.NewDecoder(strings.NewReader(str))
	jsonDecoder.UseNumber()
	err := jsonDecoder.Decode(&m)
	if err != nil {
		return nil, err
	}

	ss, err := parseJSONStruct(m)
	if err != nil {
		return nil, err
	}
	ss.Name = defaultStructName

	return ss, nil
}

func parseJSONStruct(m map[string]interface{}) (*convert.Struct, error) {
	ret := &convert.Struct{
		Fields: nil,
	}
	for k, v := range m {
		if field, err := parseJSONField(k, v); err == nil {
			ret.Fields = append(ret.Fields, *field)
		}
	}

	return ret, nil
}

func parseJSONField(key string, val interface{}) (*convert.StructField, error) {
	var fieldType convert.FieldTyp

	switch val.(type) {
	case json.Number:
		n, _ := val.(json.Number)
		if _, err := n.Int64(); err == nil {
			fieldType = convert.FieldTyp{
				Typ: convert.Int64,
			}
			break
		}
		if _, err := n.Float64(); err == nil {
			fieldType = convert.FieldTyp{
				Typ: convert.Float64,
			}
			break
		}
		fieldType = convert.FieldTyp{
			Typ: convert.String,
		}
	case bool:
		fieldType = convert.FieldTyp{
			Typ: convert.Bool,
		}
	case []interface{}:
		fieldType = convert.FieldTyp{
			Typ: convert.Array,
		}
	case map[string]interface{}:
		s, err := parseJSONStruct(val.(map[string]interface{}))
		if err != nil {
			return nil, ErrInvalidJSONNestedStructure
		}

		s.Name = key
		fieldType = convert.FieldTyp{
			Ptr: s,
			Typ: convert.StructTyp,
		}
	default:
		return nil, ErrInvalidJSONType
	}
	return &convert.StructField{
		Key: key,
		Typ: fieldType,
	}, nil
}
