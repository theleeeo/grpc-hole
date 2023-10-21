package servicecmd

import (
	"github.com/spf13/cobra"
)

func init() {
	ServiceCmd.AddCommand(ScanCmd)
	ServiceCmd.AddCommand(ListCmd)
	ServiceCmd.AddCommand(InfoCmd)
}

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
