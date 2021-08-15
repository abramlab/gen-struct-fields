package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type generator struct {
	buf     bytes.Buffer
	tagName string
	genTpls []*template.Template

	pkgName string
	structs []*Struct
}

func (gen *generator) parse(needed map[string]*options) error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("reading dir entries: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".go" {
			continue
		}
		if err = gen.parseFile(entry.Name(), needed); err != nil {
			return fmt.Errorf("parsing file \"%s\": %w", entry.Name(), err)
		}
	}
	return nil
}

func (gen *generator) parseFile(path string, needed map[string]*options) error {
	node, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if _, ok := needed[typeSpec.Name.Name]; !ok {
				continue
			}
			gen.pkgName = node.Name.Name
			gen.structs = append(gen.structs, newStruct(
				typeSpec, structType, gen.tagName, needed[typeSpec.Name.Name]))
		}
	}
	return nil
}

func (gen *generator) generate() error {
	for _, s := range gen.structs {
		gen.printf("// Code generated by \"gen-struct-fields %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
		gen.printf("\n")
		gen.printf("package %s", gen.pkgName)
		gen.printf("\n")

		err := gen.printTemplates(s)
		if err != nil {
			return fmt.Errorf("generating templates: %w", err)
		}

		var (
			outputData = gen.format()
			baseName   = fmt.Sprintf("%s_%s.go", strings.ToLower(s.Name), namePostfix)
			outputName = filepath.Join(".", strings.ToLower(baseName))
		)
		gen.buf.Reset()

		tmpName := fmt.Sprintf("%s_%s_", strings.ToLower(s.Name), namePostfix)
		tmpFile, err := ioutil.TempFile(filepath.Dir("."), tmpName)
		if err != nil {
			return fmt.Errorf("creating temporary file for output: %w", err)
		}
		if _, err = tmpFile.Write(outputData); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return fmt.Errorf("writing output: %w", err)
		}
		tmpFile.Close()

		if err = os.Rename(tmpFile.Name(), outputName); err != nil {
			return fmt.Errorf("moving tempfile to output file: %w", err)
		}
	}
	return nil
}

func (gen *generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(&gen.buf, format, args...)
}

func (gen *generator) printTemplates(s *Struct) error {
	for _, tpl := range gen.genTpls {
		if err := printTemplate(&gen.buf, tpl, s); err != nil {
			return err
		}
	}
	return nil
}

func (gen *generator) format() []byte {
	src, err := format.Source(gen.buf.Bytes())
	if err != nil {
		log.Printf("WARNING: internal error: invalid Go generated: %s", err)
		log.Printf("WARNING: compile the package to analyze the error")
		return gen.buf.Bytes()
	}
	return src
}

func printTemplate(buf *bytes.Buffer, tpl *template.Template, s *Struct) error {
	if err := tpl.Execute(buf, s); err != nil {
		return fmt.Errorf("execute tpl err: %w", err)
	}
	buf.WriteString("\n")
	return nil
}
