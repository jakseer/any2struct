package convert

type Struct struct {
	Name    string
	Comment string
	Fields  []StructField
}

type StructField struct {
	Key     string
	Typ     FieldType
	Tags    []StructFieldTag
	Comment string
}

type StructFieldTag struct {
	Typ     string
	Content string
}

type FieldType string

const (
	Unknown FieldType = "unknown"
	Int     FieldType = "int"
	String  FieldType = "string"
)
