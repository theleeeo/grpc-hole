package servercmd

import (
	"fmt"
	"os"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/runner"
	"github.com/TheLeeeo/grpc-hole/server"
	"github.com/fatih/color"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	StaticCmd.AddCommand(CreateDefaultsCmd)
}

var StaticCmd = &cobra.Command{
	Use:   "static [service]",
	Short: "start a static server",
	RunE: func(cmd *cobra.Command, args []string) error {
		level := hclog.LevelFromString(viper.GetString("log-level"))
		if level == hclog.NoLevel {
			return fmt.Errorf("invalid log level: %s", viper.GetString("log-level"))
		}

		service := viper.GetString(vars.SerivceKey)
		if service == "" {
			return fmt.Errorf("no service specified")
		}

		cfg := &runner.Config{
			Address: viper.GetString("host") + ":" + viper.GetString("port"),
			Logging: &hclog.LoggerOptions{
				Name:  "grpc-hole",
				Level: level,
				Color: hclog.AutoColor,
			},
			ServiceName: service,
			ServerType:  server.StaticServer,
		}

		r, err := runner.New(cfg)
		if err != nil {
			color.Red(fmt.Errorf("failed to create server: %w", err).Error())
			os.Exit(1)
		}

		if err = r.Run(); err != nil {
			color.Red(fmt.Errorf("failed to run server: %w", err).Error())
			os.Exit(1)
		}

		return nil
	},
}
