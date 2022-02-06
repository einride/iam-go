package examplecmd

import (
	"github.com/spf13/cobra"
	"go.einride.tech/iam/cmd/iamctl/internal/examplecmd/exampleclientcmd"
	"go.einride.tech/iam/cmd/iamctl/internal/examplecmd/exampleservercmd"
)

// Command for the IAM example server and client.
// nolint: gochecknoglobals
var Command = &cobra.Command{
	Use:   "example",
	Short: "IAM example server and client",
}

// nolint: gochecknoinits
func init() {
	Command.AddCommand(exampleclientcmd.Command)
	Command.AddCommand(exampleservercmd.Command)
	Command.Hidden = true
}
