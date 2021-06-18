package iamtoken

import (
	"testing"
	"time"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"gotest.tools/v3/assert"
)

func TestValidateIdentityToken(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		mockTimeNow(t, 1623998472)
		identityToken := &iamv1.IdentityToken{
			Exp: 1623998480,
			Nbf: 1623998470,
		}
		assert.NilError(t, ValidateIdentityToken(identityToken))
	})

	t.Run("after exp", func(t *testing.T) {
		mockTimeNow(t, 1623998472)
		identityToken := &iamv1.IdentityToken{
			Exp: 1623998471,
			Nbf: 1623998470,
		}
		assert.Error(
			t,
			ValidateIdentityToken(identityToken),
			"identity token expired at 2021-06-18T06:41:11Z (1s ago)",
		)
	})

	t.Run("before nbf", func(t *testing.T) {
		mockTimeNow(t, 1623998472)
		identityToken := &iamv1.IdentityToken{
			Exp: 1623998480,
			Nbf: 1623998473,
		}
		assert.Error(
			t,
			ValidateIdentityToken(identityToken),
			"identity token not valid before 2021-06-18T06:41:13Z (in 1s)",
		)
	})
}

func mockTimeNow(t *testing.T, tt int64) {
	timeNow = func() time.Time {
		return time.Unix(tt, 0).UTC()
	}
	t.Cleanup(func() {
		timeNow = time.Now
	})
}
