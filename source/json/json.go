package json

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/convert"
)

var (
	// ErrInvalidJSONType invalid json type
	ErrInvalidJSONType = errors.New("invalid json type")
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

	ret := &convert.Struct{
		Name:    defaultStructName,
		Comment: "",
		Fields:  nil,
	}

	for k, v := range m {
		if field, err := parseJSONField(k, v); err == nil {
			ret.Fields = append(ret.Fields, *field)
		}
	}

	return ret, nil
}

func parseJSONField(key string, val interface{}) (*convert.StructField, error) {
	var fieldType convert.FieldType

	switch val.(type) {
	case json.Number:
		n, _ := val.(json.Number)
		if _, err := n.Int64(); err == nil {
			fieldType = convert.Int64
			break
		}
		if _, err := n.Float64(); err == nil {
			fieldType = convert.Float64
			break
		}
		fieldType = convert.String
	case bool:
		fieldType = convert.Bool
	case []interface{}:
		return nil, ErrInvalidJSONType
	case map[string]interface{}:
		return nil, ErrInvalidJSONType
	default:
		return nil, ErrInvalidJSONType
	}
	return &convert.StructField{
		Key: strcase.ToCamel(key),
		Typ: fieldType,
	}, nil
}
