package iamcel

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"go.einride.tech/aip/resourcename"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// AncestorFunction is the name of the CEL ancestor function.
const AncestorFunction = "ancestor"

const ancestorFunctionOverload = "ancestor_string_string"

// NewAncestorFunctionDeclaration creates a new declaration for the ancestor function.
func NewAncestorFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		AncestorFunction,
		decls.NewOverload(
			ancestorFunctionOverload,
			[]*expr.Type{decls.String, decls.String},
			decls.String,
		),
	)
}

// NewAncestorFunctionImplementation creates a new implementation for the ancestor function.
func NewAncestorFunctionImplementation() *functions.Overload {
	return &functions.Overload{
		Operator: ancestorFunctionOverload,
		Binary: func(nameVal, patternVal ref.Val) ref.Val {
			name, ok := nameVal.Value().(string)
			if !ok {
				return types.NewErr("parent: unexpected type of arg 1, expected string but got %T", nameVal.Value())
			}
			pattern, ok := patternVal.Value().(string)
			if !ok {
				return types.NewErr("parent: unexpected type of arg 2, expected string but got %T", patternVal.Value())
			}
			result, ok := resourcename.Ancestor(name, pattern)
			if !ok {
				return types.String("")
			}
			return types.String(result)
		},
	}
}
