package iamexample

import (
	"context"

	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorization struct {
	Next iamexamplev1.FreightServiceServer
}

var _ iamexamplev1.FreightServiceServer = &Authorization{}

func (a *Authorization) GetShipper(
	ctx context.Context,
	request *iamexamplev1.GetShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.get"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetShipper(ctx, request)
}

func (a *Authorization) ListShippers(
	ctx context.Context,
	request *iamexamplev1.ListShippersRequest,
) (*iamexamplev1.ListShippersResponse, error) {
	const permission = "freight.shippers.list"
	if err := a.test(ctx, permission, "*"); err != nil {
		return nil, err
	}
	return a.Next.ListShippers(ctx, request)
}

func (a *Authorization) CreateShipper(
	ctx context.Context,
	request *iamexamplev1.CreateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.create"
	if err := a.test(ctx, permission, "*"); err != nil {
		return nil, err
	}
	return a.Next.CreateShipper(ctx, request)
}

func (a *Authorization) UpdateShipper(
	ctx context.Context,
	request *iamexamplev1.UpdateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.update"
	if err := a.test(ctx, permission, request.GetShipper().GetName()); err != nil {
		return nil, err
	}
	return a.Next.UpdateShipper(ctx, request)
}

func (a *Authorization) DeleteShipper(
	ctx context.Context,
	request *iamexamplev1.DeleteShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.delete"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteShipper(ctx, request)
}

func (a *Authorization) GetSite(
	ctx context.Context,
	request *iamexamplev1.GetSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.get"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetSite(ctx, request)
}

func (a *Authorization) ListSites(
	ctx context.Context,
	request *iamexamplev1.ListSitesRequest,
) (*iamexamplev1.ListSitesResponse, error) {
	const permission = "freight.sites.list"
	if err := a.test(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.ListSites(ctx, request)
}

func (a *Authorization) CreateSite(
	ctx context.Context,
	request *iamexamplev1.CreateSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.create"
	if err := a.test(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.CreateSite(ctx, request)
}

func (a *Authorization) UpdateSite(
	ctx context.Context,
	request *iamexamplev1.UpdateSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.update"
	if err := a.test(ctx, permission, request.GetSite().GetName()); err != nil {
		return nil, err
	}
	return a.Next.UpdateSite(ctx, request)
}

func (a *Authorization) DeleteSite(
	ctx context.Context,
	request *iamexamplev1.DeleteSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.delete"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteSite(ctx, request)
}

func (a *Authorization) BatchGetSites(
	ctx context.Context,
	request *iamexamplev1.BatchGetSitesRequest,
) (*iamexamplev1.BatchGetSitesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}

func (a *Authorization) SearchSites(
	ctx context.Context,
	request *iamexamplev1.SearchSitesRequest,
) (*iamexamplev1.SearchSitesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}

func (a *Authorization) GetShipment(
	ctx context.Context,
	request *iamexamplev1.GetShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.get"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetShipment(ctx, request)
}

func (a *Authorization) ListShipments(
	ctx context.Context,
	request *iamexamplev1.ListShipmentsRequest,
) (*iamexamplev1.ListShipmentsResponse, error) {
	const permission = "freight.shipments.list"
	if err := a.test(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.ListShipments(ctx, request)
}

func (a *Authorization) CreateShipment(
	ctx context.Context,
	request *iamexamplev1.CreateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.create"
	if err := a.test(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.CreateShipment(ctx, request)
}

func (a *Authorization) UpdateShipment(
	ctx context.Context,
	request *iamexamplev1.UpdateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.update"
	if err := a.test(ctx, permission, request.GetShipment().GetName()); err != nil {
		return nil, err
	}
	return a.Next.UpdateShipment(ctx, request)
}

func (a *Authorization) DeleteShipment(
	ctx context.Context,
	request *iamexamplev1.DeleteShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.delete"
	if err := a.test(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteShipment(ctx, request)
}

func (a *Authorization) BatchGetShipments(
	ctx context.Context,
	request *iamexamplev1.BatchGetShipmentsRequest,
) (*iamexamplev1.BatchGetShipmentsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}

func (a *Authorization) SetIamPolicy(
	ctx context.Context,
	request *iam.SetIamPolicyRequest,
) (*iam.Policy, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}

func (a *Authorization) GetIamPolicy(
	ctx context.Context,
	request *iam.GetIamPolicyRequest,
) (*iam.Policy, error) {
	return nil, status.Error(codes.Unimplemented, "TODO: implement me")
}

func (a *Authorization) TestIamPermissions(
	ctx context.Context,
	request *iam.TestIamPermissionsRequest,
) (*iam.TestIamPermissionsResponse, error) {
	return a.Next.TestIamPermissions(ctx, request)
}

func (a *Authorization) test(ctx context.Context, permission, resource string) error {
	response, err := a.Next.TestIamPermissions(ctx, &iam.TestIamPermissionsRequest{
		Resource:    resource,
		Permissions: []string{permission},
	})
	if err != nil {
		return err
	}
	if len(response.Permissions) != 1 || response.Permissions[0] != permission {
		return status.Errorf(codes.PermissionDenied, "caller must have permission `%s`", permission)
	}
	return nil
}
