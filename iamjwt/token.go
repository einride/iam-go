package iamjwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Payload is a generic JWT token payload.
type Payload struct {
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Expires  int64  `json:"exp"`
	Subject  string `json:"sub"`
}

// UnmarshalToken unmarshals the provided JWT token.
// No claims verification or signature validation is performed.
func (t *Payload) UnmarshalToken(token string) error {
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
	return nil
}
