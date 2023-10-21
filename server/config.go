package server

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
)

type Config struct {
	// The address to listen on
	Address string

	// The service descriptor to use
	Service *desc.ServiceDescriptor
}

func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}

	if c.Address == "" {
		return fmt.Errorf("no address specified")
	}

	if c.Service == nil {
		return fmt.Errorf("no service specified")
	}

	return nil
}
