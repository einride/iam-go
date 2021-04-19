package generator

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/compiler/protogen"
)

type MiddlewareMethod struct {
	Method *protogen.Method
	Policy *iamv1.Policy
}

func (g MiddlewareMethod) Generate(f *protogen.GeneratedFile) {
	middleware := Middleware{Service: g.Method.Parent}
	contextContext := f.QualifiedGoIdent(protogen.GoIdent{
		GoImportPath: "context",
		GoName:       "Context",
	})
	f.P()
	f.P("func (m *", middleware.GoName(), ") ", g.Method.GoName, "(")
	f.P("ctx ", contextContext, ",")
	f.P("request *", g.Method.Input.GoIdent, ",")
	f.P(") ( ")
	f.P("*", g.Method.Output.GoIdent, ", error,")
	f.P(") {")
	switch g.Policy.GetDecisionPoint() {
	case iamv1.PolicyDecisionPoint_BEFORE:
		f.P("caller, err := m.callerFn(ctx)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("val, _, err := m.program", g.Method.GoName, ".Eval(map[string]interface{}{")
		f.P(`"caller": caller,`)
		f.P(`"request": request,`)
		f.P("})")
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
	case iamv1.PolicyDecisionPoint_AFTER:
		f.P("caller, err := m.callerFn(ctx)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("response, err := m.next.", g.Method.GoName, "(ctx, request)")
		f.P("if err != nil {")
		f.P("return nil, err")
		f.P("}")
		f.P("val, _, err := m.program", g.Method.GoName, ".Eval(map[string]interface{}{")
		f.P(`"caller": caller,`)
		f.P(`"request": request,`)
		f.P(`"response": response,`)
		f.P("})")
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
