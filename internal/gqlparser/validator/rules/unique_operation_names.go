package validator

import (
	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/ast"
	. "github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/validator"
)

func init() {
	AddRule("UniqueOperationNames", func(observers *Events, addError AddErrFunc) {
		seen := map[string]bool{}

		observers.OnOperation(func(walker *Walker, operation *ast.OperationDefinition) {
			if seen[operation.Name] {
				addError(
					Message(`There can be only one operation named "%s".`, operation.Name),
					At(operation.Position),
				)
			}
			seen[operation.Name] = true
		})
	})
}
