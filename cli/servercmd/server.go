package servercmd

import (
	"github.com/spf13/cobra"
)

func init() {
	ServerCmd.AddCommand(StartCmd)
	ServerCmd.AddCommand(CreateDefaultsCmd)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
