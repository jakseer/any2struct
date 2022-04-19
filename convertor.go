package main

import (
	"bytes"
	"errors"
	"github.com/jakseer/any2struct/destination"
	"github.com/jakseer/any2struct/destination/json"
	"github.com/jakseer/any2struct/source"
	"github.com/jakseer/any2struct/source/sql"
	"text/template"
)

const (
	EncodeTypeSQL string = "encode_sql"

	DecodeTypeJson string = "decode_json"
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

func Convert(input string, decodeType string, encodeType string) (string, error) {
	var decoder source.Source

	switch decodeType {
	case DecodeTypeJson:
		decoder = sql.New()
	default:
		return "", ErrInvalidDecodeType
	}

	var encoder destination.Destination

	switch encodeType {
	case EncodeTypeSQL:
		encoder = json.New()
	default:
		return "", ErrInvalidEncodeType
	}

	s, err := decoder.Convert(input)
	if err != nil {
		return "", ErrDecode
	}

	s = encoder.Convert(s)

	// print using template
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
