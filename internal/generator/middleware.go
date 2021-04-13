package generator

import (
	authorizationv1 "go.einride.tech/protoc-gen-go-authorization-policy/proto/gen/einride/authorization/v1"
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

func (g Middleware) PluginGoName() string {
	return g.Service.GoName + "AuthorizationMiddlewarePlugin"
}

func (g Middleware) Generate(f *protogen.GeneratedFile) {
	contextContext := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "context",
		GoName:       "Context",
	})
	celEnvOption := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "EnvOption",
	})
	celProgramOption := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "ProgramOption",
	})
	celProgram := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "Program",
	})
	interpreterActivation := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/interpreter",
		GoName:       "Activation",
	})
	f.P()
	f.P("type ", g.Service.GoName, "AuthorizationMiddleware struct {")
	f.P("next ", g.ServerGoName())
	f.P("plugin ", g.PluginGoName())
	for _, method := range g.Service.Methods {
		switch getPolicy(method).GetDecisionPoint() {
		case authorizationv1.PolicyDecisionPoint_BEFORE, authorizationv1.PolicyDecisionPoint_AFTER:
			f.P("program", method.GoName, " ", celProgram)
		}
	}
	f.P("}")
	f.P("")
	f.P("var _ ", g.ServerGoName(), " = &", g.GoName(), "{}")
	f.P()
	f.P("type ", g.PluginGoName(), " interface {")
	f.P("EnvOptions() []", celEnvOption)
	f.P("ProgramOptions() []", celProgramOption)
	f.P("ActivationForContext(", contextContext, ") (", interpreterActivation, ", error)")
	f.P("}")
	f.P()
	f.P("func New", g.GoName(), "(")
	f.P("next ", g.ServerGoName(), ",")
	f.P("plugin ", g.PluginGoName(), ",")
	f.P(") (")
	f.P("_ *", g.GoName(), ", err error,")
	f.P(") {")
	f.P("m := ", g.GoName(), "{")
	f.P("next: next,")
	f.P("plugin: plugin,")
	f.P("}")
	for _, method := range g.Service.Methods {
		switch policy := getPolicy(method); policy.GetDecisionPoint() {
		case authorizationv1.PolicyDecisionPoint_BEFORE, authorizationv1.PolicyDecisionPoint_AFTER:
			f.P("m.program", method.GoName, ", err = m.newProgram", method.GoName, "()")
			f.P("if err != nil {")
			f.P("return nil, err")
			f.P("}")
		}
	}
	f.P("return &m, nil")
	f.P("}")
	f.P()
	f.P("func (m *", g.GoName(), ") mustEmbedUnimplemented", g.Service.GoName, "Server() {}")
	for _, method := range g.Service.Methods {
		switch policy := getPolicy(method); policy.GetDecisionPoint() {
		case authorizationv1.PolicyDecisionPoint_BEFORE, authorizationv1.PolicyDecisionPoint_AFTER:
			g.generateNewProgram(f, method, policy)
		}
	}
}

func (g Middleware) generateNewProgram(
	f *protogen.GeneratedFile,
	method *protogen.Method,
	policy *authorizationv1.Policy,
) {
	celProgram := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "Program",
	})
	celNewEnv := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "NewEnv",
	})
	celTypes := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "Types",
	})
	celDeclarations := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/cel",
		GoName:       "Declarations",
	})
	protoEqual := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "google.golang.org/protobuf/proto",
		GoName:       "Equal",
	})
	declsBool := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/checker/decls",
		GoName:       "Bool",
	})
	declsNewVar := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/checker/decls",
		GoName:       "NewVar",
	})
	declsNewObjectType := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/checker/decls",
		GoName:       "NewObjectType",
	})
	fmtErrorf := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "fmt",
		GoName:       "Errorf",
	})
	f.P()
	f.P("func (m *", g.GoName(), ") newProgram", method.GoName, "() (_ ", celProgram, ", err error) {")
	f.P("defer func() {")
	f.P("if err != nil {")
	f.P("err = ", fmtErrorf, `("authorization policy on `, method.Desc.FullName(), `: %w", err)`)
	f.P("}")
	f.P("}()")
	f.P("const expression = `", policy.Expression, "`")
	f.P("env, err := ", celNewEnv, "(append(m.plugin.EnvOptions(),")
	f.P(celTypes, "(")
	f.P("&", method.Input.GoIdent, "{},")
	if policy.DecisionPoint == authorizationv1.PolicyDecisionPoint_AFTER {
		f.P("&", method.Output.GoIdent, "{},")
	}
	f.P("),")
	f.P(celDeclarations, "(")
	f.P(declsNewVar, `("request", `, declsNewObjectType, `("`, method.Input.Desc.FullName(), `")),`)
	if policy.DecisionPoint == authorizationv1.PolicyDecisionPoint_AFTER {
		f.P(declsNewVar, `("response", `, declsNewObjectType, `("`, method.Output.Desc.FullName(), `")),`)
	}
	f.P("),")
	f.P(")...)")
	f.P("if err != nil {")
	f.P("return nil, err")
	f.P("}")
	f.P("ast, issues := env.Compile(expression)")
	f.P("if err := issues.Err(); err != nil {")
	f.P("return nil, err")
	f.P("}")
	f.P("if !", protoEqual, "(ast.ResultType(), ", declsBool, ") {")
	f.P("return nil, ", fmtErrorf, `("non-bool result type: %v", ast.ResultType())`)
	f.P("}")
	f.P("return env.Program(ast, m.plugin.ProgramOptions()...)")
	f.P("}")
}

func getPolicy(method *protogen.Method) *authorizationv1.Policy {
	return proto.GetExtension(method.Desc.Options(), authorizationv1.E_Policy).(*authorizationv1.Policy)
}
