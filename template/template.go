package template

import "fmt"

var tmplContent = fmt.Sprint(`
{{- /* struct-template */ -}}
{{if .Comment -}} // {{.Name}} {{.Comment}} {{- end}}
type {{.Name}} struct { {{range .Fields}}
    {{.Key}} {{.Typ}} {{template "tags" .Tags}} {{if .Comment}}// {{.Comment}}{{end}}
{{- end}}
}

{{- /* tags-template */ -}}
{{define "tags" -}}
{{if . -}}
` + "`{{range $i, $v := . -}}\n" +
	"{{if ne $i 0}} {{end -}}\n" +
	"{{$v.Typ}}:\"{{$v.Content}}\"\n" +
	"{{- end}}`\n" +
	"{{- end}}\n" +
	"{{- end}}")
