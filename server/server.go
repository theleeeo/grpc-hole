package server

import (
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

	methodHandlers := make(map[string]methodhandler.Handler)
	for _, method := range cfg.Service.GetMethods() {
		cfg.Logger.Debug("Registering method", "Method", method.GetName())
		methodHandlers[method.GetName()] = methodhandler.NewDynamicHandler(method, cfg.Logger)
	}

	return &Server{
		service: cfg.Service,
		methods: methodHandlers,
		lg:      cfg.Logger,
	}, nil
}
