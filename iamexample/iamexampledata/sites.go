package iamexampledata

import (
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// EinrideGothenburgOffice returns an iamexamplev1.Site representing Einride's Gothenburg office.
func EinrideGothenburgOffice() *iamexamplev1.Site {
	return &iamexamplev1.Site{
		Name:        Einride().GetName() + "/sites/gothenburg",
		DisplayName: "Einride Gothenburg Office",
		LatLng: &latlng.LatLng{
			Latitude:  57.70775726491335,
			Longitude: 11.94977470756508,
		},
	}
}

// EinrideStockholmOffice returns an iamexamplev1.Site representing Einride's Stockholm office.
func EinrideStockholmOffice() *iamexamplev1.Site {
	return &iamexamplev1.Site{
		Name:        Einride().GetName() + "/sites/sthlm",
		DisplayName: "Einride Gothenburg Office",
		LatLng: &latlng.LatLng{
			Latitude:  59.33749110496606,
			Longitude: 18.063672598779984,
		},
	}
}

// EinrideBatcave returns an iamexamplev1.Site representing Einride's Batcave.
func EinrideBatcave() *iamexamplev1.Site {
	return &iamexamplev1.Site{
		Name:        Einride().GetName() + "/sites/batcave",
		DisplayName: "Einride Batcave",
		LatLng: &latlng.LatLng{
			Latitude:  59.33749110496606,
			Longitude: 18.063672598779984,
		},
	}
}
