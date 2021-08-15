package main

import "text/template"

var basicTemplates = []*template.Template{
	structNameTpl,
	structFieldsTpl,
	structFieldsArrayTpl,
}

var (
	structNameTpl = template.Must(template.New("nameTpl").Parse(
		`const {{.Name}}Name = "{{.CustomName}}"
`))

	structFieldsTpl = template.Must(template.New("structFieldsTpl").Parse(
		`var {{.Name}}Fields = struct {
		{{range $field := .Fields -}}
		{{$field.Name}} string
		{{end}}
	}{
		{{range $field := .Fields -}}
		{{$field.Name}}: "{{ .Value }}",
		{{end}}
	}
`))

	structFieldsArrayTpl = template.Must(template.New("structFieldsArrayTpl").Parse(
		`var {{.Name}}FieldsArray = []string{
		{{range $field := .Fields -}}
		{{$.Name}}Fields.{{$field.Name}},
		{{end}}
	}
`))
)
