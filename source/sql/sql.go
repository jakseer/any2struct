package sql

import (
	"errors"

	"github.com/iancoleman/strcase"
	"github.com/jakseer/any2struct/convert"
	"github.com/jakseer/any2struct/source"
	"github.com/xwb1989/sqlparser"
)

// the input source with sql statement

var _ source.Source = &Source{}

var (
	// ErrInvalidDdlFormat invalid ddl format
	ErrInvalidDdlFormat = errors.New("invalid ddl format")

	// ErrNotCreateSQL not create sql
	ErrNotCreateSQL = errors.New("not create sql")
)

type Source struct{}

func New() *Source {
	return &Source{}
}

func (s Source) Convert(sql string) (*convert.Struct, error) {
	stmtTree, err := sqlparser.ParseStrictDDL(sql)
	if err != nil {
		return nil, err
	}

	ddlTree, ok := stmtTree.(*sqlparser.DDL)
	if !ok {
		return nil, ErrInvalidDdlFormat
	}

	if ddlTree.Action != "create" {
		return nil, ErrNotCreateSQL
	}

	goStruct := &convert.Struct{
		Name: ddlTree.NewName.Name.String(),
	}

	for _, field := range ddlTree.TableSpec.Columns {
		goStruct.Fields = append(goStruct.Fields, convert.StructField{
			Key:     strcase.ToCamel(field.Name.String()),
			Typ:     parseSQLType(field.Type.Type),
			Comment: string(field.Type.Comment.Val),
		})
	}

	return goStruct, nil
}

func parseSQLType(sqlType string) convert.FieldType {
	switch sqlType {
	case "int":
		return convert.Int
	case "varchar":
		return convert.String
	default:
		return convert.Unknown
	}
}
