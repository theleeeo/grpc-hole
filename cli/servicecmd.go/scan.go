package servicecmd

import (
	"fmt"
	"os"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/scanning"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	targetAddr string
)

func init() {
	ScanCmd.PersistentFlags().StringVarP(&targetAddr, "target", "t", "", "address of the target grpc server to scan")
}

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan a running grpc server for services",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := scanning.ScanServer(targetAddr)
		if err != nil {
			color.Red(fmt.Errorf("failed to scan service: %w", err).Error())
			os.Exit(1)
		}

		fmt.Println("Found services:")
		for _, s := range services {
			fmt.Println(s.GetFullyQualifiedName())
		}

		var errors []error
		for _, s := range services {
			err := service.Save(viper.GetString(vars.SerivceDirKey), s)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			color.Red("failed to save services:")
			for _, err := range errors {
				fmt.Println(err.Error())
			}
			os.Exit(1)
		}
	},
}
