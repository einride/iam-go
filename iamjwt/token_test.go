package iamjwt

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestToken_UnmarshalString(t *testing.T) {
	const input = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	var actual Token
	assert.NilError(t, actual.UnmarshalString(input))
	assert.Equal(t, "1234567890", actual.Subject)
	assert.Equal(t, int64(1516239022), actual.IssuedAt)
	assert.Equal(t, uint64(3411411760009329079), actual.Checksum)
	t.Log(actual)
}
