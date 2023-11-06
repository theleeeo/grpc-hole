package main

import (
	"os"

	"github.com/TheLeeeo/grpc-hole/cli"
	"github.com/fatih/color"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
