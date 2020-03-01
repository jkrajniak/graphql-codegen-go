package internal

import (
	"fmt"
	"github.com/jkrajniak/graphql-codegen-go/internal/readers"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
	"os"
)

type InputSchema struct {
	Data       string
	SourcePath string
}

func ReadSchemas(schemaPaths []string) ([]InputSchema, error) {
	var outs []InputSchema
	for _, s := range schemaPaths {
		r := readers.DiscoverReader(s)
		o, err := r.Read()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read from %s", s)
		}
		outs = append(outs, InputSchema{
			Data:       string(o),
			SourcePath: s,
		})
	}
	return outs, nil
}

func LoadSchemas(inputSchemas []InputSchema) (*ast.SchemaDocument, error) {
	sourceSchemas := []*ast.Source{validator.Prelude}  // include types
	for _, inputSchema := range inputSchemas {
		sourceSchemas = append(sourceSchemas, &ast.Source{
			Name:    inputSchema.SourcePath,
			Input:   inputSchema.Data,
			BuiltIn: false,
		})
	}
	doc, err := parser.ParseSchemas(sourceSchemas...)
	if err != nil {
		return nil, err
	}

	if _, err := validator.ValidateSchemaDocument(doc); err != nil {
		f := formatter.NewFormatter(os.Stderr)
		os.Stderr.WriteString("Parsed schema:\n")
		f.FormatSchemaDocument(doc)
		return nil, fmt.Errorf(err.Message)
	}

	return doc, nil
}
