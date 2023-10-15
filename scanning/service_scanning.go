package scanning

import (
	"context"
	"strings"

	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ScanService(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	refClient := grpcreflect.NewClientAuto(context.Background(), conn)
	defer refClient.Reset()

	serviceNames, err := refClient.ListServices()
	if err != nil {
		return err
	}

	mainServices := getNonReflectServices(serviceNames)

	for _, mainService := range mainServices {
		serviceDescr, _ := refClient.ResolveService(mainService)

		service.Save(serviceDescr)
	}

	return nil
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
