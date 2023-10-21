package servercmd

import (
	"fmt"
	"os"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/server"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	StartCmd.Flags().StringP("addr", "a", "", "the host:port to listen on")
	viper.BindPFlag("addr", StartCmd.Flags().Lookup("addr"))

	StartCmd.Flags().StringP("log-level", "l", "info", "[trace|debug|info|warn|error|off]")
	viper.BindPFlag("log-level", StartCmd.Flags().Lookup("log-level"))
}

var StartCmd = &cobra.Command{
	Use:   "start [service]",
	Short: "start a server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// When an error with the command syntax occurs, return the error.
		// This will print the error and the usage information.
		// If it is a runtime error, print the error and exit with code 1.
		if len(args) == 0 {
			return fmt.Errorf("no service specified")
		}

		if len(args) > 1 {
			return (fmt.Errorf("too many arguments"))
		}

		level := hclog.LevelFromString(viper.GetString("log-level"))
		if level == hclog.NoLevel {
			return fmt.Errorf("invalid log level: %s", viper.GetString("log-level"))
		}

		logger := hclog.New(&hclog.LoggerOptions{
			Name:  "grpc-hole",
			Level: level,
			Color: hclog.AutoColor,
		})

		serviceDescr, err := service.Load(viper.GetString(vars.SerivceDirKey), args[0])
		if err != nil {
			color.Red(fmt.Errorf("failed to load service: %w", err).Error())
			os.Exit(1)
		}

		cfg := &server.Config{
			Address: viper.GetString("addr"),
			Service: serviceDescr,
		}

		s, err := server.New(cfg, logger)
		if err != nil {
			color.Red(fmt.Errorf("failed to create server: %w", err).Error())
			os.Exit(1)
		}

		if err = s.Run(); err != nil {
			color.Red(fmt.Errorf("failed to run server: %w", err).Error())
			os.Exit(1)
		}

		return nil
	},
}
