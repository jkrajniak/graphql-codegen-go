package validator

import (
	"github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/ast"
	. "github.com/jkrajniak/graphql-codegen-go/internal/gqlparser/validator"
)

func init() {
	AddRule("KnownFragmentNames", func(observers *Events, addError AddErrFunc) {
		observers.OnFragmentSpread(func(walker *Walker, fragmentSpread *ast.FragmentSpread) {
			if fragmentSpread.Definition == nil {
				addError(
					Message(`Unknown fragment "%s".`, fragmentSpread.Name),
					At(fragmentSpread.Position),
				)
			}
		})
	})
}