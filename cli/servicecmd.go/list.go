package servicecmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all saved services",
	Run: func(cmd *cobra.Command, args []string) {
		serviceDirs, err := os.ReadDir(viper.GetString(vars.SerivceDirKey))
		if err != nil {
			color.Red(fmt.Errorf("failed to load service directory: %w", err).Error())
			os.Exit(1)
		}

		for _, f := range serviceDirs {
			if f.IsDir() {
				dataFilePath := filepath.Join(viper.GetString(vars.SerivceDirKey), f.Name())
				serviceData, err := service.LoadDataFile(dataFilePath)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}

					color.Red(fmt.Errorf("failed to load service data for service %s: error=%w", f.Name(), err).Error())
					os.Exit(1)
				}

				serviceName := protoreflect.FullName(serviceData.Name)
				fmt.Println(serviceName)
			}
		}
	},
}
