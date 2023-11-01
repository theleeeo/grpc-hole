package runner

import (
	"fmt"
	"net"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/server"
	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Runner struct {
	// the address to listen on
	addr string

	// the logger used
	lg hclog.Logger

	server *server.Server
}

func New(cfg *Config) (*Runner, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	serviceDir := filepath.Join(viper.GetString(vars.SerivceDirKey), cfg.ServiceName)
	serviceDescr, err := service.Load(serviceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load service: %w", err)
	}

	logger := hclog.New(cfg.Logging)

	serverCfg := &server.Config{
		Service:      serviceDescr,
		ServiceDir:   serviceDir,
		Logger:       logger,
		ServerType:   cfg.ServerType,
		ProxyAddress: cfg.ProxyAddress,
	}

	s, err := server.New(serverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return &Runner{
		addr:   cfg.Address,
		lg:     logger,
		server: s,
	}, nil
}

func (r *Runner) Run() error {
	lis, err := net.Listen("tcp", r.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	srv := grpc.NewServer(grpc.UnknownServiceHandler(r.server.Handler))

	r.lg.Info("Server started", "address", r.addr, "service", r.server.Name())
	if err := srv.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
