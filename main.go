package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/ast"
	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/parser"
)

const (
	StructTPL = `type %s struct {
%s
}`

	FieldTPL     = "  %s %s `json:\"%s\""
	ListFieldTPL = "  %s []%s `json:\"%s\""
	EnumTypeDefTPL = "type %s %s"
	EnumDefConstTPL = "const %s %s = \"%s\""
)

var GQLTypesToGoTypes = map[string]string{
	"Int":     "int64",
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"ID":      "string",
}

func main() {
	schemaFile := flag.String("schema", "", "schema")
	flag.Parse()

	of, err := ioutil.ReadFile(*schemaFile)
	if err != nil {
		panic(err)
	}

	doc, _ := parser.ParseSchema(&ast.Source{
		Name:    *schemaFile,
		Input:   string(of),
		BuiltIn: false,
	})

	//remap enums
	enumMap := map[string]string{}
	for _, i := range doc.Definitions {
		if i.Kind == ast.Enum {
			enumTypeName := fmt.Sprintf("%sEnum", i.Name)
			enumMap[i.Name] = enumTypeName
			fmt.Printf(EnumTypeDefTPL, enumTypeName, "string")
			fmt.Println()
			for _, e := range i.EnumValues {
				fmt.Printf(EnumDefConstTPL, e.Name, enumTypeName, e.Name)
				fmt.Println()
			}
		}
	}

	for _, i := range doc.Definitions {
		var fields []string
		if i.Kind == ast.Object || i.Kind == ast.InputObject {
			for _, f := range i.Fields {
				typeName := resolveType(f.Type.Name(), enumMap, f.Type.NonNull)
				fieldName := strings.Title(f.Name)
				jsonFieldName := f.Name
				if f.Type.Elem != nil { // list type
					elemTypeName := resolveType(f.Type.Elem.Name(), enumMap, f.Type.Elem.NonNull)
					fields = append(fields, fmt.Sprintf(ListFieldTPL, fieldName, elemTypeName, jsonFieldName))
				} else {
					fields = append(fields, fmt.Sprintf(FieldTPL, fieldName, typeName, jsonFieldName))
				}
			}
			fmt.Printf(StructTPL, i.Name, strings.Join(fields, "\n"))
		}
		fmt.Println()
	}
}

func resolveType(typeName string, enumMap map[string]string, notNull bool) string {
	if tName, hasType := GQLTypesToGoTypes[typeName]; hasType {
		typeName = tName
	}
	if tName, hasEnumType := enumMap[typeName]; hasEnumType {
		typeName = tName
	}
	if !notNull { // if type can be nullable, use pointer
		typeName = strings.Join([]string{"*", typeName}, "")
	}
	return typeName
}
