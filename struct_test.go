package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewStruct(t *testing.T) {
	var (
		testFilePath = filepath.Join("testdata", "user.go")
		tagName      = "custom_tag"
		customName   = "user"
		need         = &Struct{
			Name:       "User",
			CustomName: customName,
			Fields: []*Field{
				{
					Name:  "Name",
					Value: "name",
				},
				{
					Name:  "Login",
					Value: "username",
				},
				{
					Name:  "AuthType",
					Value: "auth_type",
				},
			},
		}
	)
	node, err := parser.ParseFile(token.NewFileSet(), testFilePath, nil, 0)
	if err != nil {
		t.Fatalf("can't parse file: %v", err)
	}
	var got *Struct
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				t.Fatalf("spec is not *ast.TypeSpec")
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("spec type is not *ast.StructType")
			}
			got = newStruct(typeSpec, structType, tagName, &options{
				customName: customName,
			})
		}
	}
	if !reflect.DeepEqual(need, got) {
		t.Fatalf("got struct is not equal to needed")
	}
}

func TestTagValue(t *testing.T) {
	var (
		testFilePath = filepath.Join("testdata", "user.go")
		tagName      = "custom_tag"
		testCases    = []struct {
			FieldName string
			NeedValue string
			Error     bool
		}{
			{
				FieldName: "Name",
				NeedValue: "name",
			},
			{
				FieldName: "Login",
				NeedValue: "username",
			},
			{
				FieldName: "Password",
				Error:     true,
			},
			{
				FieldName: "CustomField",
				Error:     true,
			},
			{
				FieldName: "CustomField1",
				Error:     true,
			},
			{
				FieldName: "AuthType",
				NeedValue: "auth_type",
			},
		}
	)
	node, err := parser.ParseFile(token.NewFileSet(), testFilePath, nil, 0)
	if err != nil {
		t.Fatalf("can't parse file: %v", err)
	}
	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				t.Fatalf("spec is not *ast.TypeSpec")
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				t.Fatalf("spec type is not *ast.StructType")
			}
			for i, field := range structType.Fields.List {
				var (
					tag      = field.Tag
					index    = i
					caseName = fmt.Sprintf("tag from field %s", testCases[i].FieldName)
				)
				t.Run(caseName, func(t *testing.T) {
					if tag == nil {
						t.Fatal("field tag is nil")
					}
					value, err := tagValue(tag, tagName)
					if testCases[index].Error {
						if err != nil {
							return
						}
						t.Fatalf("need error but it is nil")
					}
					if err != nil {
						t.Fatalf("getting tag value: %v", err)
					}
					if testCases[index].NeedValue != value {
						t.Errorf("tag value must be: %s, got: %s", testCases[i].NeedValue, value)
					}
				})
			}
		}
	}
}
