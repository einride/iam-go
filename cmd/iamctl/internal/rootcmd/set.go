package rootcmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// nolint: gochecknoglobals
var setIAMPolicyCommand = &cobra.Command{
	Use:   "set-iam-policy <resource> <policy-file>",
	Short: "Set an IAM policy",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags setIAMPolicyFlags
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
		return runSetIAMPolicyCommand(cmd.Context(), client, args[0], args[1])
	},
}

type setIAMPolicyFlags struct {
	connection.Flags `mapstructure:",squash"`
}

func runSetIAMPolicyCommand(
	ctx context.Context,
	client iam.IAMPolicyClient,
	resource string,
	policyFile string,
) error {
	policyJSON, err := os.ReadFile(policyFile)
	if err != nil {
		return err
	}
	var policy iam.Policy
	if err := protojson.Unmarshal(policyJSON, &policy); err != nil {
		return err
	}
	newPolicy, err := client.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
		Resource: resource,
		Policy:   &policy,
	})
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(newPolicy))
	return nil
}
