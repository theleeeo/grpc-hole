package servercmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	ServerCmd.PersistentFlags().StringP("port", "p", "", "the port to listen on")
	if err := viper.BindPFlag("port", ServerCmd.PersistentFlags().Lookup("port")); err != nil {
		panic(err)
	}

	ServerCmd.PersistentFlags().String("host", "0.0.0.0", "the host to listen on")
	if err := viper.BindPFlag("host", ServerCmd.PersistentFlags().Lookup("host")); err != nil {
		panic(err)
	}

	ServerCmd.PersistentFlags().StringP("log-level", "l", "info", "[trace|debug|info|warn|error|off]")
	if err := viper.BindPFlag("log-level", ServerCmd.PersistentFlags().Lookup("log-level")); err != nil {
		panic(err)
	}

	ServerCmd.AddCommand(StaticCmd)
	ServerCmd.AddCommand(ProxyCmd)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
