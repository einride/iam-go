package geniam

import (
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func getPredefinedRoles(
	service *protogen.Service,
) *iamv1.PredefinedRoles {
	return proto.GetExtension(service.Desc.Options(), iamv1.E_PredefinedRoles).(*iamv1.PredefinedRoles)
}

func getMethodAuthorizationOptions(
	method *protogen.Method,
) *iamv1.MethodAuthorizationOptions {
	return proto.GetExtension(method.Desc.Options(), iamv1.E_MethodAuthorization).(*iamv1.MethodAuthorizationOptions)
}

func getLongRunningOperationsAuthorizationOptions(
	service *protogen.Service,
) *iamv1.LongRunningOperationsAuthorizationOptions {
	return proto.GetExtension(
		service.Desc.Options(), iamv1.E_LongRunningOperationsAuthorization,
	).(*iamv1.LongRunningOperationsAuthorizationOptions)
}
