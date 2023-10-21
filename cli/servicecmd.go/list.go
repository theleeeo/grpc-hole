package servicecmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all saved services",
	Run: func(cmd *cobra.Command, args []string) {
		serviceDir, err := os.ReadDir(viper.GetString(vars.SerivceDirKey))
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, f := range serviceDir {
			if f.IsDir() {
				dataFilePath := filepath.Join(viper.GetString(vars.SerivceDirKey), f.Name())
				serviceData, err := service.LoadDataFile(dataFilePath)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}

					fmt.Println(err)
					return
				}

				serviceName := protoreflect.FullName(serviceData.Name)
				fmt.Println(serviceName.Name())
			}
		}
	},
}
