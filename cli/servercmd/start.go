package servercmd

import (
	"github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start a server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
