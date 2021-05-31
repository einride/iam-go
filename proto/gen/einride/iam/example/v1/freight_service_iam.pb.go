// Code generated by protoc-gen-go-iam. DO NOT EDIT.
// versions:
// 	protoc            v3.15.2

package iamexamplev1

import (
	context "context"
	fmt "fmt"
	iamauthz "go.einride.tech/iam/iamauthz"
	iammember "go.einride.tech/iam/iammember"
	iamreflect "go.einride.tech/iam/iamreflect"
	v11 "google.golang.org/genproto/googleapis/iam/admin/v1"
	v1 "google.golang.org/genproto/googleapis/iam/v1"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// NewFreightServiceIAMDescriptor returns a new FreightService IAM descriptor.
func NewFreightServiceIAMDescriptor() (*iamreflect.IAMDescriptor, error) {
	descriptor, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService")
	if err != nil {
		return nil, fmt.Errorf("new FreightService IAM descriptor: %w", err)
	}
	service, ok := descriptor.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightServiceIAM descriptor: got non-service descriptor")
	}
	return iamreflect.NewIAMDescriptor(service, protoregistry.GlobalFiles)
}

// NewFreightServiceAuthorization creates a new authorization middleware for FreightService.
func NewFreightServiceAuthorization(
	next FreightServiceServer,
	permissionTester iamauthz.PermissionTester,
	memberResolver iammember.Resolver,
) (*FreightServiceAuthorization, error) {
	var result FreightServiceAuthorization
	result.next = next
	descriptorGetShipper, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.GetShipper")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for GetShipper")
	}
	methodGetShipper, ok := descriptorGetShipper.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for GetShipper")
	}
	beforeGetShipper, err := iamauthz.NewBeforeMethodAuthorization(methodGetShipper, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeGetShipper = beforeGetShipper
	descriptorListShippers, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.ListShippers")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for ListShippers")
	}
	methodListShippers, ok := descriptorListShippers.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for ListShippers")
	}
	beforeListShippers, err := iamauthz.NewBeforeMethodAuthorization(methodListShippers, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeListShippers = beforeListShippers
	descriptorCreateShipper, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.CreateShipper")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for CreateShipper")
	}
	methodCreateShipper, ok := descriptorCreateShipper.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for CreateShipper")
	}
	beforeCreateShipper, err := iamauthz.NewBeforeMethodAuthorization(methodCreateShipper, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeCreateShipper = beforeCreateShipper
	descriptorUpdateShipper, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.UpdateShipper")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for UpdateShipper")
	}
	methodUpdateShipper, ok := descriptorUpdateShipper.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for UpdateShipper")
	}
	beforeUpdateShipper, err := iamauthz.NewBeforeMethodAuthorization(methodUpdateShipper, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeUpdateShipper = beforeUpdateShipper
	descriptorDeleteShipper, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.DeleteShipper")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for DeleteShipper")
	}
	methodDeleteShipper, ok := descriptorDeleteShipper.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for DeleteShipper")
	}
	beforeDeleteShipper, err := iamauthz.NewBeforeMethodAuthorization(methodDeleteShipper, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeDeleteShipper = beforeDeleteShipper
	descriptorGetSite, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.GetSite")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for GetSite")
	}
	methodGetSite, ok := descriptorGetSite.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for GetSite")
	}
	beforeGetSite, err := iamauthz.NewBeforeMethodAuthorization(methodGetSite, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeGetSite = beforeGetSite
	descriptorListSites, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.ListSites")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for ListSites")
	}
	methodListSites, ok := descriptorListSites.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for ListSites")
	}
	beforeListSites, err := iamauthz.NewBeforeMethodAuthorization(methodListSites, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeListSites = beforeListSites
	descriptorCreateSite, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.CreateSite")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for CreateSite")
	}
	methodCreateSite, ok := descriptorCreateSite.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for CreateSite")
	}
	beforeCreateSite, err := iamauthz.NewBeforeMethodAuthorization(methodCreateSite, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeCreateSite = beforeCreateSite
	descriptorUpdateSite, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.UpdateSite")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for UpdateSite")
	}
	methodUpdateSite, ok := descriptorUpdateSite.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for UpdateSite")
	}
	beforeUpdateSite, err := iamauthz.NewBeforeMethodAuthorization(methodUpdateSite, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeUpdateSite = beforeUpdateSite
	descriptorDeleteSite, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.DeleteSite")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for DeleteSite")
	}
	methodDeleteSite, ok := descriptorDeleteSite.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for DeleteSite")
	}
	beforeDeleteSite, err := iamauthz.NewBeforeMethodAuthorization(methodDeleteSite, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeDeleteSite = beforeDeleteSite
	descriptorBatchGetSites, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.BatchGetSites")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for BatchGetSites")
	}
	methodBatchGetSites, ok := descriptorBatchGetSites.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for BatchGetSites")
	}
	beforeBatchGetSites, err := iamauthz.NewBeforeMethodAuthorization(methodBatchGetSites, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeBatchGetSites = beforeBatchGetSites
	descriptorSearchSites, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.SearchSites")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for SearchSites")
	}
	methodSearchSites, ok := descriptorSearchSites.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for SearchSites")
	}
	afterSearchSites, err := iamauthz.NewAfterMethodAuthorization(methodSearchSites, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.afterSearchSites = afterSearchSites
	descriptorGetShipment, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.GetShipment")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for GetShipment")
	}
	methodGetShipment, ok := descriptorGetShipment.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for GetShipment")
	}
	afterGetShipment, err := iamauthz.NewAfterMethodAuthorization(methodGetShipment, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.afterGetShipment = afterGetShipment
	descriptorListShipments, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.ListShipments")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for ListShipments")
	}
	methodListShipments, ok := descriptorListShipments.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for ListShipments")
	}
	beforeListShipments, err := iamauthz.NewBeforeMethodAuthorization(methodListShipments, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeListShipments = beforeListShipments
	descriptorCreateShipment, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.CreateShipment")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for CreateShipment")
	}
	methodCreateShipment, ok := descriptorCreateShipment.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for CreateShipment")
	}
	beforeCreateShipment, err := iamauthz.NewBeforeMethodAuthorization(methodCreateShipment, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeCreateShipment = beforeCreateShipment
	descriptorDeleteShipment, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.DeleteShipment")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for DeleteShipment")
	}
	methodDeleteShipment, ok := descriptorDeleteShipment.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for DeleteShipment")
	}
	beforeDeleteShipment, err := iamauthz.NewBeforeMethodAuthorization(methodDeleteShipment, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeDeleteShipment = beforeDeleteShipment
	descriptorBatchGetShipments, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.BatchGetShipments")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for BatchGetShipments")
	}
	methodBatchGetShipments, ok := descriptorBatchGetShipments.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for BatchGetShipments")
	}
	afterBatchGetShipments, err := iamauthz.NewAfterMethodAuthorization(methodBatchGetShipments, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.afterBatchGetShipments = afterBatchGetShipments
	descriptorSetIamPolicy, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.SetIamPolicy")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for SetIamPolicy")
	}
	methodSetIamPolicy, ok := descriptorSetIamPolicy.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for SetIamPolicy")
	}
	beforeSetIamPolicy, err := iamauthz.NewBeforeMethodAuthorization(methodSetIamPolicy, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeSetIamPolicy = beforeSetIamPolicy
	descriptorGetIamPolicy, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.GetIamPolicy")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for GetIamPolicy")
	}
	methodGetIamPolicy, ok := descriptorGetIamPolicy.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for GetIamPolicy")
	}
	beforeGetIamPolicy, err := iamauthz.NewBeforeMethodAuthorization(methodGetIamPolicy, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeGetIamPolicy = beforeGetIamPolicy
	descriptorListRoles, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.ListRoles")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for ListRoles")
	}
	methodListRoles, ok := descriptorListRoles.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for ListRoles")
	}
	beforeListRoles, err := iamauthz.NewBeforeMethodAuthorization(methodListRoles, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeListRoles = beforeListRoles
	descriptorGetRole, err := protoregistry.GlobalFiles.FindDescriptorByName("einride.iam.example.v1.FreightService.GetRole")
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: failed to find descriptor for GetRole")
	}
	methodGetRole, ok := descriptorGetRole.(protoreflect.MethodDescriptor)
	if !ok {
		return nil, fmt.Errorf("new FreightService authorization: got non-method descriptor for GetRole")
	}
	beforeGetRole, err := iamauthz.NewBeforeMethodAuthorization(methodGetRole, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeGetRole = beforeGetRole
	iamDescriptor, err := NewFreightServiceIAMDescriptor()
	if err != nil {
		return nil, err
	}
	beforeLongRunningOperationMethod, err := iamauthz.NewBeforeLongRunningOperationMethodAuthorization(iamDescriptor.LongRunningOperationsAuthorization.OperationPermissions, permissionTester, memberResolver)
	if err != nil {
		return nil, fmt.Errorf("new FreightService authorization: %w", err)
	}
	result.beforeLongRunningOperationMethod = beforeLongRunningOperationMethod
	return &result, nil
}

type FreightServiceAuthorization struct {
	next                             FreightServiceServer
	beforeGetShipper                 *iamauthz.BeforeMethodAuthorization
	beforeListShippers               *iamauthz.BeforeMethodAuthorization
	beforeCreateShipper              *iamauthz.BeforeMethodAuthorization
	beforeUpdateShipper              *iamauthz.BeforeMethodAuthorization
	beforeDeleteShipper              *iamauthz.BeforeMethodAuthorization
	beforeGetSite                    *iamauthz.BeforeMethodAuthorization
	beforeListSites                  *iamauthz.BeforeMethodAuthorization
	beforeCreateSite                 *iamauthz.BeforeMethodAuthorization
	beforeUpdateSite                 *iamauthz.BeforeMethodAuthorization
	beforeDeleteSite                 *iamauthz.BeforeMethodAuthorization
	beforeBatchGetSites              *iamauthz.BeforeMethodAuthorization
	afterSearchSites                 *iamauthz.AfterMethodAuthorization
	afterGetShipment                 *iamauthz.AfterMethodAuthorization
	beforeListShipments              *iamauthz.BeforeMethodAuthorization
	beforeCreateShipment             *iamauthz.BeforeMethodAuthorization
	beforeDeleteShipment             *iamauthz.BeforeMethodAuthorization
	afterBatchGetShipments           *iamauthz.AfterMethodAuthorization
	beforeSetIamPolicy               *iamauthz.BeforeMethodAuthorization
	beforeGetIamPolicy               *iamauthz.BeforeMethodAuthorization
	beforeListRoles                  *iamauthz.BeforeMethodAuthorization
	beforeGetRole                    *iamauthz.BeforeMethodAuthorization
	beforeLongRunningOperationMethod *iamauthz.BeforeLongRunningOperationMethodAuthorization
}

func (a *FreightServiceAuthorization) GetShipper(
	ctx context.Context,
	request *GetShipperRequest,
) (*Shipper, error) {
	ctx, err := a.beforeGetShipper.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.GetShipper(ctx, request)
}

func (a *FreightServiceAuthorization) ListShippers(
	ctx context.Context,
	request *ListShippersRequest,
) (*ListShippersResponse, error) {
	ctx, err := a.beforeListShippers.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.ListShippers(ctx, request)
}

func (a *FreightServiceAuthorization) CreateShipper(
	ctx context.Context,
	request *CreateShipperRequest,
) (*Shipper, error) {
	ctx, err := a.beforeCreateShipper.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.CreateShipper(ctx, request)
}

func (a *FreightServiceAuthorization) UpdateShipper(
	ctx context.Context,
	request *UpdateShipperRequest,
) (*Shipper, error) {
	ctx, err := a.beforeUpdateShipper.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.UpdateShipper(ctx, request)
}

func (a *FreightServiceAuthorization) DeleteShipper(
	ctx context.Context,
	request *DeleteShipperRequest,
) (*longrunning.Operation, error) {
	ctx, err := a.beforeDeleteShipper.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.DeleteShipper(ctx, request)
}

func (a *FreightServiceAuthorization) GetSite(
	ctx context.Context,
	request *GetSiteRequest,
) (*Site, error) {
	ctx, err := a.beforeGetSite.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.GetSite(ctx, request)
}

func (a *FreightServiceAuthorization) ListSites(
	ctx context.Context,
	request *ListSitesRequest,
) (*ListSitesResponse, error) {
	ctx, err := a.beforeListSites.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.ListSites(ctx, request)
}

func (a *FreightServiceAuthorization) CreateSite(
	ctx context.Context,
	request *CreateSiteRequest,
) (*Site, error) {
	ctx, err := a.beforeCreateSite.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.CreateSite(ctx, request)
}

func (a *FreightServiceAuthorization) UpdateSite(
	ctx context.Context,
	request *UpdateSiteRequest,
) (*Site, error) {
	ctx, err := a.beforeUpdateSite.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.UpdateSite(ctx, request)
}

func (a *FreightServiceAuthorization) DeleteSite(
	ctx context.Context,
	request *DeleteSiteRequest,
) (*Site, error) {
	ctx, err := a.beforeDeleteSite.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.DeleteSite(ctx, request)
}

func (a *FreightServiceAuthorization) BatchGetSites(
	ctx context.Context,
	request *BatchGetSitesRequest,
) (*BatchGetSitesResponse, error) {
	ctx, err := a.beforeBatchGetSites.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.BatchGetSites(ctx, request)
}

func (a *FreightServiceAuthorization) SearchSites(
	ctx context.Context,
	request *SearchSitesRequest,
) (*SearchSitesResponse, error) {
	response, err := a.next.SearchSites(ctx, request)
	_, errAuth := a.afterSearchSites.AuthorizeRequestAndResponse(ctx, request, response)
	if errAuth != nil {
		return nil, errAuth
	}
	return response, err
}

func (a *FreightServiceAuthorization) GetShipment(
	ctx context.Context,
	request *GetShipmentRequest,
) (*Shipment, error) {
	response, err := a.next.GetShipment(ctx, request)
	_, errAuth := a.afterGetShipment.AuthorizeRequestAndResponse(ctx, request, response)
	if errAuth != nil {
		return nil, errAuth
	}
	return response, err
}

func (a *FreightServiceAuthorization) ListShipments(
	ctx context.Context,
	request *ListShipmentsRequest,
) (*ListShipmentsResponse, error) {
	ctx, err := a.beforeListShipments.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.ListShipments(ctx, request)
}

func (a *FreightServiceAuthorization) CreateShipment(
	ctx context.Context,
	request *CreateShipmentRequest,
) (*Shipment, error) {
	ctx, err := a.beforeCreateShipment.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.CreateShipment(ctx, request)
}

func (a *FreightServiceAuthorization) UpdateShipment(
	ctx context.Context,
	request *UpdateShipmentRequest,
) (*Shipment, error) {
	return nil, status.Error(codes.Internal, "custom authorization not implemented")
}

func (a *FreightServiceAuthorization) DeleteShipment(
	ctx context.Context,
	request *DeleteShipmentRequest,
) (*Shipment, error) {
	ctx, err := a.beforeDeleteShipment.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.DeleteShipment(ctx, request)
}

func (a *FreightServiceAuthorization) BatchGetShipments(
	ctx context.Context,
	request *BatchGetShipmentsRequest,
) (*BatchGetShipmentsResponse, error) {
	response, err := a.next.BatchGetShipments(ctx, request)
	_, errAuth := a.afterBatchGetShipments.AuthorizeRequestAndResponse(ctx, request, response)
	if errAuth != nil {
		return nil, errAuth
	}
	return response, err
}

func (a *FreightServiceAuthorization) SetIamPolicy(
	ctx context.Context,
	request *v1.SetIamPolicyRequest,
) (*v1.Policy, error) {
	ctx, err := a.beforeSetIamPolicy.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.SetIamPolicy(ctx, request)
}

func (a *FreightServiceAuthorization) GetIamPolicy(
	ctx context.Context,
	request *v1.GetIamPolicyRequest,
) (*v1.Policy, error) {
	ctx, err := a.beforeGetIamPolicy.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.GetIamPolicy(ctx, request)
}

func (a *FreightServiceAuthorization) TestIamPermissions(
	ctx context.Context,
	request *v1.TestIamPermissionsRequest,
) (*v1.TestIamPermissionsResponse, error) {
	iamauthz.Authorize(ctx)
	return a.next.TestIamPermissions(ctx, request)
}

func (a *FreightServiceAuthorization) ListRoles(
	ctx context.Context,
	request *v11.ListRolesRequest,
) (*v11.ListRolesResponse, error) {
	ctx, err := a.beforeListRoles.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.ListRoles(ctx, request)
}

func (a *FreightServiceAuthorization) GetRole(
	ctx context.Context,
	request *v11.GetRoleRequest,
) (*v11.Role, error) {
	ctx, err := a.beforeGetRole.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return a.next.GetRole(ctx, request)
}

func (a *FreightServiceAuthorization) ListOperations(
	ctx context.Context,
	request *longrunning.ListOperationsRequest,
) (*longrunning.ListOperationsResponse, error) {
	ctx, err := a.beforeLongRunningOperationMethod.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	impl, ok := a.next.(interface {
		ListOperations(context.Context, *longrunning.ListOperationsRequest) (*longrunning.ListOperationsResponse, error)
	})
	if !ok {
		return nil, status.Error(codes.Unimplemented, "ListOperations not implemented")
	}
	return impl.ListOperations(ctx, request)
}

func (a *FreightServiceAuthorization) GetOperation(
	ctx context.Context,
	request *longrunning.GetOperationRequest,
) (*longrunning.Operation, error) {
	ctx, err := a.beforeLongRunningOperationMethod.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	impl, ok := a.next.(interface {
		GetOperation(context.Context, *longrunning.GetOperationRequest) (*longrunning.Operation, error)
	})
	if !ok {
		return nil, status.Error(codes.Unimplemented, "GetOperation not implemented")
	}
	return impl.GetOperation(ctx, request)
}

func (a *FreightServiceAuthorization) DeleteOperation(
	ctx context.Context,
	request *longrunning.DeleteOperationRequest,
) (*emptypb.Empty, error) {
	ctx, err := a.beforeLongRunningOperationMethod.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	impl, ok := a.next.(interface {
		DeleteOperation(context.Context, *longrunning.DeleteOperationRequest) (*emptypb.Empty, error)
	})
	if !ok {
		return nil, status.Error(codes.Unimplemented, "DeleteOperation not implemented")
	}
	return impl.DeleteOperation(ctx, request)
}

func (a *FreightServiceAuthorization) CancelOperation(
	ctx context.Context,
	request *longrunning.CancelOperationRequest,
) (*emptypb.Empty, error) {
	ctx, err := a.beforeLongRunningOperationMethod.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	impl, ok := a.next.(interface {
		CancelOperation(context.Context, *longrunning.CancelOperationRequest) (*emptypb.Empty, error)
	})
	if !ok {
		return nil, status.Error(codes.Unimplemented, "CancelOperation not implemented")
	}
	return impl.CancelOperation(ctx, request)
}

func (a *FreightServiceAuthorization) WaitOperation(
	ctx context.Context,
	request *longrunning.WaitOperationRequest,
) (*longrunning.Operation, error) {
	ctx, err := a.beforeLongRunningOperationMethod.AuthorizeRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	impl, ok := a.next.(interface {
		WaitOperation(context.Context, *longrunning.WaitOperationRequest) (*longrunning.Operation, error)
	})
	if !ok {
		return nil, status.Error(codes.Unimplemented, "WaitOperation not implemented")
	}
	return impl.WaitOperation(ctx, request)
}
