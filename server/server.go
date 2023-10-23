package server

import (
	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
)

type Server struct {
	service *desc.ServiceDescriptor

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

	return &Server{
		service: cfg.Service,
		lg:      cfg.Logger,
	}, nil
}
