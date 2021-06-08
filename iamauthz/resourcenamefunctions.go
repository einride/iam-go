package iamauthz

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"go.einride.tech/aip/resourcename"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

const (
	ancestorFunction         = "ancestor"
	ancestorFunctionOverload = "ancestor_string_string"
)

type ResourceNameFunctions struct{}

func (ResourceNameFunctions) Declarations() []*expr.Decl {
	return []*expr.Decl{
		decls.NewFunction(
			ancestorFunction,
			decls.NewOverload(
				ancestorFunctionOverload,
				[]*expr.Type{decls.String, decls.String},
				decls.String,
			),
		),
	}
}

func (f ResourceNameFunctions) Functions() []*functions.Overload {
	return []*functions.Overload{
		{Operator: ancestorFunctionOverload, Binary: f.ancestorStringString},
	}
}

func (ResourceNameFunctions) ancestorStringString(nameVal ref.Val, patternVal ref.Val) ref.Val {
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
}
