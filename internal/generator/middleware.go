package generator

import (
	"fmt"

	"go.einride.tech/authorization-aip/authorization"
	authorizationv1 "go.einride.tech/authorization-aip/proto/gen/einride/authorization/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

type Middleware struct {
	Service *protogen.Service
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
	iamPolicyServer := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/genproto/googleapis/iam/v1",
		GoName:       "IAMPolicyServer",
	})
	authorizationV1Caller := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/authorization-aip/proto/gen/einride/authorization/v1",
		GoName:       "Caller",
	})
	authorizationV1Policy := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/authorization-aip/proto/gen/einride/authorization/v1",
		GoName:       "Policy",
	})
	newPolicyProgram := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "go.einride.tech/authorization-aip/authorization",
		GoName:       "NewPolicyProgram",
	})
	f.P()
	f.P("type ", g.Service.GoName, "AuthorizationMiddleware struct {")
	f.P("next ", g.ServerGoName())
	f.P("callerFn func(", contextContext, ") (", authorizationV1Caller, ", error)")
	for _, method := range g.Service.Methods {
		switch getPolicy(method).GetDecisionPoint() {
		case authorizationv1.PolicyDecisionPoint_BEFORE, authorizationv1.PolicyDecisionPoint_AFTER:
			f.P("program", method.GoName, " ", celProgram)
		}
	}
	f.P("}")
	f.P("")
	f.P("var (")
	f.P("_ ", g.ServerGoName(), " = &", g.GoName(), "{}")
	f.P("_ ", iamPolicyServer, " = &", g.GoName(), "{}")
	f.P(")")
	f.P()
	f.P("func New", g.GoName(), "(")
	f.P("next ", g.ServerGoName(), ",")
	f.P("callerFn func(", contextContext, ") (", authorizationV1Caller, ", error),")
	f.P(") (")
	f.P("_ *", g.GoName(), ", err error,")
	f.P(") {")
	f.P("m := ", g.GoName(), "{")
	f.P("next: next,")
	f.P("callerFn: callerFn,")
	f.P("}")
	for _, method := range g.Service.Methods {
		switch policy := getPolicy(method); policy.GetDecisionPoint() {
		case authorizationv1.PolicyDecisionPoint_BEFORE, authorizationv1.PolicyDecisionPoint_AFTER:
			env, err := authorization.NewPolicyEnv(method.Input.Desc, method.Output.Desc, policy)
			if err != nil {
				panic(fmt.Errorf("%s: %v", method.Desc.FullName(), err))
			}
			if _, issues := env.Compile(policy.Expression); issues.Err() != nil {
				panic(fmt.Errorf("%s: %v", method.Desc.FullName(), issues.Err()))
			}
			f.P("m.program", method.GoName, ", err = ", newPolicyProgram, "(")
			f.P("(&", method.Input.GoIdent, "{}).ProtoReflect().Descriptor(),")
			f.P("(&", method.Output.GoIdent, "{}).ProtoReflect().Descriptor(),")
			f.P("&", authorizationV1Policy, "{")
			f.P("Description: `", policy.Description, "`,")
			f.P("Expression: `", policy.Expression, "`,")
			f.P("DecisionPoint: ", policy.DecisionPoint.Number(), ",")
			f.P("Permission: `", policy.Permission, "`,")
			f.P("},")
			f.P("next,")
			f.P(")")
			f.P("if err != nil {")
			f.P("return nil, err")
			f.P("}")
		}
	}
	f.P("return &m, nil")
	f.P("}")
	f.P()
	f.P("func (m *", g.GoName(), ") mustEmbedUnimplemented", g.Service.GoName, "Server() {}")
}

func getPolicy(method *protogen.Method) *authorizationv1.Policy {
	return proto.GetExtension(method.Desc.Options(), authorizationv1.E_Policy).(*authorizationv1.Policy)
}
