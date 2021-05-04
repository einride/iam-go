package exampleclientcmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.einride.tech/iam/iamexample"
	iamexamplev1 "go.einride.tech/iam/proto/gen/einride/iam/example/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var Command = &cobra.Command{
	Use:   "client",
	Short: "IAM example client",
}

type commandConfig struct {
	Address string   `mapstructure:"address"`
	Members []string `mapstructure:"members"`
}

func init() {
	// shippers
	Command.AddCommand(createShipperCommand)
	Command.AddCommand(deleteShipperCommand)
	Command.AddCommand(getShipperCommand)
	Command.AddCommand(listShippersCommand)
	Command.AddCommand(updateShipperCommand)
	// sites
	Command.AddCommand(batchGetSitesCommand)
	Command.AddCommand(createSiteCommand)
	Command.AddCommand(deleteSiteCommand)
	Command.AddCommand(getSiteCommand)
	Command.AddCommand(listSitesCommand)
	Command.AddCommand(searchSitesCommand)
	Command.AddCommand(updateSiteCommand)
	// shipments
	Command.AddCommand(batchGetShipmentsCommand)
	Command.AddCommand(createShipmentCommand)
	Command.AddCommand(deleteShipmentCommand)
	Command.AddCommand(getShipmentCommand)
	Command.AddCommand(listShipmentsCommand)
	Command.AddCommand(updateShipmentCommand)
	Command.PersistentFlags().String("address", "localhost:8080", "address to connect to")
	Command.PersistentFlags().StringSlice("members", nil, "IAM members to set for the caller")
}

func (c *commandConfig) connect(ctx context.Context) (iamexamplev1.FreightServiceClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		c.Address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(c.unaryClientInterceptor),
	)
	if err != nil {
		return nil, err
	}
	return iamexamplev1.NewFreightServiceClient(conn), nil
}

func (c *commandConfig) unaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if err := invoker(iamexample.WithOutgoingMembers(ctx, c.Members...), method, req, reply, cc, opts...); err != nil {
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
