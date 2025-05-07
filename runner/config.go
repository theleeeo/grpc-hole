package runner

import (
	"fmt"

	"github.com/TheLeeeo/grpc-hole/server"
	"github.com/hashicorp/go-hclog"
)

type Config struct {
	// The address to listen on
	Address string

	// The name of the service to run
	ServiceNames []string

	// The configuration for the logger
	Logging *hclog.LoggerOptions

	ServerType server.ServerType

	// ProxyAddress is the address to proxy requests to
	ProxyAddress string
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.Address == "" {
		return fmt.Errorf("no address specified")
	}

	if len(c.ServiceNames) == 0 {
		return fmt.Errorf("no service specified")
	}

	return nil
}
