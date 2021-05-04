package examplecmd

import (
	"github.com/spf13/cobra"
	"go.einride.tech/iam/cmd/iamctl/examplecmd/exampleclientcmd"
	"go.einride.tech/iam/cmd/iamctl/examplecmd/exampleservercmd"
)

var Command = &cobra.Command{
	Use:   "example",
	Short: "IAM example server and client",
}

func init() {
	Command.AddCommand(exampleclientcmd.Command)
	Command.AddCommand(exampleservercmd.Command)
}
