package any2struct

import (
	"errors"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/convert"
	"github.com/jakseer/any2struct/destination"
	"github.com/jakseer/any2struct/destination/gorm"
	"github.com/jakseer/any2struct/destination/json"
	yaml2 "github.com/jakseer/any2struct/destination/yaml"
	"github.com/jakseer/any2struct/source"
	json2 "github.com/jakseer/any2struct/source/json"
	"github.com/jakseer/any2struct/source/sql"
	"github.com/jakseer/any2struct/source/yaml"

	template2 "github.com/jakseer/any2struct/template"
)

const (
	DecodeTypeSQL  string = "decode_sql"
	DecodeTypeJSON string = "decode_json"
	DecodeTypeYaml string = "decode_yaml"

	EncodeTypeJSON string = "encode_json"
	EncodeTypeGorm string = "encode_gorm"
	EncodeTypeYaml string = "encode_yaml"
)

var (
	// ErrInvalidDecodeType invalid decode type
	ErrInvalidDecodeType = errors.New("invalid decode type")

	// ErrInvalidEncodeType invalid encode type
	ErrInvalidEncodeType = errors.New("invalid encode type")

	// ErrDecode decode error
	ErrDecode = errors.New("decode error")
)

const (
	defaultTmpPath = "./template/struct.tmpl"
)

type Convertor struct {
	classNameExist map[string]struct{} // judge the class name whether has been used
	classNameMap   map[*convert.Struct]string
	tmplPath       string
}

func NewConvertor() *Convertor {
	return &Convertor{
		classNameMap:   make(map[*convert.Struct]string),
		classNameExist: make(map[string]struct{}),
		tmplPath:       defaultTmpPath,
	}
}

// Convert input(which is decodeType struct) to go struct(with encodeTypes tag)
func (c *Convertor) Convert(input string, decodeType string, encodeTypes []string) (string, error) {
	cs, err := c.parseInput(input, decodeType)
	if err != nil {
		return "", err
	}

	ts := c.convertStruct(cs)

	tmplStructs, err := c.buildTags(ts, encodeTypes)
	if err != nil {
		return "", err
	}

	return c.parseWithTemp(tmplStructs)
}

// parseInput parse input with decodeType decoder
func (c *Convertor) parseInput(input string, decodeType string) (*convert.Struct, error) {
	var decoder source.Source
	switch decodeType {
	case DecodeTypeSQL:
		decoder = sql.New()
	case DecodeTypeJSON:
		decoder = json2.New()
	case DecodeTypeYaml:
		decoder = yaml.New()
	default:
		return nil, ErrInvalidDecodeType
	}

	// decode input
	s, err := decoder.Convert(input)
	if err != nil {
		return nil, ErrDecode
	}

	return s, nil
}

// buildTags generate tags according encodeTypes
func (c *Convertor) buildTags(input []*template2.Struct, encodeTypes []string) ([]*template2.Struct, error) {
	var encoders []destination.Destination
	for _, v := range encodeTypes {
		switch v {
		case EncodeTypeJSON:
			encoders = append(encoders, json.New())
		case EncodeTypeGorm:
			encoders = append(encoders, gorm.New())
		case EncodeTypeYaml:
			encoders = append(encoders, yaml2.New())
		default:
			return nil, ErrInvalidEncodeType
		}
	}

	// build tag and generate output with template
	for k := range input {
		for _, encoder := range encoders {
			input[k] = encoder.Convert(input[k])
		}
	}

	return input, nil
}

// convertStruct convert *convert.Struct to []*template2.Struct.
// Spreading multi-level struct to one-level array, but not copy tags
func (c *Convertor) convertStruct(input *convert.Struct) []*template2.Struct {
	var ret []*template2.Struct

	var needProcessStructList []*convert.Struct
	needProcessStructList = append(needProcessStructList, input)
	c.registerClassName(input)

	for len(needProcessStructList) > 0 {
		// pop from queue
		s := needProcessStructList[0]
		needProcessStructList = needProcessStructList[1:]

		var fieldList []template2.StructField
		for _, field := range s.Fields {
			typString := string(field.Typ.Typ)

			// rename and push nested struct into queue for subsequent process
			if field.Typ.Typ == convert.StructTyp && field.Typ.Ptr != nil {
				// generate unique class name and get it
				c.registerClassName(field.Typ.Ptr)
				typString, _ = c.classNameMap[field.Typ.Ptr]

				// typString should be camel format because it is also class name
				typString = strcase.ToCamel(typString)

				field.Typ.Ptr.Name = typString
				needProcessStructList = append(needProcessStructList, field.Typ.Ptr)
			}

			fieldList = append(fieldList, template2.StructField{
				Key:     strcase.ToCamel(field.Key),
				Typ:     typString,
				Comment: field.Comment,
			})
		}

		// copy whole struct
		ret = append(ret, &template2.Struct{
			Name:    strcase.ToCamel(s.Name),
			Comment: s.Comment,
			Fields:  fieldList,
		})
	}

	return ret
}

// registerClassName  register struct ptr, it will generate unique class name with it
func (c *Convertor) registerClassName(p *convert.Struct) {
	className := strcase.ToCamel(p.Name)

	_, ok := c.classNameExist[className]
	for ; ok; _, ok = c.classNameExist[className] {
		// if duplicated, rename it with _1 postfix
		className = className + "_1"
	}
	c.classNameMap[p] = className
	c.classNameExist[className] = struct{}{}
}

// parseWithTemp generate Go Struct
func (c *Convertor) parseWithTemp(ss []*template2.Struct) (string, error) {
	tmplResp := make([]string, len(ss))
	for _, s := range ss {
		resp, err := s.ParseWithTmpl(c.tmplPath)
		if err != nil {
			return "", err
		}

		tmplResp = append(tmplResp, resp)
	}

	ret := strings.Join(tmplResp, "\n")
	ret = strings.Trim(ret, "\n")

	return ret, nil
}
