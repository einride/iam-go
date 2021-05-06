package rootcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	"go.einride.tech/iam/iampolicy"
	"google.golang.org/genproto/googleapis/iam/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var removeIAMPolicyBindingCommand = &cobra.Command{
	Use:   "remove-iam-policy-binding <resource>",
	Short: "Remove an IAM policy binding",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags removeIAMPolicyBindingFlags
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
		return runRemoveIAMPolicyBindingCommand(cmd.Context(), client, args[0], &flags)
	},
}

type removeIAMPolicyBindingFlags struct {
	connection.Flags `mapstructure:",squash"`
	Member           string `mapstructure:"member"`
	Role             string `mapstructure:"role"`
}

func init() {
	removeIAMPolicyBindingCommand.Flags().String("member", "", "member to remove")
	removeIAMPolicyBindingCommand.Flags().String("role", "", "role to remove")
	_ = removeIAMPolicyBindingCommand.MarkFlagRequired("member")
	_ = removeIAMPolicyBindingCommand.MarkFlagRequired("role")
}

func runRemoveIAMPolicyBindingCommand(
	ctx context.Context,
	client iam.IAMPolicyClient,
	resource string,
	flags *removeIAMPolicyBindingFlags,
) error {
	policy, err := client.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: resource,
		Options:  nil,
	})
	if err != nil {
		return err
	}
	iampolicy.RemoveBinding(policy, flags.Role, flags.Member)
	newPolicy, err := client.SetIamPolicy(ctx, &iam.SetIamPolicyRequest{
		Resource: resource,
		Policy:   policy,
	})
	if err != nil {
		return err
	}
	fmt.Println(protojson.Format(newPolicy))
	return nil
}
