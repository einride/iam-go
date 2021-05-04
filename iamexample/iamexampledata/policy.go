package iamexampledata

import "google.golang.org/genproto/googleapis/iam/v1"

// Example IAM policy members.
const (
	RootAdminMember                       = "user:root"
	EinrideAdminMember                    = "user:einride-admin"
	EinrideGothenburgFreightPlannerMember = "user:einride-gothenburg-freight-planner"
	EinrideBatcaveWorkerMember            = "user:einride-batcave-worker"
)

// EinrideSetIamPolicyRequest returns an iam.SetIamPolicyRequest for the
// Einride shipper resource.
func EinrideSetIamPolicyRequest() *iam.SetIamPolicyRequest {
	return &iam.SetIamPolicyRequest{
		Resource: Einride().Name,
		Policy: &iam.Policy{
			Bindings: []*iam.Binding{
				{
					Role:    "roles/freight.admin",
					Members: []string{EinrideAdminMember},
				},
			},
		},
	}
}

// EinrideGothenburgOfficeSetIamPolicyRequest returns an iam.SetIamPolicyRequest for the
// Einride Gothenburg Office site resource.
func EinrideGothenburgOfficeSetIamPolicyRequest() *iam.SetIamPolicyRequest {
	return &iam.SetIamPolicyRequest{
		Resource: EinrideGothenburgOffice().Name,
		Policy: &iam.Policy{
			Bindings: []*iam.Binding{
				{
					Role:    "roles/freight.editor",
					Members: []string{EinrideGothenburgFreightPlannerMember},
				},
			},
		},
	}
}

// EinrideBatcaveSetIamPolicyRequest returns an iam.SetIamPolicyRequest for the
// Einride Batcave site resource.
func EinrideBatcaveSetIamPolicyRequest() *iam.SetIamPolicyRequest {
	return &iam.SetIamPolicyRequest{
		Resource: EinrideBatcave().Name,
		Policy: &iam.Policy{
			Bindings: []*iam.Binding{
				{
					Role:    "roles/freight.viewer",
					Members: []string{EinrideBatcaveWorkerMember},
				},
			},
		},
	}
}
