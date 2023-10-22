package runner

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
)

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.Address == "" {
		return fmt.Errorf("no address specified")
	}

	if c.ServiceName == "" {
		return fmt.Errorf("no service specified")
	}

	return nil
}

type Config struct {
	// The address to listen on
	Address string

	// The name of the service to run
	ServiceName string

	// The configuration for the logger
	Logging *hclog.LoggerOptions
}
