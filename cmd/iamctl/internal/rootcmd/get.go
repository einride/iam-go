package rootcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// nolint: gochecknoglobals
var getIAMPolicyCommand = &cobra.Command{
	Use:   "get-iam-policy <resource>",
	Short: "Get an IAM policy",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags getIAMPolicyFlags
		if err := viperCfg.Unmarshal(&flags); err != nil {
			return err
		}
		conn, err := flags.Connect(cmd.Context())
		if err != nil {
			return err
		}
		defer func() {
			_ = conn.Close()
		}()
		client := iam.NewIAMPolicyClient(conn)
		return runGetIAMPolicyCommand(cmd.Context(), client, args[0])
	},
}

type getIAMPolicyFlags struct {
	connection.Flags `mapstructure:",squash"`
}

func runGetIAMPolicyCommand(ctx context.Context, client iam.IAMPolicyClient, resource string) error {
	policy, err := client.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: resource,
		Options:  nil,
	})
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(policy))
	return nil
}
