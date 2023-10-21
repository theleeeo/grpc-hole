package servicecmd

import (
	"fmt"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/scanning"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	address string
)

func init() {
	ScanCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "address of the grpc server")
}

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan a running grpc server for services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := scanning.ScanServer(address)
		if err != nil {
			return err
		}

		fmt.Println("Found services:")
		for _, s := range services {
			fmt.Println(s.GetName())
		}

		for _, s := range services {
			service.Save(viper.GetString(vars.SerivceDirKey), s)
		}

		return nil
	},
}
