package iamexampledata

import (
	"context"

	"cloud.google.com/go/spanner"
	"go.einride.tech/aip/resourcename"
	"go.einride.tech/iam/iamexample/iamexampledb"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/iam/v1"
)

// BootstrapRootAdmin bootstraps an IAM database with RootAdminMember as iamexamplev1.FreightServiceServer root admin.
func BootstrapRootAdmin(ctx context.Context, spannerClient *spanner.Client) error {
	if _, err := spannerClient.Apply(ctx, []*spanner.Mutation{
		spanner.Insert((&iamexampledb.IamPolicyBindingsRow{
			Resource:     "*",
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
func InitializeResources(ctx context.Context, client iamexamplev1.FreightServiceClient) error {
	einride := Einride()
	var shipperID string
	if err := resourcename.Sscan(einride.Name, "shippers/{shipper}", &shipperID); err != nil {
		return err
	}
	if _, err := client.CreateShipper(ctx, &iamexamplev1.CreateShipperRequest{
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
			site.Name,
			"shippers/{shipper}/sites/{site}",
			&shipperID,
			&siteID,
		); err != nil {
			return err
		}
		if _, err := client.CreateSite(ctx, &iamexamplev1.CreateSiteRequest{
			Parent: einride.Name,
			Site:   site,
			SiteId: siteID,
		}); err != nil {
			return err
		}
	}
	for _, request := range []*iam.SetIamPolicyRequest{
		EinrideSetIamPolicyRequest(),
		EinrideGothenburgOfficeSetIamPolicyRequest(),
		EinrideBatcaveSetIamPolicyRequest(),
	} {
		if _, err := client.SetIamPolicy(ctx, request); err != nil {
			return err
		}
	}
	return nil
}
