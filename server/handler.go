package server

import (
	"fmt"
	"strings"

	"google.golang.org/grpc"
)

func (s *Server) Handler(_ any, stream grpc.ServerStream) error {
	// Extract the fullMethodName name from the stream's context.
	fullMethodName, ok := grpc.Method(stream.Context())
	if !ok {
		s.lg.Error("Failed to get method from context")
		return fmt.Errorf("unknown method")
	}

	method, ok := s.methods[getMethodName(fullMethodName)]
	if !ok {
		s.lg.Error("Failed to find method", "Method", fullMethodName)
		return fmt.Errorf("unknown method")
	}

	err := method.Handle(stream)
	if err != nil {
		s.lg.Error("Failed to handle request", "Method", method.Name(), "Error", err)
		return err
	}

	return nil
}

func getMethodName(fullName string) string {
	nameParts := strings.Split(fullName, "/")
	return nameParts[len(nameParts)-1]
}
