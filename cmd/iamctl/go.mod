module go.einride.tech/iam/cmd/iamctl

go 1.16

require (
	cloud.google.com/go/spanner v1.22.0
	firebase.google.com/go/v4 v4.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.einride.tech/iam v0.0.0-00010101000000-000000000000
	google.golang.org/api v0.49.0
	google.golang.org/genproto v0.0.0-20210624174822-c5cf32407d0a
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
)

replace go.einride.tech/iam => ../../
