package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
)

type Middleware struct {
	Service *protogen.Service
}

func (g Middleware) GoName() string {
	return g.Service.GoName + "AuthorizationMiddleware"
}

func (g Middleware) Generate(f *protogen.GeneratedFile) {
	f.P()
	f.P("type ", g.Service.GoName, "AuthorizationMiddleware struct {")
	f.P("next ", g.Service.GoName, "Server")
	f.P("}")
	f.P("")
	f.P("var _ ", g.Service.GoName, "Server = &", g.GoName(), "{}")
	f.P()
	f.P("func (*", g.GoName(), ") mustEmbedUnimplemented", g.Service.GoName, "Server() {}")
}
