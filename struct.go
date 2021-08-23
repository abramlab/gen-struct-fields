package main

import (
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type Struct struct {
	Name       string
	CustomName string
	Fields     []*Field
}

type Field struct {
	Name  string
	Value string
}

func newStruct(
	typeSpec *ast.TypeSpec,
	structType *ast.StructType,
	tagName string,
	opts *options,
) *Struct {
	s := &Struct{
		Name:       typeSpec.Name.Name,
		CustomName: opts.customName,
	}
	for _, field := range structType.Fields.List {
		if field.Tag != nil {
			value, err := tagValue(field.Tag, tagName)
			if err != nil {
				continue
			}
			s.Fields = append(s.Fields, &Field{
				Name:  field.Names[0].Name,
				Value: value,
			})
		}
	}
	return s
}

func tagValue(tag *ast.BasicLit, tagName string) (string, error) {
	tagPart := reflect.StructTag(tag.Value[1 : len(tag.Value)-1])
	value, ok := tagPart.Lookup(tagName)
	if !ok {
		return "", fmt.Errorf("field without %s tag", tagName)
	}

	tagValues := strings.Split(value, ",")
	if len(tagValues) == 0 || tagValues[0] == "" || tagValues[0] == "-" {
		return "", fmt.Errorf("empty or '-' %s tag value", tagName)
	}
	return tagValues[0], nil
}
