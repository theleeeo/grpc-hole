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

	// The type of server to start
	ServerType ServerType

	// ProxyAddress is the address to proxy requests to
	ProxyAddress string
}

type ServerType string

const (
	// The server type has not been set, this is an illegal value
	UnsetServer ServerType = ""
	// Proxy the requests to another server
	ProxyServer ServerType = "proxy"
	// Serve the requests using predefined response-templates
	StaticServer ServerType = "static"
)

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

	switch c.ServerType {
	case UnsetServer:
		return fmt.Errorf("no server type specified")
	case ProxyServer:
		if c.ProxyAddress == "" {
			return fmt.Errorf("no proxy address specified")
		}
	case StaticServer:
		// VALID
	default:
		return fmt.Errorf("invalid server type specified: %q", c.ServerType)
	}

	return nil
}
