package validator

import (
	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/ast"
	. "github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/validator"
)

func init() {
	AddRule("NoUndefinedVariables", func(observers *Events, addError AddErrFunc) {
		observers.OnValue(func(walker *Walker, value *ast.Value) {
			if walker.CurrentOperation == nil || value.Kind != ast.Variable || value.VariableDefinition != nil {
				return
			}

			if walker.CurrentOperation.Name != "" {
				addError(
					Message(`Variable "%s" is not defined by operation "%s".`, value, walker.CurrentOperation.Name),
					At(walker.CurrentOperation.Position),
				)
			} else {
				addError(
					Message(`Variable "%s" is not defined.`, value),
					At(value.Position),
				)
			}
		})
	})
}
