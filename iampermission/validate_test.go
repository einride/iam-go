package iampermission

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestValidate(t *testing.T) {
	for _, tt := range []struct {
		name          string
		permission    string
		errorContains string
	}{
		{
			name:       "valid",
			permission: "pubsub.subscriptions.consume",
		},

		{
			name:          "invalid service segment",
			permission:    "FOO.subscriptions.consume",
			errorContains: "invalid service segment",
		},

		{
			name:          "invalid number starting segment",
			permission:    "pubsub.42subscriptions.consume",
			errorContains: "invalid resource segment",
		},

		{
			name:          "invalid resource segment",
			permission:    "pubsub.SUBSCRIPTIONS.consume",
			errorContains: "invalid resource segment",
		},

		{
			name:       "camelCase segment",
			permission: "pubsub.awesomeStuff.consume",
		},

		{
			name:          "invalid verb segment",
			permission:    "pubsub.subscriptions.Consume",
			errorContains: "invalid verb segment",
		},

		{
			name:       "valid wildcard verb",
			permission: "pubsub.subscriptions.*",
		},

		{
			name:       "valid wildcard resource",
			permission: "pubsub.*",
		},

		{
			name:          "invalid wildcard service",
			permission:    "*",
			errorContains: "service segment must not be wildcard",
		},

		{
			name:          "all wildcards",
			permission:    "*.*.*",
			errorContains: "service segment must not be wildcard",
		},

		{
			name:          "multiple wildcards",
			permission:    "pubsub.*.*",
			errorContains: "only final segment can be wildcard",
		},

		{
			name:          "multiple wildcards",
			permission:    "pubsub.subscriptions.*.consume",
			errorContains: "only final segment can be wildcard",
		},

		{
			name:          "too many segments",
			permission:    "pubsub.subscriptions.consume.foobar",
			errorContains: "too many segments",
		},

		{
			name:          "missing verb",
			permission:    "pubsub.subscriptions",
			errorContains: "missing verb segment",
		},

		{
			name:          "missing resource",
			permission:    "pubsub",
			errorContains: "missing resource segment",
		},

		{
			name:          "empty",
			permission:    "",
			errorContains: "empty",
		},

		{
			name:          "empty resource",
			permission:    "pubsub..consume",
			errorContains: "invalid resource segment",
		},

		{
			name:          "empty service",
			permission:    "..consume",
			errorContains: "invalid service segment",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.permission)
			if tt.errorContains != "" {
				assert.ErrorContains(t, err, tt.errorContains)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
