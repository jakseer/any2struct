package any2struct

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"github.com/jakseer/any2struct/convert"
	template2 "github.com/jakseer/any2struct/template"
	"github.com/jinzhu/copier"

	"github.com/jakseer/any2struct/destination"
	"github.com/jakseer/any2struct/destination/gorm"
	"github.com/jakseer/any2struct/destination/json"
	"github.com/jakseer/any2struct/source"
	json2 "github.com/jakseer/any2struct/source/json"
	"github.com/jakseer/any2struct/source/sql"
)

const (
	DecodeTypeSQL  string = "decode_sql"
	DecodeTypeJSON string = "decode_json"

	EncodeTypeJSON string = "encode_json"
	EncodeTypeGorm string = "encode_gorm"
)

var (
	// ErrInvalidDecodeType invalid decode type
	ErrInvalidDecodeType = errors.New("invalid decode type")

	// ErrInvalidEncodeType invalid encode type
	ErrInvalidEncodeType = errors.New("invalid encode type")

	// ErrDecode decode error
	ErrDecode = errors.New("decode error")

	// ErrLoadTmpl load template
	ErrLoadTmpl = errors.New("load template")

	// ErrParseTmpl parse template
	ErrParseTmpl = errors.New("parse template")
)

// Convert input(which is decodeType struct) to go struct(with encodeTypes tag)
func Convert(input string, decodeType string, encodeTypes []string) (string, error) {
	cs, err := parseInput(input, decodeType)
	if err != nil {
		return "", err
	}

	ts := convertStruct(cs)

	tmplStructs, err := buildTags(ts, encodeTypes)
	if err != nil {
		return "", err
	}

	return parseTemplate(tmplStructs)
}

func parseInput(input string, decodeType string) (*convert.Struct, error) {
	var decoder source.Source
	switch decodeType {
	case DecodeTypeSQL:
		decoder = sql.New()
	case DecodeTypeJSON:
		decoder = json2.New()
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

func buildTags(input []*template2.Struct, encodeTypes []string) ([]*template2.Struct, error) {
	var encoders []destination.Destination
	for _, v := range encodeTypes {
		switch v {
		case EncodeTypeJSON:
			encoders = append(encoders, json.New())
		case EncodeTypeGorm:
			encoders = append(encoders, gorm.New())
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

// convertStruct convert *convert.Struct to []*template2.Struct. Spreading multi-level struct to one-level array
func convertStruct(input *convert.Struct) []*template2.Struct {
	var ret []*template2.Struct

	var needProcessStructList []*convert.Struct
	needProcessStructList = append(needProcessStructList, input)

	for len(needProcessStructList) > 0 {
		// pop from queue
		s := needProcessStructList[0]
		needProcessStructList = needProcessStructList[1:]

		var fieldList []template2.StructField
		for _, field := range s.Fields {
			typString := string(field.Typ.Typ)

			// push nested struct into queue for subsequent process
			if field.Typ.Typ == convert.StructTyp && field.Typ.Ptr != nil {
				needProcessStructList = append(needProcessStructList, field.Typ.Ptr)
			}

			// copy tag list
			var tagList []template2.StructFieldTag
			if err := copier.Copy(&tagList, &field.Tags); err != nil {
				continue
			}

			fieldList = append(fieldList, template2.StructField{
				Key:     field.Key,
				Typ:     typString,
				Tags:    tagList,
				Comment: field.Comment,
			})
		}

		// copy whole struct
		ret = append(ret, &template2.Struct{
			Name:    s.Name,
			Comment: s.Comment,
			Fields:  fieldList,
		})
	}

	return ret
}

// parseTemplate
func parseTemplate(ss []*template2.Struct) (string, error) {
	tmplResp := make([]string, len(ss))
	for _, s := range ss {
		tmpl, err := template.ParseFiles("./template/struct.tmpl")
		if err != nil {
			return "", ErrLoadTmpl
		}

		b := bytes.Buffer{}
		err = tmpl.Execute(&b, s)
		if err != nil {
			return "", ErrParseTmpl
		}

		tmplResp = append(tmplResp, b.String())
	}

	ret := strings.Join(tmplResp, "\n")
	ret = strings.Trim(ret, "\n")

	return ret, nil
}
