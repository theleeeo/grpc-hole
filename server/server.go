package server

import (
	"github.com/TheLeeeo/grpc-hole/server/methodhandler"
	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
)

type Server struct {
	service *desc.ServiceDescriptor

	methods map[string]methodhandler.MethodHandler

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

	methodHandlers := make(map[string]methodhandler.MethodHandler)
	for _, method := range cfg.Service.GetMethods() {
		cfg.Logger.Debug("Registering method", "Method", method.GetName())
		methodHandlers[method.GetName()] = methodhandler.New(method, cfg.Logger)
	}

	return &Server{
		service: cfg.Service,
		methods: methodHandlers,
		lg:      cfg.Logger,
	}, nil
}
