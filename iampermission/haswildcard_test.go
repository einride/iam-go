package iampermission

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHasWildcard(t *testing.T) {
	for _, tt := range []struct {
		name       string
		permission string
		expected   bool
	}{
		{
			name:       "no wildcard",
			permission: "pubsub.subscriptions.consume",
			expected:   false,
		},

		{
			name:       "empty",
			permission: "",
			expected:   false,
		},

		{
			name:       "service wildcard",
			permission: "*.subscriptions.consume",
			expected:   true,
		},

		{
			name:       "resource wildcard",
			permission: "pubsub.*.consume",
			expected:   true,
		},

		{
			name:       "multiple wildcards",
			permission: "pubsub.*.*",
			expected:   true,
		},

		{
			name:       "only wildcard",
			permission: "*",
			expected:   true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, HasWildcard(tt.permission))
		})
	}
}
