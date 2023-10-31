package cli

import (
	"github.com/TheLeeeo/grpc-hole/cli/servercmd"
	"github.com/TheLeeeo/grpc-hole/cli/servicecmd.go"
	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	// where to store all the services
	RootCmd.PersistentFlags().StringP(vars.SerivceDirKey, "d", ".services", "directory to store services")
	if err := viper.BindPFlag(vars.SerivceDirKey, RootCmd.PersistentFlags().Lookup(vars.SerivceDirKey)); err != nil {
		panic(err)
	}

	RootCmd.AddCommand(servicecmd.ServiceCmd)
	RootCmd.AddCommand(servercmd.ServerCmd)
}

var RootCmd = &cobra.Command{
	Use:   "grpc-hole",
	Short: "A dynamic grpc agent",
	Long:  `A dynamic grpc agent that can be used to test grpc clients and servers`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
