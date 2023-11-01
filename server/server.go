package server

import (
	"fmt"

	"github.com/TheLeeeo/grpc-hole/methodhandler"
	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
)

type Server struct {
	service *desc.ServiceDescriptor

	methods map[string]methodhandler.Handler

	lg hclog.Logger
}

func (s *Server) Name() string {
	return s.service.GetName()
}

func New(cfg *Config) (*Server, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	methodHandlers, err := createMethodHandlers(cfg)
	if err != nil {
		return nil, err
	}

	return &Server{
		service: cfg.Service,
		methods: methodHandlers,
		lg:      cfg.Logger,
	}, nil
}

func createMethodHandlers(cfg *Config) (map[string]methodhandler.Handler, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	methodHandlers := make(map[string]methodhandler.Handler)
	for _, method := range cfg.Service.GetMethods() {

		h, err := createMethod(cfg, method)
		if err != nil {
			return nil, err
		}

		cfg.Logger.Debug("Registering method", "Method", method.GetName())
		methodHandlers[method.GetName()] = h
	}

	return methodHandlers, nil
}

func createMethod(cfg *Config, method *desc.MethodDescriptor) (methodhandler.Handler, error) {
	switch cfg.ServerType {
	case StaticServer:
		h := methodhandler.NewStaticHandler(method, cfg.Logger)
		return h, nil
	case ProxyServer:
		h, err := methodhandler.NewProxyHandler(method, cfg.Logger, cfg.ProxyAddress)
		if err != nil {
			return nil, err
		}
		return h, nil
	}

	return nil, fmt.Errorf("invalid server type specified: %q", cfg.ServerType)
}
