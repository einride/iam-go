package iamtoken

import (
	"encoding/base64"
	"fmt"
	"strings"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// ParseIdentityToken parses a JWT identity token.
func ParseIdentityToken(token string) (*iamv1.IdentityToken, error) {
	s := strings.Split(token, ".")
	if len(s) < 2 {
		return nil, fmt.Errorf("parse identity token: not a valid token")
	}
	payloadData, err := base64.RawURLEncoding.DecodeString(s[1])
	if err != nil {
		return nil, fmt.Errorf("parse identity token: %w", err)
	}
	var result iamv1.IdentityToken
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(payloadData, &result); err != nil {
		return nil, fmt.Errorf("parse identity token: %w", err)
	}
	result.Raw = token
	return &result, nil
}
