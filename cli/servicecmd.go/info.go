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
)

var InfoCmd = &cobra.Command{
	Use:   "info [service]",
	Short: "show info about a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no service specified")
		}

		if len(args) > 1 {
			return (fmt.Errorf("too many arguments"))
		}

		sericeDir := viper.GetString(vars.SerivceDirKey)
		serviceData, err := service.LoadDataFile(filepath.Join(sericeDir, args[0]))
		if err != nil {
			color.Red(fmt.Errorf("failed to load service: %w", err).Error())
			os.Exit(1)
		}

		service, err := service.Load(sericeDir, args[0])
		if err != nil {
			color.Red(fmt.Errorf("failed to load service: %w", err).Error())
			os.Exit(1)
		}

		fmt.Printf("Full name: %s\n", serviceData.Name)
		fmt.Printf("Saved at: %s\n", serviceData.SavedAt)

		fmt.Printf("Methods:\n")
		for _, method := range service.GetMethods() {
			fmt.Printf("  %s\n", method.GetName())
		}

		return nil
	},
}
