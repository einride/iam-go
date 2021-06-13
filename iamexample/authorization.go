package iamexample

import (
	"context"

	"go.einride.tech/iam/iamauthz"
	"go.einride.tech/iam/iamspanner"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authorization struct {
	*iamexamplev1.FreightServiceAuthorization
	Next          iamexamplev1.FreightServiceServer
	IAMServer     *iamspanner.IAMServer
	IAMDescriptor *iamexamplev1.FreightServiceIAMDescriptor
}

var _ iamexamplev1.FreightServiceServer = &Authorization{}

func (a *Authorization) UpdateShipment(
	ctx context.Context,
	request *iamexamplev1.UpdateShipmentRequest,
) (*iamexamplev1.Shipment, error) {
	iamauthz.Authorize(ctx)
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

func (a *Authorization) BatchGetShipments(
	ctx context.Context,
	request *iamexamplev1.BatchGetShipmentsRequest,
) (*iamexamplev1.BatchGetShipmentsResponse, error) {
	iamauthz.Authorize(ctx)
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
	results, err := a.IAMServer.TestPermissionOnResources(ctx, permission, resources)
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

func (a *Authorization) test(ctx context.Context, permission, resource string) (bool, error) {
	return a.IAMServer.TestPermissionOnResource(ctx, permission, resource)
}

func (a *Authorization) testAny(ctx context.Context, permission string, resources []string) (bool, error) {
	results, err := a.IAMServer.TestPermissionOnResources(ctx, permission, resources)
	if err != nil {
		return false, err
	}
	for _, resource := range resources {
		if results[resource] {
			return true, nil
		}
	}
	return false, nil
}

func (a *Authorization) requireAny(ctx context.Context, permission string, resources []string) error {
	if ok, err := a.testAny(ctx, permission, resources); err != nil {
		return err
	} else if !ok {
		return status.Errorf(codes.PermissionDenied, "caller must have permission `%s` on any of the resources", permission)
	}
	return nil
}
