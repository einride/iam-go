package generator

import (
	authorizationv1 "go.einride.tech/protoc-gen-go-authorization-policy/proto/gen/einride/authorization/v1"
	"google.golang.org/protobuf/compiler/protogen"
)

type MiddlewareMethod struct {
	Method *protogen.Method
	Policy *authorizationv1.Policy
}

func (g MiddlewareMethod) Generate(f *protogen.GeneratedFile) {
	middleware := Middleware{Service: g.Method.Parent}
	contextContext := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "context",
		GoName:       "Context",
	})
	interpreterNewActivation := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/interpreter",
		GoName:       "NewActivation",
	})
	interpreterNewHierarchicalActivation := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "github.com/google/cel-go/interpreter",
		GoName:       "NewHierarchicalActivation",
	})
	f.P()
	f.P("func (m *", middleware.GoName(), ") ", g.Method.GoName, "(")
	f.P("ctx ", contextContext, ",")
	f.P("request *", g.Method.Input.GoIdent, ",")
	f.P(") ( ")
	f.P("*", g.Method.Output.GoIdent, ", error,")
	f.P(") {")
	switch g.Policy.GetDecisionPoint() {
	case authorizationv1.PolicyDecisionPoint_BEFORE:
		f.P("contextActivation, err := m.plugin.ActivationForContext(ctx)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("requestActivation, err := ", interpreterNewActivation, "(map[string]interface{}{")
		f.P(`"request": request,`)
		f.P("})")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("activation := ", interpreterNewHierarchicalActivation, "(contextActivation, requestActivation)")
		f.P("val, _, err := m.program", g.Method.GoName, ".Eval(activation)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("boolVal, ok := val.Value().(bool)")
		f.P("if !ok {")
		f.P("return nil, nil // TODO: Return error.")
		f.P("}")
		f.P("if !boolVal {")
		f.P("return nil, nil // TODO: Return error.")
		f.P("}")
		f.P("")
		f.P("return m.next.", g.Method.GoName, "(ctx, request)")
	case authorizationv1.PolicyDecisionPoint_AFTER:
		f.P("response, err := m.next.", g.Method.GoName, "(ctx, request)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("contextActivation, err := m.plugin.ActivationForContext(ctx)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("responseActivation, err := ", interpreterNewActivation, "(map[string]interface{}{")
		f.P(`"request": request,`)
		f.P(`"response": response,`)
		f.P("})")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("activation := ", interpreterNewHierarchicalActivation, "(contextActivation, responseActivation)")
		f.P("val, _, err := m.program", g.Method.GoName, ".Eval(activation)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("boolVal, ok := val.Value().(bool)")
		f.P("if !ok {")
		f.P("return nil, nil // TODO: Return error.")
		f.P("}")
		f.P("if !boolVal {")
		f.P("return nil, nil // TODO: Return error.")
		f.P("}")
		f.P("return response, nil")
	// TODO: Implement CUSTOM.
	default:
		f.P("return m.next.", g.Method.GoName, "(ctx, request)")
	}
	f.P("}")
}
