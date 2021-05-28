package iamgooglemember

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// UserInfo from a Google ID token.
//
// See: https://developers.google.com/identity/protocols/oauth2/openid-connect
type UserInfo struct {
	// Issuer is an identifier for the Issuer of the response.
	// Always https://accounts.google.com or accounts.google.com for Google ID tokens.
	Issuer string `json:"iss,omitempty"`
	// ClientID of the authorized presenter.
	// This claim is only needed when the party requesting the ID token is not the same as the audience of the ID token.
	// This may be the case at Google for hybrid apps where a web application and Android app have a different OAuth 2.0
	// client ID but share the same Google APIs project.
	ClientID string `json:"azp,omitempty"`
	// Audience that this ID token is intended for.
	Audience string `json:"aud,omitempty"`
	// Subject is an identifier for the user, unique among all Google accounts and never reused.
	Subject string `json:"sub,omitempty"`
	// HostedDomain is the hosted G Suite domain of the user. Provided only if the user belongs to a hosted domain.
	HostedDomain string `json:"hd,omitempty"`
	// Email is user's email address. May be unset.
	Email string `json:"email,omitempty"`
	// EmailVerified is true if the user's e-mail address has been verified; otherwise false.
	EmailVerified bool `json:"email_verified,omitempty"`
	// AccessTokenHash provides validation that the access token is tied to the identity token.
	// If the ID token is issued with an access token value in the server flow, this claim is always included.
	// This claim can be used as an alternate mechanism to protect against cross-site request forgery attacks,
	// but if you use CSRF it is not necessary to verify the access token.
	AccessTokenHash string `json:"at_hash,omitempty"`
	// Name is the user's full name, in a displayable form.
	// When name claims are present, you can use them to update your app's user records.
	// Note that this claim is never guaranteed to be present.
	Name string `json:"name,omitempty"`
	// Picture is the URL of the user's profile picture.
	// When picture claims are present, you can use them to update your app's user records.
	// Note that this claim is never guaranteed to be present.
	Picture string `json:"picture,omitempty"`
	// GivenName is the user's given name(s) or first name(s). Might be provided when a name claim is present.
	GivenName string `json:"given_name,omitempty"`
	// FamilyName is the user's surname(s) or last name(s). Might be provided when a name claim is present.
	FamilyName string `json:"family_name,omitempty"`
	// The user's locale, represented by a BCP 47 language tag. Might be provided when a name claim is present.
	Locale string `json:"locale,omitempty"`
	// IssuedAt is the time the ID token was issued.
	// Represented in Unix time (integer seconds).
	IssuedAt int64 `json:"iat,omitempty"`
	// Expires is the expiration time on or after which the ID token must not be accepted.
	// Represented in Unix time (integer seconds).
	Expires int64 `json:"exp,omitempty"`
	// JWTID is the JWT ID of the ID token.
	JWTID string `json:"jti,omitempty"`
}

// Validate returns an error if the UserInfo is missing any required fields or has invalid values for known fields.
func (u *UserInfo) Validate() error {
	if u.Issuer == "" {
		return fmt.Errorf("validate user info: missing required field: 'iss'")
	}
	if u.Audience == "" {
		return fmt.Errorf("validate user info: missing required field: 'aud'")
	}
	if u.Expires == 0 {
		return fmt.Errorf("validate user info: missing required field: 'exp'")
	}
	if u.IssuedAt == 0 {
		return fmt.Errorf("validate user info: missing required field: 'iat'")
	}
	if u.Subject == "" {
		return fmt.Errorf("validate user info: missing required field: 'sub'")
	}
	if strings.TrimPrefix(u.Issuer, "https://") != Issuer {
		return fmt.Errorf("validate: unsupported issuer '%s'", u.Issuer)
	}
	return nil
}

// UnmarshalBase64 unmarshals the UserInfo from the provided Base64-URL-encoded string.
func (u *UserInfo) UnmarshalBase64(value string) error {
	decoder := json.NewDecoder(base64.NewDecoder(base64.URLEncoding, strings.NewReader(value)))
	if err := decoder.Decode(u); err != nil {
		return fmt.Errorf("unmarshal Google user info from base64: %w", err)
	}
	if err := u.Validate(); err != nil {
		return fmt.Errorf("unmarshal Google user info from base64: %w", err)
	}
	return nil
}

// UnmarshalJWT unmarshals the UserInfo from the provided JWT token.
func (u *UserInfo) UnmarshalJWT(token string) error {
	s := strings.Split(token, ".")
	if len(s) < 2 {
		return fmt.Errorf("unmarshal user info from JWT: invalid token")
	}
	if err := u.UnmarshalBase64(s[1]); err != nil {
		return fmt.Errorf("unmarshal user info from JWT: %w", err)
	}
	return nil
}
