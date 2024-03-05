package iamgoogle

import (
	"strings"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// IsSignatureRemoved checks if the ID token's signature has been removed by Google.
// See: https://cloud.google.com/run/docs/troubleshooting#signature-removed
func IsSignatureRemoved(identityToken *iamv1.IdentityToken) bool {
	return strings.HasSuffix(identityToken.GetRaw(), ".SIGNATURE_REMOVED_BY_GOOGLE")
}
