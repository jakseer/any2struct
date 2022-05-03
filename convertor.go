package main

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/jakseer/any2struct/destination"
	"github.com/jakseer/any2struct/destination/gorm"
	"github.com/jakseer/any2struct/destination/json"
	"github.com/jakseer/any2struct/source"
	"github.com/jakseer/any2struct/source/sql"
)

const (
	DecodeTypeSQL string = "decode_sql"

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

func Convert(input string, decodeType string, encodeTypes []string) (string, error) {
	var decoder source.Source
	switch decodeType {
	case DecodeTypeSQL:
		decoder = sql.New()
	default:
		return "", ErrInvalidDecodeType
	}

	var encoders []destination.Destination
	for _, v := range encodeTypes {
		switch v {
		case EncodeTypeJSON:
			encoders = append(encoders, json.New())
		case EncodeTypeGorm:
			encoders = append(encoders, gorm.New())
		default:
			return "", ErrInvalidEncodeType
		}
	}

	// parse input
	s, err := decoder.Convert(input)
	if err != nil {
		return "", ErrDecode
	}

	// build tag and generate output with template
	for _, encoder := range encoders {
		s = encoder.Convert(s)
	}

	tmpl, err := template.ParseFiles("./template/struct.tmpl")
	if err != nil {
		return "", ErrLoadTmpl
	}

	b := bytes.Buffer{}
	err = tmpl.Execute(&b, s)
	if err != nil {
		return "", ErrParseTmpl
	}

	return b.String(), nil
}
