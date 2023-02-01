package iampolicy

import (
	"testing"

	"cloud.google.com/go/iam/apiv1/iampb"
	"google.golang.org/protobuf/testing/protocmp"
	"gotest.tools/v3/assert"
)

func TestRemoveBinding(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		actual := &iampb.Policy{
			Bindings: []*iampb.Binding{
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
		expected := &iampb.Policy{
			Bindings: []*iampb.Binding{
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
