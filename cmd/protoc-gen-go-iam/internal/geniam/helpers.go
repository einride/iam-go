package geniam

import (
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

func getLongRunningOperationsAuthorization(service *protogen.Service) *iamv1.LongRunningOperationsAuthorization {
	return proto.GetExtension(
		service.Desc.Options(), iamv1.E_LongRunningOperationsAuthorization,
	).(*iamv1.LongRunningOperationsAuthorization)
}
