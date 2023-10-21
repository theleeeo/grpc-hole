package server

import (
	"fmt"
	"net"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

type server struct {
	cfg *Config

	lg hclog.Logger
}

func New(cfg *Config, logger hclog.Logger) (*server, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	if logger == nil {
		logger = hclog.Default()
		logger.Info("no logger was provided, using default")
	}

	return &server{
		cfg: cfg,
		lg:  logger,
	}, nil
}

func (s *server) Run() error {
	srv := grpc.NewServer(grpc.UnknownServiceHandler(createProxyHandler(s.cfg.Service)))

	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.lg.Info("Server started", "address", s.cfg.Address, "service", s.cfg.Service.GetName())
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
