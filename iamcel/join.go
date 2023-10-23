package iamcel

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"go.einride.tech/aip/resourcename"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// JoinFunction is the name of the CEL descendant function.
const JoinFunction = "join"

const joinFunctionOverload = "join_string_string"

// NewJoinFunctionDeclaration creates a new declaration for the descendant function.
func NewJoinFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		JoinFunction,
		// TODO: if ever possible in CEL-go, declare this as a variadic function.
		decls.NewOverload(
			joinFunctionOverload,
			[]*expr.Type{decls.String, decls.String},
			decls.String,
		),
	)
}

// NewJoinFunctionImplementation creates a new implementation for the descendant function.
func NewJoinFunctionImplementation() *functions.Overload {
	return &functions.Overload{
		Operator: joinFunctionOverload,
		Binary: func(parentVal, childVal ref.Val) ref.Val {
			parent, ok := parentVal.Value().(string)
			if !ok {
				return types.NewErr("parent: unexpected type of arg 1, expected string but got %T", parentVal.Value())
			}
			child, ok := childVal.Value().(string)
			if !ok {
				return types.NewErr("child: unexpected type of arg 2, expected string but got %T", childVal.Value())
			}
			return types.String(resourcename.Join(parent, child))
		},
	}
}
