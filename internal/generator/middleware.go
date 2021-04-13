package generator

import (
	authorizationv1 "go.einride.tech/protoc-gen-go-authorization-policy/proto/gen/einride/authorization/v1"
	"google.golang.org/protobuf/compiler/protogen"
)

type Middleware struct {
	Method *protogen.Method
	Policy *authorizationv1.Policy
}

func (g Middleware) Generate(f *protogen.GeneratedFile) {
	// TODO: Implement me.
}
