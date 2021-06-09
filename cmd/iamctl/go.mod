module go.einride.tech/iam/cmd/iamctl

go 1.16

require (
	cloud.google.com/go/spanner v1.20.0
	firebase.google.com/go/v4 v4.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.einride.tech/iam v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.47.0
	google.golang.org/genproto v0.0.0-20210604141403-392c879c8b08
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace go.einride.tech/iam => ../../
