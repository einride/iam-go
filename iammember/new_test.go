package iammember

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNew(t *testing.T) {
	assert.Equal(t, "", New("", ""))
	assert.Equal(t, "", New("a", ""))
	assert.Equal(t, "", New("", "b"))
	assert.Equal(t, "a:b", New("a", "b"))
}
