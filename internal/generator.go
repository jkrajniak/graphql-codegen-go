package internal

import (
	"fmt"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
	"strings"
)

const (
	StructTPL = `type %s struct {
%s
}`

	FieldTPL        = "  %s %s `json:\"%s\""
	ListFieldTPL    = "  %s []%s `json:\"%s\""
	EnumTypeDefTPL  = "type %s %s"
	EnumDefConstTPL = "const %s %s = \"%s\""
)

var GQLTypesToGoTypes = map[string]string{
	"Int":     "int64",
	"Float":   "float64",
	"String":  "string",
	"Boolean": "bool",
	"ID":      "string",
}

type GoGenerator struct {
	entities    []string
	enumMapType map[string]string

	output Outputer
}

func NewGoGenerator(output Outputer, entities []string) *GoGenerator {
	return &GoGenerator{output: output, entities: entities}
}

func (g *GoGenerator) Generate(inputSchema string) {
	doc, _ := parser.ParseSchema(&ast.Source{
		Input:   inputSchema,
		BuiltIn: false,
	})

	//remap enums
	enumMap := map[string]string{}
	for _, i := range doc.Definitions {
		if i.Kind == ast.Enum {
			enumTypeName := fmt.Sprintf("%sEnum", i.Name)
			enumMap[i.Name] = enumTypeName
			if err := g.output.Write(fmt.Sprintf(EnumTypeDefTPL, enumTypeName, "string")); err != nil {
				panic(err)
			}
			if err := g.output.Write("\n"); err != nil {
				panic(err)
			}
			for _, e := range i.EnumValues {
				if err := g.output.Write(fmt.Sprintf(EnumDefConstTPL, e.Name, enumTypeName, e.Name)); err != nil {
					panic(err)
				}
				if err := g.output.Write("\n"); err != nil {
					panic(err)
				}
			}
		}
	}

	for _, i := range doc.Definitions {
		var fields []string
		if i.Kind == ast.Object || i.Kind == ast.InputObject {
			if g.entities != nil && len(g.entities) > 0 && !inArray(i.Name, g.entities) {
				continue
			}
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

func inArray(item string, items []string) bool {
	itemLower := strings.ToLower(item)
	for _, v := range items {
		itemVal := strings.TrimSpace(strings.ToLower(v))
		if itemVal == itemLower {
			return true
		}
	}
	return false
}