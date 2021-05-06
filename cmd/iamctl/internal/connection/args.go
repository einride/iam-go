package connection

import (
	"context"
	"crypto/x509"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"go.einride.tech/iam/iamexample"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Flags struct {
	Address        string   `mapstructure:"address"`
	Insecure       bool     `mapstructure:"insecure"`
	Token          string   `mapstructure:"token"`
	ExampleMembers []string `mapstructure:"example-members"`
}

func AddToFlagSet(flagSet *pflag.FlagSet) error {
	flagSet.String("address", "localhost:8080", "address to connect to")
	flagSet.String("token", "", "bearer token used by the client")
	flagSet.Bool("insecure", false, "make insecure connection")
	flagSet.StringSlice("example-members", nil, "example IAM members to set for the caller")
	if err := flagSet.MarkHidden("example-members"); err != nil {
		return err
	}
	return nil
}

func (f *Flags) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithUnaryInterceptor(f.unaryClientInterceptor))
	if f.Token != "" {
		audience := "https://" + trimPort(f.Address)
		idTokenSource, err := idtoken.NewTokenSource(ctx, audience)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithPerRPCCredentials(
			&oauth.TokenSource{TokenSource: idTokenSource},
		))
	}
	if f.Insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		systemCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(systemCertPool, "")))
	}
	const tlsPort = 443
	return grpc.DialContext(ctx, withDefaultPort(f.Address, tlsPort), opts...)
}

func (c *Flags) unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if len(c.ExampleMembers) > 0 {
		ctx = iamexample.WithOutgoingMembers(ctx, c.ExampleMembers...)
	}
	if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
		return &printDetailsError{err: err}
	}
	return nil
}

type printDetailsError struct {
	err error
}

func (e *printDetailsError) Error() string {
	s, ok := status.FromError(e.err)
	if !ok {
		return e.err.Error()
	}
	details := s.Details()
	if len(details) == 0 {
		return e.err.Error()
	}
	var result strings.Builder
	_, _ = result.WriteString(e.err.Error())
	for _, details := range details {
		_ = result.WriteByte('\n')
		if protoDetails, ok := details.(proto.Message); ok {
			_, _ = result.WriteString(protojson.Format(protoDetails))
		} else {
			_, _ = result.WriteString(fmt.Sprintf("%v", details))
		}
	}
	return result.String()
}

func trimPort(target string) string {
	parts := strings.Split(target, ":")
	if len(parts) == 1 {
		return target
	}
	return strings.Join(parts[:len(parts)-1], ":")
}

func withDefaultPort(target string, port int) string {
	parts := strings.Split(target, ":")
	if len(parts) == 1 {
		return target + ":" + strconv.Itoa(port)
	}
	return target
}
