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

var addIAMPolicyBindingCommand = &cobra.Command{
	Use:   "add-iam-policy-binding <resource>",
	Short: "Add an IAM policy binding",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags addIAMPolicyBindingFlags
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
		return runAddIAMPolicyBindingCommand(cmd.Context(), client, args[0], &flags)
	},
}

type addIAMPolicyBindingFlags struct {
	connection.Flags `mapstructure:",squash"`
	Member           string `mapstructure:"member"`
	Role             string `mapstructure:"role"`
}

func init() {
	addIAMPolicyBindingCommand.Flags().String("member", "", "member to add")
	addIAMPolicyBindingCommand.Flags().String("role", "", "role to add")
	_ = addIAMPolicyBindingCommand.MarkFlagRequired("member")
	_ = addIAMPolicyBindingCommand.MarkFlagRequired("role")
}

func runAddIAMPolicyBindingCommand(
	ctx context.Context,
	client iam.IAMPolicyClient,
	resource string,
	flags *addIAMPolicyBindingFlags,
) error {
	policy, err := client.GetIamPolicy(ctx, &iam.GetIamPolicyRequest{
		Resource: resource,
		Options:  nil,
	})
	if err != nil {
		return err
	}
	iampolicy.AddBinding(policy, flags.Role, flags.Member)
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
