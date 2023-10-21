package scanning

import (
	"context"
	"errors"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrNoAddress = errors.New("no address provided")
)

func ScanServer(addr string) ([]*desc.ServiceDescriptor, error) {
	if addr == "" {
		return nil, ErrNoAddress
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	refClient := grpcreflect.NewClientAuto(context.Background(), conn)
	defer refClient.Reset()

	serviceNames, err := refClient.ListServices()
	if err != nil {
		return nil, err
	}

	mainServiceNames := getNonReflectServices(serviceNames)
	services := make([]*desc.ServiceDescriptor, 0)
	for _, mainService := range mainServiceNames {
		serviceDescr, _ := refClient.ResolveService(mainService)
		services = append(services, serviceDescr)
	}

	return services, nil
}

func getFirstNonReflectionService(services []string) string {
	for _, service := range services {
		if !strings.Contains(service, "grpc.reflection") {
			return service
		}
	}
	return ""
}

func getNonReflectServices(services []string) []string {
	nonReflectServices := make([]string, 0)
	for _, service := range services {
		if !strings.Contains(service, "grpc.reflection") {
			nonReflectServices = append(nonReflectServices, service)
		}
	}
	return nonReflectServices
}
