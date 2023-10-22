package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
)

func (s *Server) Handler(srv any, stream grpc.ServerStream) error {
	// Extract the fullMethodName name from the stream's context.
	fullMethodName, ok := grpc.Method(stream.Context())
	if !ok {
		s.lg.Error("Failed to get method from context")
		return fmt.Errorf("unknown method")
	}

	method := s.service.FindMethodByName(getMethodName(fullMethodName))
	if method == nil {
		s.lg.Error("Failed to find method", "Method", fullMethodName)
		return fmt.Errorf("unknown method")
	}

	err := s.handleRequest(stream, method)
	if err != nil {
		s.lg.Error("Failed to handle request", "Method", method.GetName(), "Error", err)
		return err
	}

	// return stream.SendMsg(&dyn) // You may want to customize the response.
	return nil
}

func getMethodName(fullName string) string {
	nameParts := strings.Split(fullName, "/")
	return nameParts[len(nameParts)-1]
}

func (s *Server) handleRequest(stream grpc.ServerStream, method *desc.MethodDescriptor) error {
	msg := dynamic.NewMessage(method.GetInputType())
	if err := stream.RecvMsg(msg); err != nil {
		return err
	}

	m, _ := msg.MarshalJSON()
	s.lg.Info("Received request", "Method", method.GetName(), "Input", string(m))

	outType := method.GetOutputType()

	out, err := service.LoadResponse(method.GetService().GetName(), method.GetName(), outType)
	if err != nil {
		// If the error is something else than "file not found", return it.
		if !os.IsNotExist(err) {
			return err
		}
		out = CreatePopulatedMessage(outType)
	}

	return stream.SendMsg(out)
}
