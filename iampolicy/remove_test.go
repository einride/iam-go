package iampolicy

import (
	"testing"

	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestRemoveBinding(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		actual := &iam.Policy{
			Bindings: []*iam.Binding{
				{
					Role:    "roles/test",
					Members: []string{"foo", "bar"},
				},
				{
					Role:    "roles/test2",
					Members: []string{"foo", "bar"},
				},
			},
		}
		RemoveBinding(actual, "roles/test2", "bar")
		expected := &iam.Policy{
			Bindings: []*iam.Binding{
				{
					Role:    "roles/test",
					Members: []string{"foo", "bar"},
				},
				{
					Role:    "roles/test2",
					Members: []string{"foo"},
				},
			},
		}
		assert.DeepEqual(t, expected, actual, protocmp.Transform())
	})
}
