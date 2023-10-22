package server

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
)

type Config struct {
	// The service descriptor to use
	Service *desc.ServiceDescriptor

	// The directory that contains the service
	ServiceDir string

	// The logger to use
	Logger hclog.Logger
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.Service == nil {
		return fmt.Errorf("no service specified")
	}

	if c.ServiceDir == "" {
		return fmt.Errorf("no service directory specified")
	}

	if c.Logger == nil {
		c.Logger = hclog.Default()
		c.Logger.Info("no logger was provided, using default")
	}

	return nil
}
