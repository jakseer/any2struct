package template

type TmplStruct struct {
	StructName    string
	StructComment string
	Fields        []TmplStructField
}

type TmplStructField struct {
	Key  string
	Typ  string
	Tags []TmplStructFieldTag
	Comment string
}

type TmplStructFieldTag struct {
	Typ     string
	Content string
}
