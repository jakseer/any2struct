package convert

type GoStruct struct {
	Name   string
	Fields []GoStructField
}

type GoStructField struct {
	Name    string
	Type    FieldType
	Comment string
}

type FieldType string

const (
	Unknown FieldType = "unknown"
	Int     FieldType = "int"
	String  FieldType = "string"
)
