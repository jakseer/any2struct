package convert

// Struct represent the whole struct
type Struct struct {
	Name    string
	Comment string
	Fields  []StructField
}

// StructField is the field in struct
type StructField struct {
	Key     string
	Typ     FieldTyp
	Tags    []StructFieldTag
	Comment string
}

// StructFieldTag is the struct tag
type StructFieldTag struct {
	Typ     string
	Content string
}

type fieldTypConst string

// FieldTyp is the struct field type
type FieldTyp struct {
	Ptr *Struct // if typ is struct, ptr point the struct; otherwise nil
	Typ fieldTypConst
}

const (
	Unknown   fieldTypConst = "unknown"
	Int       fieldTypConst = "int"
	Int64     fieldTypConst = "int64"
	Float64   fieldTypConst = "float64"
	String    fieldTypConst = "string"
	Bool      fieldTypConst = "bool"
	StructTyp fieldTypConst = "struct"
	Array     fieldTypConst = "[]interface"
)
