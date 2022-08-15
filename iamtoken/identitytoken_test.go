package iamtoken

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseIdentityToken(t *testing.T) {
	//nolint: lll
	const input = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	actual, err := ParseIdentityToken(input)
	assert.NilError(t, err)
	assert.Equal(t, "1234567890", actual.Sub)
	assert.Equal(t, int64(1516239022), actual.Iat)
}
