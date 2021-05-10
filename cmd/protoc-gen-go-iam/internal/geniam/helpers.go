package geniam

import (
	"strings"
	"unicode/utf8"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func getPredefinedRoles(service *protogen.Service) *iamv1.Roles {
	return proto.GetExtension(service.Desc.Options(), iamv1.E_PredefinedRoles).(*iamv1.Roles)
}

func getMethodAuthorizationOptions(method *protogen.Method) *iamv1.MethodAuthorizationOptions {
	return proto.GetExtension(method.Desc.Options(), iamv1.E_MethodAuthorization).(*iamv1.MethodAuthorizationOptions)
}

func private(s string) string {
	_, n := utf8.DecodeRuneInString(s)
	return strings.ToLower(s[:n]) + s[n:]
}

func fileVarName(f *protogen.File, suffix string) string {
	return private(f.GoDescriptorIdent.GoName) + "_" + suffix
}
