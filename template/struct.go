package template

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
