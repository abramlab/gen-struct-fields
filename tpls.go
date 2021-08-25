package main

import (
	"sync"
	"text/template"
)

type onceTemplate struct {
	tpl  *template.Template
	once sync.Once
}

type Template struct {
	onceGen *onceTemplate
	tpls    []*template.Template
}

var (
	basicTemplates = &Template{
		tpls: []*template.Template{
			structNameTpl,
			structFieldsTpl,
			structFieldsArrayTpl,
		},
	}
	genjiTemplates = &Template{
		onceGen: &onceTemplate{
			tpl: genjiMainTpl,
		},
		tpls: []*template.Template{
			genjiCastTpl,
		},
	}
)

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

var (
	genjiMainTpl = template.Must(template.New("genjiMainTpl").Parse(
		`type castHelper struct{ field string }

		func (w castHelper) AsBool() string   { return "(CAST " + w.field + " AS BOOL)" }
		func (w castHelper) AsInt() string    { return "(CAST " + w.field + " AS INTEGER)" }
		func (w castHelper) AsDouble() string { return "(CAST " + w.field + " AS DOUBLE)" }
		func (w castHelper) AsBlob() string   { return "(CAST " + w.field + " AS BLOB)" }
		func (w castHelper) AsText() string   { return "(CAST " + w.field + " AS TEXT)" }
		func (w castHelper) AsArray() string  { return "(CAST " + w.field + " AS ARRAY)" }
		func (w castHelper) AsDoc() string    { return "(CAST " + w.field + " AS DOCUMENT)" }
`))

	genjiCastTpl = template.Must(template.New("genjiCastTpl").Parse(
		`var {{.Name}}CastFields = struct {
		{{range $field := .Fields -}}
		{{$field.Name}} castHelper
		{{end}}
	}{
		{{range $field := .Fields -}}
		{{$field.Name}}: castHelper{field: "{{ .Value }}"},
		{{end}}
	}
`))
)
