package iamexampledata

import (
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// PredefinedRoles returns the pre-defined roles for iamexamplev1.FreightServiceServer.
func PredefinedRoles() *iamv1.Roles {
	desc, err := protoregistry.GlobalFiles.FindDescriptorByName(
		protoreflect.FullName(iamexamplev1.FreightService_ServiceDesc.ServiceName),
	)
	if err != nil {
		return nil
	}
	serviceDesc, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil
	}
	return proto.GetExtension(serviceDesc.Options(), iamv1.E_PredefinedRoles).(*iamv1.Roles)
}
