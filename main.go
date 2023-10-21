package main

import (
	"os"

	"github.com/TheLeeeo/grpc-hole/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
