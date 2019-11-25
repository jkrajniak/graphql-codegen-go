package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

const (
	StructTPL = `type %s struct {
%s
}`

	FieldTPL     = "  %s %s `json:\"%s\""
	ListFieldTPL = "  %s []%s `json:\"%s\""
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

	for _, i := range doc.Definitions {
		var fields []string
		if i.IsCompositeType() {
			for _, f := range i.Fields {
				typeName := f.Type.Name()
				if tName, hasType := GQLTypesToGoTypes[f.Type.Name()]; hasType {
					typeName = tName
				}
				if !f.Type.NonNull {
					typeName = strings.Join([]string{"*", typeName}, "")
				}
				fields = append(fields, fmt.Sprintf(FieldTPL, strings.Title(f.Name), typeName, f.Name))
			}
			fmt.Printf(StructTPL, i.Name, strings.Join(fields, "\n"))
		} else if i.IsLeafType() {
			for _, e := range i.EnumValues {
				fmt.Println(e.Directives)
			}
		}
		fmt.Println()
	}
}
