module go.einride.tech/iam/cmd/iamctl

go 1.16

require (
	cloud.google.com/go/spanner v1.18.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.0
	go.einride.tech/iam v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.46.0
	google.golang.org/genproto v0.0.0-20210429181445-86c259c2b4ab
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace go.einride.tech/iam => ../../
