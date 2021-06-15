package iammember

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParse(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		kind, value, ok := Parse("a:b")
		assert.Assert(t, ok)
		assert.Equal(t, "a", kind)
		assert.Equal(t, "b", value)
	})

	t.Run("empty", func(t *testing.T) {
		kind, value, ok := Parse("")
		assert.Assert(t, !ok)
		assert.Equal(t, "", kind)
		assert.Equal(t, "", value)
	})

	t.Run("missing kind", func(t *testing.T) {
		kind, value, ok := Parse(":b")
		assert.Assert(t, !ok)
		assert.Equal(t, "", kind)
		assert.Equal(t, "", value)
	})

	t.Run("missing value", func(t *testing.T) {
		kind, value, ok := Parse("a:")
		assert.Assert(t, !ok)
		assert.Equal(t, "", kind)
		assert.Equal(t, "", value)
	})
}
