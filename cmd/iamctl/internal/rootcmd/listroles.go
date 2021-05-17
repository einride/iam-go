package rootcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.einride.tech/iam/cmd/iamctl/internal/connection"
	"google.golang.org/genproto/googleapis/iam/admin/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

var listRolesCommand = &cobra.Command{
	Use:   "list-roles",
	Short: "List IAM roles",
	RunE: func(cmd *cobra.Command, args []string) error {
		viperCfg := viper.New()
		if err := viperCfg.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		if err := viperCfg.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}
		var flags listRolesFlags
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
		client := admin.NewIAMClient(conn)
		return runListRolesCommand(cmd.Context(), client, &flags)
	},
}

type listRolesFlags struct {
	connection.Flags `mapstructure:",squash"`
	Full             bool `mapstructure:"full"`
}

func init() {
	listRolesCommand.Flags().Bool("full", false, "list full roles")
}

func runListRolesCommand(
	ctx context.Context,
	client admin.IAMClient,
	flags *listRolesFlags,
) error {
	var nextPageToken string
	var view admin.RoleView
	if flags.Full {
		view = admin.RoleView_FULL
	}
	for {
		response, err := client.ListRoles(ctx, &admin.ListRolesRequest{
			PageToken: nextPageToken,
			View:      view,
		})
		if err != nil {
			return err
		}
		for _, role := range response.Roles {
			fmt.Println(protojson.Format(role))
		}
		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}
	return nil
}
