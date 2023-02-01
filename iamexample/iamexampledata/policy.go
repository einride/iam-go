package iamexampledata

import (
	"cloud.google.com/go/iam/apiv1/iampb"
)

// Example IAM policy members.
const (
	RootAdminMember                       = "user:root"
	EinrideAdminMember                    = "user:einride-admin"
	EinrideGothenburgFreightPlannerMember = "user:einride-gothenburg-freight-planner"
	EinrideBatcaveWorkerMember            = "user:einride-batcave-worker"
)

// EinrideSetIamPolicyRequest returns an iampb.SetIamPolicyRequest for the
// Einride shipper resource.
func EinrideSetIamPolicyRequest() *iampb.SetIamPolicyRequest {
	return &iampb.SetIamPolicyRequest{
		Resource: Einride().Name,
		Policy: &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/freight.admin",
					Members: []string{EinrideAdminMember},
				},
			},
		},
	}
}

// EinrideGothenburgOfficeSetIamPolicyRequest returns an iampb.SetIamPolicyRequest for the
// Einride Gothenburg Office site resource.
func EinrideGothenburgOfficeSetIamPolicyRequest() *iampb.SetIamPolicyRequest {
	return &iampb.SetIamPolicyRequest{
		Resource: EinrideGothenburgOffice().Name,
		Policy: &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/freight.editor",
					Members: []string{EinrideGothenburgFreightPlannerMember},
				},
			},
		},
	}
}

// EinrideBatcaveSetIamPolicyRequest returns an iampb.SetIamPolicyRequest for the
// Einride Batcave site resource.
func EinrideBatcaveSetIamPolicyRequest() *iampb.SetIamPolicyRequest {
	return &iampb.SetIamPolicyRequest{
		Resource: EinrideBatcave().Name,
		Policy: &iampb.Policy{
			Bindings: []*iampb.Binding{
				{
					Role:    "roles/freight.viewer",
					Members: []string{EinrideBatcaveWorkerMember},
				},
			},
		},
	}
}
