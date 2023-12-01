# IAM Go

An opinionated Open Source implementation of the
[google.iam.v1.IAMPolicy](https://github.com/googleapis/googleapis/blob/master/google/iam/v1/iam_policy.proto)
service API, using [Cloud Spanner](https://cloud.google.com/spanner) for
storage.

## Usage

### 1) Install

```bash
$ go get go.einride.tech/iam
```

### 2) Include the IAMPolicy mixin in your gRPC service

See
[google.iam.v1.IAMPolicy](https://github.com/googleapis/googleapis/blob/master/google/iam/v1/iam_policy.proto).

```proto
package your.pkg;

import "google/iam/v1/iam_policy.proto";
import "google/iam/v1/policy.proto";

service YourService {
  /* ... */

  rpc SetIamPolicy(google.iam.v1.SetIamPolicyRequest) returns (google.iam.v1.Policy);
  rpc GetIamPolicy(google.iam.v1.GetIamPolicyRequest) returns (google.iam.v1.Policy);
  rpc TestIamPermissions(google.iam.v1.TestIamPermissionsRequest) returns (google.iam.v1.TestIamPermissionsResponse);
}
```

### 3) Embed the IAMServer implementation in your server

See [iamspanner.IAMServer](./iamspanner/server.go).

```go
// Server implements your gRPC API.
type Server struct {
	*iamspanner.IAMServer
	// ...
}

// Server now also implements the iam.IAMPolicyServer mixin.
var _ iam.IAMPolicyServer = &Server{}
```

### 4) Include the IAM policy bindings table in your Spanner SQL schema

See [schema.sql](./iamspanner/schema.sql).

### 5) Annotate your gRPC methods

Buf annotations for rpc method authorization are described in
[annotations.proto](../proto/einride/iam/v1/annotations.proto)

```proto
package your.pkg;

import "einride/iam/v1/annotations.proto";

service YourService {
  rpc YourMethod(YourMethodRequest) returns YourMethodResponse {
      option (einride.iam.v1.method_authorization) = {
        permission: "namespace.entity.method"
        before: {
          expression: "test(caller, request.entity)" // iamcel expression
          description: "The caller must have method permission against the entity"
        }
      };
    };
}

message YourMethodRequest {
  string entity = 1 [
    (google.api.resource_reference) = {
      type: "example.com/Entity"
    }
  ];
};
```

```proto
package your.pkg;

import "einride/iam/v1/annotations.proto";

service YourService {
  rpc YourMethod(YourMethodRequest) returns YourMethodResponse {
      option (einride.iam.v1.method_authorization) = {
        resource_permissions {
          resource_permission {
            resource: {
              type: "example.com/Entity1"
            }
            permission: "namespace.entity1.method"
          }
          resource_permission {
            resource: {
              type: "example.com/Entity2"
            }
            permission: "namespace.entity2.method"
          }
        }
        after: {
          expression: "test_all(caller, response.entities)" // iamcel expression
          description: "The caller must have method permission against all entities"
        }
      };
    };
}

message YourMethodResponse {
  // Elements in this list are either Entity1 or Entity2 references
  repeated string entities = 1;
};
```

Expresssions in the `method_authorization` annotation use
[cel-go](https://github.com/google/cel-go) with [iamcel](./iamcel) extensions.
The `iamcel` extensions provide the following cel functions.

#### [`test(caller Caller, resource string) bool`](./iamcel/test.go)

Tests `caller`s permissions against `resource`.

#### [`test_all(caller Caller, resources []string) bool`](./iamcel/testall.go)

Tests `caller`s permissions against all `resources`. This test asserts that the
caller has the permission against all resources.

#### [`test_any(caller Caller, resources []string) bool`](./iamcel/testany.go)

Tests `caller`s permissions against any `resources`. This test asserts that the
caller has the permission against at least one resource.

#### [`ancestor(resource string, pattern string) string`](./iamcel/ancestor.go)

Resolves an ancestor of `resource` using `pattern`. An input of
`ancestor("foo/1/bar/2", "foo/{foo}")` will yield the result `"foo/1"`.

#### [`join(parent string, resource string) string`](./iamcel/join.go)

Joins a `resource` name with a `parent` resource name. An input of
`join("foo/1", "bar/2")` will yield the result `"foo/1/bar/2"`.

#### [`caller.member(kind string) string`](./iamcel/member.go)

Returns the first IAM member value from the caller's member list which matches
the member kind, or fails if there are no such kind.

### 6) Generate authorization middleware

Coming soon.
