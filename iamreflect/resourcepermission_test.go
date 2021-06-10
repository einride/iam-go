package iamreflect

import (
	"testing"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
	"google.golang.org/genproto/googleapis/api/annotations"
	"gotest.tools/v3/assert"
)

func TestResolveMethodPermission(t *testing.T) {
	t.Run("permission", func(t *testing.T) {
		options := &iamv1.MethodAuthorizationOptions{
			Permissions: &iamv1.MethodAuthorizationOptions_Permission{
				Permission: "foo.bar.baz",
			},
		}
		result, ok := ResolveMethodPermission(options, "foos/foo")
		assert.Assert(t, ok)
		assert.Equal(t, "foo.bar.baz", result)
	})

	t.Run("resource permissions", func(t *testing.T) {
		options := &iamv1.MethodAuthorizationOptions{
			Permissions: &iamv1.MethodAuthorizationOptions_ResourcePermissions{
				ResourcePermissions: &iamv1.ResourcePermissions{
					ResourcePermission: []*iamv1.ResourcePermission{
						{
							Resource: &annotations.ResourceDescriptor{
								Type: "foo.test.com/Test",
								Pattern: []string{
									"foos/{foo}",
								},
							},
							Permission: "foo.bar.baz",
						},
					},
				},
			},
		}
		result, ok := ResolveMethodPermission(options, "foos/foo")
		assert.Assert(t, ok)
		assert.Equal(t, "foo.bar.baz", result)
	})
}
