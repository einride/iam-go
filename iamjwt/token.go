package iamjwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"strings"
)

// Token is a generic JWT token.
type Token struct {
	// Checksum is a checksum of the raw token, including header and signature.
	Checksum uint64
	// Issuer is the "iss" claim, that identifies the principal that issued the JWT.
	// The processing of this claim is generally application specific.
	// The "iss" value is a case-sensitive string containing a string or URI value.
	Issuer string `json:"iss"`
	// Audience is the "aud" claim identifies the recipients that the JWT is intended for.
	// Each principal intended to process the JWT MUST identify itself with a value in the audience claim.
	// If the principal processing the claim does not identify itself with a value in the
	// "aud" claim when this claim is present, then the JWT MUST be rejected.
	// In the general case, the "aud" value is an array of case-sensitive strings, each containing a
	// string or URI value.
	// In the special case when the JWT has one audience, the "aud" value MAY be a
	// single case-sensitive string containing a string or URI value.
	// The interpretation of audience values is generally application specific.
	Audience string `json:"aud"`
	// Expires is the "exp" claim that identifies the expiration time on or after
	// which the JWT MUST NOT be accepted for processing.
	// The processing of the "exp" claim requires that the current date/time
	// MUST be before the expiration date/time listed in the "exp" claim.
	Expires int64 `json:"exp"`
	// IssuedAt is the "iat" claim that identifies the time at which the JWT was
	// issued.  This claim can be used to determine the age of the JWT.
	// Its value MUST be a number containing a numeric date value.
	IssuedAt int64 `json:"iat"`
	// Subject is the "sub" claim that identifies the principal that is the
	// subject of the JWT.  The claims in a JWT are normally statements
	// about the subject.  The subject value MUST either be scoped to be
	// locally unique in the context of the issuer or be globally unique.
	// The processing of this claim is generally application specific.
	// The "sub" value is a case-sensitive string containing a string or URI value.
	Subject string `json:"sub"`
}

// crc64ISOTable used for calculating token checksums.
var crc64ISOTable = crc64.MakeTable(crc64.ISO)

// UnmarshalString unmarshals the JWT token from the provided string.
// No claims verification or signature validation is performed.
func (t *Token) UnmarshalString(token string) error {
	s := strings.Split(token, ".")
	if len(s) < 2 {
		return fmt.Errorf("unmarshal JWT payload: invalid token")
	}
	payloadData, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return fmt.Errorf("unmarshal JWT payload: %w", err)
	}
	if err := json.Unmarshal(payloadData, t); err != nil {
		return fmt.Errorf("unmarshal JWT payload: %w", err)
	}
	t.Checksum = crc64.Checksum([]byte(token), crc64ISOTable)
	return nil
}
