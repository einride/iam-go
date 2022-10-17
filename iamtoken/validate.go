package iamtoken

import (
	"fmt"
	"time"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// timeNow is for mocking in tests.
var timeNow = time.Now //nolint: gochecknoglobals

// ValidateIdentityToken checks that an identity token is valid and acceptable for processing.
func ValidateIdentityToken(token *iamv1.IdentityToken) error {
	if token.GetExp() > 0 {
		now := timeNow()
		expired := time.Unix(token.GetExp(), 0).UTC()
		if expired.Before(now) {
			return fmt.Errorf(
				"identity token expired at %s (%s ago)",
				expired.Format(time.RFC3339),
				now.Sub(expired),
			)
		}
	}
	if token.GetNbf() > 0 {
		now := timeNow()
		notBefore := time.Unix(token.GetNbf(), 0).UTC()
		if now.Before(notBefore) {
			return fmt.Errorf(
				"identity token not valid before %s (in %s)",
				notBefore.Format(time.RFC3339),
				notBefore.Sub(now),
			)
		}
	}
	return nil
}
