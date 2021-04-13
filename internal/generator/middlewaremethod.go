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
	f.P()
	f.P("func (m *", middleware.GoName(), ") ", g.Method.GoName, "(")
	f.P("ctx ", contextContext, ",")
	f.P("request *", g.Method.Input.GoIdent, ",")
	f.P(") ( ")
	f.P("*", g.Method.Output.GoIdent, ", error,")
	f.P(") {")
	if g.Policy == nil {
		f.P("return m.next.", g.Method.GoName, "(ctx, request)")
		f.P("}")
		return
	}
	f.P("// TODO: Run the following policy: ", g.Policy.Expression)
	f.P("return nil, nil")
	f.P("}")
}
