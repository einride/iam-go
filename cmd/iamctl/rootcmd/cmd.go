package rootcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.einride.tech/iam/cmd/iamctl/examplecmd"
)

// Command is the CLI root command.
var Command = &cobra.Command{
	Use:   "iamctl",
	Short: "IAM command line controls",
}

// Execute the CLI root command.
func Execute() {
	if err := Command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	Command.AddCommand(examplecmd.Command)
}
