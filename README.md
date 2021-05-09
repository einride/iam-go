# IAM Go

An opinionated Open Source implementation of the
[google.iam.v1.IAMPolicy][iam-policy] service API, using [Cloud
Spanner][cloud-spanner] for storage.

[iam-policy]: https://github.com/googleapis/googleapis/blob/master/google/iam/v1/iam_policy.proto
[cloud-spanner]: https://cloud.google.com/spanner

## Usage

### 1) Install

```bash
$ go get go.einride.tech/iam
```

### 2) Include the IAMPolicy mixin in your gRPC service

See [google.iam.v1.IAMPolicy][iam-policy-proto].

```proto
/* ... */
package your.pkg;

/* ... */

import "google/iam/v1/iam_policy.proto";
import "google/iam/v1/policy.proto";

/* ... */

service YourService {
  /* ... */

  rpc SetIamPolicy(google.iam.v1.SetIamPolicyRequest)
    returns (google.iam.v1.Policy);
  rpc GetIamPolicy(google.iam.v1.GetIamPolicyRequest)
    returns (google.iam.v1.Policy);
  rpc TestIamPermissions(google.iam.v1.TestIamPermissionsRequest)
    returns (google.iam.v1.TestIamPermissionsResponse);
}
```

[iam-policy-proto]: https://github.com/googleapis/googleapis/blob/master/google/iam/v1/iam_policy.proto

### 3) Embed the IAMServer implementation in your server

See [iamspanner.IAMServer][iam-server].

[iam-server]: ./iamspanner/server.go

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

See [schema.sql][schema-sql].

[schema-sql]: ./iamspanner/schema.sql

### 5) Annotate your gRPC methods

Coming soon.

### 6) Generate authorization middleware

Coming soon.
