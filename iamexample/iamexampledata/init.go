package iamexampledata

import (
	"context"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	"go.einride.tech/iam/iamresource"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
)

// BootstrapRootAdmin bootstraps an IAM database with RootAdminMember as iamexamplev1.FreightServiceServer root admin.
func BootstrapRootAdmin(ctx context.Context, spannerClient *spanner.Client) error {
	if _, err := spannerClient.Apply(ctx, []*spanner.Mutation{
		spanner.Insert((&iamexampledb.IamPolicyBindingsRow{
			Resource:     iamresource.Root,
			BindingIndex: 0,
			Role:         "roles/freight.admin",
			MemberIndex:  0,
			Member:       RootAdminMember,
		}).Mutate()),
	}); err != nil {
		return err
	}
	return nil
}

// InitializeResources uses an iamexamplev1.FreightServiceClient to initialize the set of example resources.
func InitializeResources(ctx context.Context, server iamexamplev1.FreightServiceServer) error {
	einride := Einride()
	var shipperID string
	if err := resourcename.Sscan(einride.GetName(), "shippers/{shipper}", &shipperID); err != nil {
		return err
	}
	if _, err := server.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
		Shipper:   einride,
		ShipperId: shipperID,
	}); err != nil {
		return err
	}
	for _, site := range []*iamexamplev1.Site{
		EinrideGothenburgOffice(),
		EinrideStockholmOffice(),
		EinrideBatcave(),
	} {
		var siteID string
		if err := resourcename.Sscan(
			site.GetName(),
			"shippers/{shipper}/sites/{site}",
			&shipperID,
			&siteID,
		); err != nil {
			return err
		}
		if _, err := server.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
			Parent: einride.GetName(),
			Site:   site,
			SiteId: siteID,
		}); err != nil {
			return err
		}
	}
	for _, request := range []*iampb.SetIamPolicyRequest{
		EinrideSetIamPolicyRequest(),
		EinrideGothenburgOfficeSetIamPolicyRequest(),
		EinrideBatcaveSetIamPolicyRequest(),
	} {
		if _, err := server.SetIamPolicy(ctx, request); err != nil {
			return err
		}
	}
	return nil
}
