package servercmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/methodhandler"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/fatih/color"
	"github.com/jhump/protoreflect/desc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateDefaultsCmd = &cobra.Command{
	Use:   "create-defaults [service]",
	Short: "create default response files for a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("no service specified")
		}

		if len(args) > 1 {
			return (fmt.Errorf("too many arguments"))
		}

		path := filepath.Join(viper.GetString(vars.SerivceDirKey), args[0])
		service, err := service.Load(path)
		if err != nil {
			color.Red(fmt.Errorf("failed to load service: %w", err).Error())
			os.Exit(1)
		}

		serviceDir := filepath.Join(viper.GetString(vars.SerivceDirKey), args[0])
		for _, method := range service.GetMethods() {
			err := CreateDefaultResponseFile(method, serviceDir)
			if err != nil {
				color.Red(fmt.Errorf("failed to create default response file for method %s: %w", method.GetName(), err).Error())
			}
		}

		return nil
	},
}

func CreateDefaultResponseFile(method *desc.MethodDescriptor, methodDir string) error {
	msg := methodhandler.CreatePopulatedMessage(method.GetOutputType())
	b, err := msg.MarshalJSON()
	if err != nil {
		return err
	}

	return service.SaveResponseFile(methodDir, method.GetName(), b)
}
