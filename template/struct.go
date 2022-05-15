package template

import (
	"bytes"
	"errors"
	"text/template"
)

// Struct describe the whole struct
type Struct struct {
	Name    string
	Comment string
	Fields  []StructField
}

// StructField is the field in struct
type StructField struct {
	Key     string
	Typ     string
	Tags    []StructFieldTag
	Comment string
}

// StructFieldTag is the struct tag
type StructFieldTag struct {
	Typ     string
	Content string
}

var (
	// ErrLoadTmpl load template
	ErrLoadTmpl = errors.New("load template")

	// ErrParseTmpl parse template
	ErrParseTmpl = errors.New("parse template")
)

// ParseWithTmpl parse struct with template
func (s *Struct) ParseWithTmpl(tmplPath string) (string, error) {
	tmpl, err := template.ParseFiles(tmplPath)
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
