package iamexamplev1

import (
	v1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	proto "google.golang.org/protobuf/proto"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
)

// FreightServiceIAM returns a descriptor for the FreightService IAM configuration.
func FreightServiceIAM() FreightServiceIAMDescriptor {
	return &file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService
}

// FreightServiceIAMDescriptor describes the FreightService IAM configuration.
type FreightServiceIAMDescriptor interface {
	PredefinedRoles() *v1.Roles
	GetShipper() *v1.MethodAuthorizationOptions
	ListShippers() *v1.MethodAuthorizationOptions
	CreateShipper() *v1.MethodAuthorizationOptions
	UpdateShipper() *v1.MethodAuthorizationOptions
	DeleteShipper() *v1.MethodAuthorizationOptions
	GetSite() *v1.MethodAuthorizationOptions
	ListSites() *v1.MethodAuthorizationOptions
	CreateSite() *v1.MethodAuthorizationOptions
	UpdateSite() *v1.MethodAuthorizationOptions
	DeleteSite() *v1.MethodAuthorizationOptions
	BatchGetSites() *v1.MethodAuthorizationOptions
	SearchSites() *v1.MethodAuthorizationOptions
	GetShipment() *v1.MethodAuthorizationOptions
	ListShipments() *v1.MethodAuthorizationOptions
	CreateShipment() *v1.MethodAuthorizationOptions
	UpdateShipment() *v1.MethodAuthorizationOptions
	DeleteShipment() *v1.MethodAuthorizationOptions
	BatchGetShipments() *v1.MethodAuthorizationOptions
	SetIamPolicy() *v1.MethodAuthorizationOptions
	GetIamPolicy() *v1.MethodAuthorizationOptions
	TestIamPermissions() *v1.MethodAuthorizationOptions
}

type freightServiceIAMDescriptor struct {
	predefinedRoles    *v1.Roles
	getShipper         *v1.MethodAuthorizationOptions
	listShippers       *v1.MethodAuthorizationOptions
	createShipper      *v1.MethodAuthorizationOptions
	updateShipper      *v1.MethodAuthorizationOptions
	deleteShipper      *v1.MethodAuthorizationOptions
	getSite            *v1.MethodAuthorizationOptions
	listSites          *v1.MethodAuthorizationOptions
	createSite         *v1.MethodAuthorizationOptions
	updateSite         *v1.MethodAuthorizationOptions
	deleteSite         *v1.MethodAuthorizationOptions
	batchGetSites      *v1.MethodAuthorizationOptions
	searchSites        *v1.MethodAuthorizationOptions
	getShipment        *v1.MethodAuthorizationOptions
	listShipments      *v1.MethodAuthorizationOptions
	createShipment     *v1.MethodAuthorizationOptions
	updateShipment     *v1.MethodAuthorizationOptions
	deleteShipment     *v1.MethodAuthorizationOptions
	batchGetShipments  *v1.MethodAuthorizationOptions
	setIamPolicy       *v1.MethodAuthorizationOptions
	getIamPolicy       *v1.MethodAuthorizationOptions
	testIamPermissions *v1.MethodAuthorizationOptions
}

func (d *freightServiceIAMDescriptor) PredefinedRoles() *v1.Roles {
	return d.predefinedRoles
}

func (d *freightServiceIAMDescriptor) GetShipper() *v1.MethodAuthorizationOptions {
	return d.getShipper
}

func (d *freightServiceIAMDescriptor) ListShippers() *v1.MethodAuthorizationOptions {
	return d.listShippers
}

func (d *freightServiceIAMDescriptor) CreateShipper() *v1.MethodAuthorizationOptions {
	return d.createShipper
}

func (d *freightServiceIAMDescriptor) UpdateShipper() *v1.MethodAuthorizationOptions {
	return d.updateShipper
}

func (d *freightServiceIAMDescriptor) DeleteShipper() *v1.MethodAuthorizationOptions {
	return d.deleteShipper
}

func (d *freightServiceIAMDescriptor) GetSite() *v1.MethodAuthorizationOptions {
	return d.getSite
}

func (d *freightServiceIAMDescriptor) ListSites() *v1.MethodAuthorizationOptions {
	return d.listSites
}

func (d *freightServiceIAMDescriptor) CreateSite() *v1.MethodAuthorizationOptions {
	return d.createSite
}

func (d *freightServiceIAMDescriptor) UpdateSite() *v1.MethodAuthorizationOptions {
	return d.updateSite
}

func (d *freightServiceIAMDescriptor) DeleteSite() *v1.MethodAuthorizationOptions {
	return d.deleteSite
}

func (d *freightServiceIAMDescriptor) BatchGetSites() *v1.MethodAuthorizationOptions {
	return d.batchGetSites
}

func (d *freightServiceIAMDescriptor) SearchSites() *v1.MethodAuthorizationOptions {
	return d.searchSites
}

func (d *freightServiceIAMDescriptor) GetShipment() *v1.MethodAuthorizationOptions {
	return d.getShipment
}

func (d *freightServiceIAMDescriptor) ListShipments() *v1.MethodAuthorizationOptions {
	return d.listShipments
}

func (d *freightServiceIAMDescriptor) CreateShipment() *v1.MethodAuthorizationOptions {
	return d.createShipment
}

func (d *freightServiceIAMDescriptor) UpdateShipment() *v1.MethodAuthorizationOptions {
	return d.updateShipment
}

func (d *freightServiceIAMDescriptor) DeleteShipment() *v1.MethodAuthorizationOptions {
	return d.deleteShipment
}

func (d *freightServiceIAMDescriptor) BatchGetShipments() *v1.MethodAuthorizationOptions {
	return d.batchGetShipments
}

func (d *freightServiceIAMDescriptor) SetIamPolicy() *v1.MethodAuthorizationOptions {
	return d.setIamPolicy
}

func (d *freightServiceIAMDescriptor) GetIamPolicy() *v1.MethodAuthorizationOptions {
	return d.getIamPolicy
}

func (d *freightServiceIAMDescriptor) TestIamPermissions() *v1.MethodAuthorizationOptions {
	return d.testIamPermissions
}

var file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService freightServiceIAMDescriptor

func file_einride_iam_example_v1_freight_service_proto_init_iamDescriptor_FreightService() {
	// Init predefined roles.
	serviceDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService",
	)
	if err != nil {
		panic("unable to find service descriptor einride.iam.example.v1.FreightService")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.predefinedRoles = proto.GetExtension(
		serviceDesc.Options(),
		v1.E_PredefinedRoles,
	).(*v1.Roles)
	// Init GetShipper authorization options.
	methodDescGetShipper, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.GetShipper",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.GetShipper")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.getShipper = proto.GetExtension(
		methodDescGetShipper.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init ListShippers authorization options.
	methodDescListShippers, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.ListShippers",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.ListShippers")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.listShippers = proto.GetExtension(
		methodDescListShippers.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init CreateShipper authorization options.
	methodDescCreateShipper, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.CreateShipper",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.CreateShipper")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.createShipper = proto.GetExtension(
		methodDescCreateShipper.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init UpdateShipper authorization options.
	methodDescUpdateShipper, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.UpdateShipper",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.UpdateShipper")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.updateShipper = proto.GetExtension(
		methodDescUpdateShipper.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init DeleteShipper authorization options.
	methodDescDeleteShipper, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.DeleteShipper",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.DeleteShipper")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.deleteShipper = proto.GetExtension(
		methodDescDeleteShipper.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init GetSite authorization options.
	methodDescGetSite, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.GetSite",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.GetSite")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.getSite = proto.GetExtension(
		methodDescGetSite.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init ListSites authorization options.
	methodDescListSites, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.ListSites",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.ListSites")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.listSites = proto.GetExtension(
		methodDescListSites.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init CreateSite authorization options.
	methodDescCreateSite, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.CreateSite",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.CreateSite")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.createSite = proto.GetExtension(
		methodDescCreateSite.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init UpdateSite authorization options.
	methodDescUpdateSite, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.UpdateSite",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.UpdateSite")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.updateSite = proto.GetExtension(
		methodDescUpdateSite.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init DeleteSite authorization options.
	methodDescDeleteSite, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.DeleteSite",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.DeleteSite")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.deleteSite = proto.GetExtension(
		methodDescDeleteSite.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init BatchGetSites authorization options.
	methodDescBatchGetSites, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.BatchGetSites",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.BatchGetSites")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.batchGetSites = proto.GetExtension(
		methodDescBatchGetSites.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init SearchSites authorization options.
	methodDescSearchSites, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.SearchSites",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.SearchSites")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.searchSites = proto.GetExtension(
		methodDescSearchSites.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init GetShipment authorization options.
	methodDescGetShipment, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.GetShipment",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.GetShipment")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.getShipment = proto.GetExtension(
		methodDescGetShipment.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init ListShipments authorization options.
	methodDescListShipments, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.ListShipments",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.ListShipments")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.listShipments = proto.GetExtension(
		methodDescListShipments.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init CreateShipment authorization options.
	methodDescCreateShipment, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.CreateShipment",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.CreateShipment")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.createShipment = proto.GetExtension(
		methodDescCreateShipment.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init UpdateShipment authorization options.
	methodDescUpdateShipment, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.UpdateShipment",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.UpdateShipment")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.updateShipment = proto.GetExtension(
		methodDescUpdateShipment.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init DeleteShipment authorization options.
	methodDescDeleteShipment, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.DeleteShipment",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.DeleteShipment")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.deleteShipment = proto.GetExtension(
		methodDescDeleteShipment.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init BatchGetShipments authorization options.
	methodDescBatchGetShipments, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.BatchGetShipments",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.BatchGetShipments")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.batchGetShipments = proto.GetExtension(
		methodDescBatchGetShipments.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init SetIamPolicy authorization options.
	methodDescSetIamPolicy, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.SetIamPolicy",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.SetIamPolicy")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.setIamPolicy = proto.GetExtension(
		methodDescSetIamPolicy.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init GetIamPolicy authorization options.
	methodDescGetIamPolicy, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.GetIamPolicy",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.GetIamPolicy")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.getIamPolicy = proto.GetExtension(
		methodDescGetIamPolicy.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
	// Init TestIamPermissions authorization options.
	methodDescTestIamPermissions, err := protoregistry.GlobalFiles.FindDescriptorByName(
		"einride.iam.example.v1.FreightService.TestIamPermissions",
	)
	if err != nil {
		panic("unable to find method descriptor einride.iam.example.v1.FreightService.TestIamPermissions")
	}
	file_einride_iam_example_v1_freight_service_proto_iamDescriptor_FreightService.testIamPermissions = proto.GetExtension(
		methodDescTestIamPermissions.Options(),
		v1.E_MethodAuthorization,
	).(*v1.MethodAuthorizationOptions)
}

func init() {
	// This init function runs after the init function of go.einride.tech/iam/proto/gen/einride/iam/example/v1/freight_service.pb.go.
	// We depend on the Go initialization order to ensure this.
	// See: https://golang.org/ref/spec#Package_initialization
	file_einride_iam_example_v1_freight_service_proto_init_iamDescriptor_FreightService()
}
