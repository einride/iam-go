package iamexample

import (
	"context"

	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorization struct {
	Next iamexamplev1.FreightServiceServer
	IAM  *iamspanner.Server
}

var _ iamexamplev1.FreightServiceServer = &Authorization{}

func (a *Authorization) GetShipper(
	ctx context.Context,
	request *iamexamplev1.GetShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.get"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetShipper(ctx, request)
}

func (a *Authorization) ListShippers(
	ctx context.Context,
	request *iamexamplev1.ListShippersRequest,
) (*iamexamplev1.ListShippersResponse, error) {
	const permission = "freight.shippers.list"
	if err := a.require(ctx, permission, "*"); err != nil {
		return nil, err
	}
	return a.Next.ListShippers(ctx, request)
}

func (a *Authorization) CreateShipper(
	ctx context.Context,
	request *iamexamplev1.CreateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.create"
	if err := a.require(ctx, permission, "*"); err != nil {
		return nil, err
	}
	return a.Next.CreateShipper(ctx, request)
}

func (a *Authorization) UpdateShipper(
	ctx context.Context,
	request *iamexamplev1.UpdateShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.update"
	if err := a.require(ctx, permission, request.GetShipper().GetName()); err != nil {
		return nil, err
	}
	return a.Next.UpdateShipper(ctx, request)
}

func (a *Authorization) DeleteShipper(
	ctx context.Context,
	request *iamexamplev1.DeleteShipperRequest,
) (*iamexamplev1.Shipper, error) {
	const permission = "freight.shippers.delete"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteShipper(ctx, request)
}

func (a *Authorization) GetSite(
	ctx context.Context,
	request *iamexamplev1.GetSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.get"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetSite(ctx, request)
}

func (a *Authorization) ListSites(
	ctx context.Context,
	request *iamexamplev1.ListSitesRequest,
) (*iamexamplev1.ListSitesResponse, error) {
	const permission = "freight.sites.list"
	if err := a.require(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.ListSites(ctx, request)
}

func (a *Authorization) CreateSite(
	ctx context.Context,
	request *iamexamplev1.CreateSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.create"
	if err := a.require(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.CreateSite(ctx, request)
}

func (a *Authorization) UpdateSite(
	ctx context.Context,
	request *iamexamplev1.UpdateSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.update"
	if err := a.require(ctx, permission, request.GetSite().GetName()); err != nil {
		return nil, err
	}
	return a.Next.UpdateSite(ctx, request)
}

func (a *Authorization) DeleteSite(
	ctx context.Context,
	request *iamexamplev1.DeleteSiteRequest,
) (*iamexamplev1.Site, error) {
	const permission = "freight.sites.delete"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteSite(ctx, request)
}

func (a *Authorization) BatchGetSites(
	ctx context.Context,
	request *iamexamplev1.BatchGetSitesRequest,
) (*iamexamplev1.BatchGetSitesResponse, error) {
	const permission = "freight.sites.get"
	if request.Parent != "" {
		if ok, err := a.test(ctx, permission, request.Parent); err != nil {
			return nil, err
		} else if ok {
			return a.Next.BatchGetSites(ctx, request)
		}
	}
	if err := a.requireAll(ctx, permission, request.Names); err != nil {
		return nil, err
	}
	return a.Next.BatchGetSites(ctx, request)
}

func (a *Authorization) SearchSites(
	ctx context.Context,
	request *iamexamplev1.SearchSitesRequest,
) (*iamexamplev1.SearchSitesResponse, error) {
	const permission = "freight.sites.get"
	if request.Parent != "" {
		if ok, err := a.test(ctx, permission, request.Parent); err != nil {
			return nil, err
		} else if ok {
			return a.Next.SearchSites(ctx, request)
		}
	}
	response, err := a.Next.SearchSites(ctx, request)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(response.Sites))
	for _, site := range response.Sites {
		names = append(names, site.Name)
	}
	if err := a.requireAll(ctx, permission, names); err != nil {
		return nil, err
	}
	return response, nil
}

func (a *Authorization) GetShipment(
	ctx context.Context,
	request *iamexamplev1.GetShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.get"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.GetShipment(ctx, request)
}

func (a *Authorization) ListShipments(
	ctx context.Context,
	request *iamexamplev1.ListShipmentsRequest,
) (*iamexamplev1.ListShipmentsResponse, error) {
	const permission = "freight.shipments.list"
	if err := a.require(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.ListShipments(ctx, request)
}

func (a *Authorization) CreateShipment(
	ctx context.Context,
	request *iamexamplev1.CreateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.create"
	if err := a.require(ctx, permission, request.GetParent()); err != nil {
		return nil, err
	}
	return a.Next.CreateShipment(ctx, request)
}

func (a *Authorization) UpdateShipment(
	ctx context.Context,
	request *iamexamplev1.UpdateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.update"
	ok, err := a.test(ctx, permission, request.GetShipment().GetName())
	if err != nil {
		return nil, err
	}
	if ok {
		return a.Next.UpdateShipment(ctx, request)
	}
	shipment, err := a.GetShipment(ctx, &iamexamplev1.GetShipmentRequest{
		Name: request.GetShipment().GetName(),
	})
	if err != nil {
		return nil, err
	}
	if err := a.requireAny(ctx, permission, []string{shipment.OriginSite, shipment.DestinationSite}); err != nil {
		return nil, err
	}
	return a.Next.UpdateShipment(ctx, request)
}

func (a *Authorization) DeleteShipment(
	ctx context.Context,
	request *iamexamplev1.DeleteShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	const permission = "freight.shipments.delete"
	if err := a.require(ctx, permission, request.GetName()); err != nil {
		return nil, err
	}
	return a.Next.DeleteShipment(ctx, request)
}

func (a *Authorization) BatchGetShipments(
	ctx context.Context,
	request *iamexamplev1.BatchGetShipmentsRequest,
) (*iamexamplev1.BatchGetShipmentsResponse, error) {
	const permission = "freight.shipments.get"
	if request.Parent != "" {
		if ok, err := a.test(ctx, permission, request.Parent); err != nil {
			return nil, err
		} else if ok {
			return a.Next.BatchGetShipments(ctx, request)
		}
	}
	response, err := a.Next.BatchGetShipments(ctx, request)
	if err != nil {
		return nil, err
	}
	resources := make([]string, 0, 3*len(response.Shipments))
	for _, shipment := range response.Shipments {
		resources = append(resources, shipment.Name, shipment.OriginSite, shipment.DestinationSite)
	}
	results, err := a.IAM.TestPermissionOnResources(ctx, permission, resources)
	if err != nil {
		return nil, err
	}
	for _, shipment := range response.Shipments {
		if !(results[shipment.Name] || results[shipment.OriginSite] || results[shipment.DestinationSite]) {
			return nil, status.Errorf(codes.PermissionDenied, "missing permission %s for %s", permission, shipment.Name)
		}
	}
	return response, nil
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

func (a *Authorization) require(ctx context.Context, permission, resource string) error {
	if ok, err := a.test(ctx, permission, resource); err != nil {
		return err
	} else if !ok {
		return status.Errorf(codes.PermissionDenied, "caller must have permission `%s`", permission)
	}
	return nil
}

func (a *Authorization) test(ctx context.Context, permission, resource string) (bool, error) {
	return a.IAM.TestPermissionOnResource(ctx, permission, resource)
}

func (a *Authorization) testAll(ctx context.Context, permission string, resources []string) (bool, error) {
	results, err := a.IAM.TestPermissionOnResources(ctx, permission, resources)
	if err != nil {
		return false, err
	}
	all := true
	for _, resource := range resources {
		all = all && results[resource]
	}
	return all, nil
}

func (a *Authorization) requireAll(ctx context.Context, permission string, resources []string) error {
	if ok, err := a.testAll(ctx, permission, resources); err != nil {
		return err
	} else if !ok {
		return status.Errorf(codes.PermissionDenied, "caller must have permission `%s` on all resources", permission)
	}
	return nil
}

func (a *Authorization) requireAny(ctx context.Context, permission string, resources []string) error {
	if ok, err := a.testAll(ctx, permission, resources); err != nil {
		return err
	} else if !ok {
		return status.Errorf(codes.PermissionDenied, "caller must have permission `%s` on any of the resources", permission)
	}
	return nil
}
