package generator

import (
	"fmt"

	"go.einride.tech/iam/authorization"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Middleware struct {
	Service *protogen.Service
}

func (g Middleware) PermissionTesterGoName() string {
	return "PermissionTester"
}

func (g Middleware) GoName() string {
	return g.Service.GoName + "AuthorizationMiddleware"
}

func (g Middleware) ServerGoName() string {
	return g.Service.GoName + "Server"
}

func (g Middleware) Generate(f *protogen.GeneratedFile) {
	contextContext := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "context",
		GoName:       "Context",
	})
	celProgram := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "Program",
	})
	authorizationV1Caller := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/iam/proto/gen/einride/iam/v1",
		GoName:       "Caller",
	})
	f.P()
	f.P("type ", g.Service.GoName, "AuthorizationMiddleware struct {")
	f.P("next ", g.ServerGoName())
	f.P("callerFn func(", contextContext, ") (", authorizationV1Caller, ", error)")
	f.P("permissionTester ", g.PermissionTesterGoName())
	for _, method := range g.Service.Methods {
		switch getPolicy(method).GetDecisionPoint() {
		case iamv1.PolicyDecisionPoint_BEFORE, iamv1.PolicyDecisionPoint_AFTER:
			f.P("program", method.GoName, " ", celProgram)
		}
	}
	f.P("}")
	f.P()
	f.P("var _ ", g.ServerGoName(), " = &", g.GoName(), "{}")
	f.P()
	f.P("type ", g.PermissionTesterGoName(), " interface {")
	f.P("Test(ctx ", contextContext, ", permission string, caller *", authorizationV1Caller, ", resource string) (bool, error)")
	f.P("TestAll(ctx ", contextContext, ", permission string, caller *", authorizationV1Caller, ", resources []string) (bool, error)")
	f.P("TestAny(ctx ", contextContext, ", permission string, caller *", authorizationV1Caller, ", resources []string) (bool, error)")
	f.P("}")
	f.P()
	f.P("func New", g.GoName(), "(")
	f.P("next ", g.ServerGoName(), ",")
	f.P("callerFn func(", contextContext, ") (", authorizationV1Caller, ", error),")
	f.P("permissionTester ", g.PermissionTesterGoName(), ",")
	f.P(") (")
	f.P("_ *", g.GoName(), ", err error,")
	f.P(") {")
	f.P("m := ", g.GoName(), "{")
	f.P("next: next,")
	f.P("callerFn: callerFn,")
	f.P("permissionTester: permissionTester,")
	f.P("}")
	for _, method := range g.Service.Methods {
		switch policy := getPolicy(method); policy.GetDecisionPoint() {
		case iamv1.PolicyDecisionPoint_BEFORE, iamv1.PolicyDecisionPoint_AFTER:
			env, err := authorization.NewPolicyEnv(method.Input.Desc, method.Output.Desc, policy)
			if err != nil {
				panic(fmt.Errorf("%s: %v", method.Desc.FullName(), err))
			}
			if _, issues := env.Compile(policy.Expression); issues.Err() != nil {
				panic(fmt.Errorf("%s: %v", method.Desc.FullName(), issues.Err()))
			}
			// f.P("m.program", method.GoName, ", err = ", newPolicyProgram, "(")
			// f.P("(&", method.Input.GoIdent, "{}).ProtoReflect().Descriptor(),")
			// f.P("(&", method.Output.GoIdent, "{}).ProtoReflect().Descriptor(),")
			// f.P("&", authorizationV1Policy, "{")
			// f.P("Description: `", policy.Description, "`,")
			// f.P("Expression: `", policy.Expression, "`,")
			// f.P("DecisionPoint: ", policy.DecisionPoint.Number(), ",")
			// f.P("Permission: `", policy.Permission, "`,")
			// f.P("},")
			// f.P("next,")
			// f.P(")")
			// f.P("if err != nil {")
			// f.P("return nil, err")
			// f.P("}")
		}
	}
	f.P("return &m, nil")
	f.P("}")
	f.P()
	f.P("func (m *", g.GoName(), ") mustEmbedUnimplemented", g.Service.GoName, "Server() {}")
}

func getPolicy(method *protogen.Method) *iamv1.Policy {
	return proto.GetExtension(method.Desc.Options(), iamv1.E_Policy).(*iamv1.Policy)
}
