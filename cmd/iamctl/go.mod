module go.einride.tech/iam/cmd/iamctl

go 1.16

require (
	cloud.google.com/go/spanner v1.24.1
	firebase.google.com/go/v4 v4.6.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	go.einride.tech/iam v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.54.0
	google.golang.org/genproto v0.0.0-20210805201207-89edb61ffb67
	google.golang.org/grpc v1.39.1
	google.golang.org/protobuf v1.27.1
)

replace go.einride.tech/iam => ../../
