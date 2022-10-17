package iampermission

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestMatch(t *testing.T) {
	for _, tt := range []struct {
		name     string
		lhs      string
		rhs      string
		expected bool
	}{
		{
			name:     "equal",
			lhs:      "pubsub.subscriptions.consume",
			rhs:      "pubsub.subscriptions.consume",
			expected: true,
		},

		{
			name:     "wildcard verb",
			lhs:      "pubsub.subscriptions.*",
			rhs:      "pubsub.subscriptions.consume",
			expected: true,
		},

		{
			name:     "wildcard resource",
			lhs:      "pubsub.*",
			rhs:      "pubsub.subscriptions.consume",
			expected: true,
		},

		{
			name:     "wildcard service not allowed",
			lhs:      "*",
			rhs:      "pubsub.subscriptions.consume",
			expected: false,
		},

		{
			name:     "wildcard both",
			lhs:      "pubsub.*",
			rhs:      "pubsub.*",
			expected: false,
		},

		{
			name:     "wildcard rhs",
			lhs:      "pubsub.subscriptions.consume",
			rhs:      "pubsub.*",
			expected: false,
		},

		{
			name:     "non-matching service",
			lhs:      "pubsub.subscriptions.consume",
			rhs:      "foo.subscriptions.consume",
			expected: false,
		},

		{
			name:     "non-matching resource",
			lhs:      "pubsub.subscriptions.consume",
			rhs:      "pubsub.foo.consume",
			expected: false,
		},

		{
			name:     "non-matching verb",
			lhs:      "pubsub.subscriptions.consume",
			rhs:      "pubsub.subscriptions.foo",
			expected: false,
		},

		{
			name:     "camelCase",
			lhs:      "pubsub.fooSubscriptions.create",
			rhs:      "pubsub.fooSubscriptions.create",
			expected: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Match(tt.lhs, tt.rhs))
		})
	}
}

var boolSink bool //nolint: gochecknoglobals

func BenchmarkMatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boolSink = Match("pubsub.subscriptions.*", "pubsub.subscriptions.consume")
	}
}
