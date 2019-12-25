package validator

import (
	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/ast"
	. "github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/validator"
)

func init() {
	AddRule("NoUnusedVariables", func(observers *Events, addError AddErrFunc) {
		observers.OnOperation(func(walker *Walker, operation *ast.OperationDefinition) {
			for _, varDef := range operation.VariableDefinitions {
				if varDef.Used {
					continue
				}

				if operation.Name != "" {
					addError(
						Message(`Variable "$%s" is never used in operation "%s".`, varDef.Variable, operation.Name),
						At(varDef.Position),
					)
				} else {
					addError(
						Message(`Variable "$%s" is never used.`, varDef.Variable),
						At(varDef.Position),
					)
				}
			}
		})
	})
}
