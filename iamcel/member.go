package iamcel

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"go.einride.tech/iam/iammember"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// MemberFunction is the name of the CEL member function.
const MemberFunction = "member"

const memberFunctionOverload = "member_caller_string_string"

// NewMemberFunctionDeclaration creates a new declaration for the member function.
func NewMemberFunctionDeclaration() *expr.Decl {
	return decls.NewFunction(
		MemberFunction,
		decls.NewInstanceOverload(
			memberFunctionOverload,
			[]*expr.Type{
				decls.NewObjectType(string((&iamv1.Caller{}).ProtoReflect().Descriptor().FullName())),
				decls.String,
			},
			decls.String,
		),
	)
}

// NewMemberFunctionImplementation creates a new implementation for the member function.
func NewMemberFunctionImplementation() *functions.Overload {
	return &functions.Overload{
		Operator: memberFunctionOverload,
		Binary: func(callerVal, kindVal ref.Val) ref.Val {
			caller, ok := callerVal.Value().(*iamv1.Caller)
			if !ok {
				return types.NewErr("test: unexpected type of arg 1, expected %T but got %T", &iamv1.Caller{}, callerVal.Value())
			}

			kind, ok := kindVal.Value().(string)
			if !ok {
				return types.NewErr("test: unexpected type of arg 2, expected string but got %T", kindVal.Value())
			}

			for _, member := range caller.GetMembers() {
				memberKind, memberValue, ok := iammember.Parse(member)
				if !ok {
					return types.NewErr("member: error parsing caller member '%s'", member)
				}
				if memberKind == kind {
					return types.String(memberValue)
				}
			}

			return types.NewErr("member: no kind '%s' found in caller", kind)
		},
	}
}
