package servicecmd

import (
	"fmt"
	"os"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/fatih/color"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	targetFiles []string
	protoRoot   string
)

func init() {
	LoadCmd.PersistentFlags().StringSliceVarP(&targetFiles, "files", "f", []string{}, "the files containing the .proto files to load")
	LoadCmd.PersistentFlags().StringVarP(&protoRoot, "root", "r", ".", "the root directory of the .proto files and its dependencies")
}

var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "load a service from a .proto file",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		parser := protoparse.Parser{
			// If your .proto files include any imports, you need to provide the paths to those imports.
			// Include paths are necessary for the parser to locate and parse imported files referenced in the .proto file.
			ImportPaths:                     []string{protoRoot},
			ValidateUnlinkedFiles:           true,
			InterpretOptionsInUnlinkedFiles: true,
		}

		// The names of the files you want to parse.
		filenames := targetFiles

		// Parse the files into descriptors.
		fds, err := parser.ParseFiles(filenames...)
		if err != nil {
			color.Red("Failed to parse .proto files: %v", err)
			os.Exit(1)
		}

		// fds will be a slice of *desc.FileDescriptor. You can range over it to access each descriptor.
		services := make([]*desc.ServiceDescriptor, 0)
		for _, fd := range fds {
			services = append(services, fd.GetServices()...)
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
