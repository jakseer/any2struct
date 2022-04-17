package main

import (
	"errors"
	"github.com/jakseer/any2struct/destination"
	"github.com/jakseer/any2struct/destination/json"
	"github.com/jakseer/any2struct/source"
	"github.com/jakseer/any2struct/source/sql"
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

	// ErrDecodeError decode error
	ErrDecodeError = errors.New("decode error")
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

	gs, err := decoder.Convert(input)
	if err != nil {
		return "", ErrDecodeError
	}

	return encoder.Convert(gs), nil
}
