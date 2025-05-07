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
	ProxyCmd.Flags().String("proxy", "", "the address to proxy to")
	if err := viper.BindPFlag("proxy", ProxyCmd.Flags().Lookup("proxy")); err != nil {
		panic(err)
	}
}

var ProxyCmd = &cobra.Command{
	Use:   "proxy [service]",
	Short: "start a proxy server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// When an error with the command syntax occurs, return the error.
		// This will print the error and the usage information.
		// If it is a runtime error, print the error and exit with code 1.
		level := hclog.LevelFromString(viper.GetString("log-level"))
		if level == hclog.NoLevel {
			return fmt.Errorf("invalid log level: %s", viper.GetString("log-level"))
		}

		services := viper.GetStringSlice(vars.SerivceKey)
		if len(services) == 0 {
			return fmt.Errorf("no service specified")
		}

		cfg := &runner.Config{
			Address: viper.GetString("host") + ":" + viper.GetString("port"),
			Logging: &hclog.LoggerOptions{
				Name:  "grpc-hole",
				Level: level,
				Color: hclog.AutoColor,
			},
			ServiceNames: services,
			ServerType:   server.ProxyServer,
			ProxyAddress: viper.GetString("proxy"),
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
