package server

import (
	"encoding/json"
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
	inputMsg := dynamic.NewMessage(method.GetInputType())
	if err := stream.RecvMsg(inputMsg); err != nil {
		return err
	}

	inputJSON, _ := inputMsg.MarshalJSON()
	s.lg.Info("Received request", "Method", method.GetName(), "Input", string(inputJSON))

	outType := method.GetOutputType()
	var out *dynamic.Message

	respTemplate, err := service.LoadResponse(method.GetService().GetName(), method.GetName())
	if err != nil {
		// If the error is something else than "file not found", return it.
		if !os.IsNotExist(err) {
			return err
		}
		out = CreatePopulatedMessage(outType)
	}

	if respTemplate != nil {
		var inputMap map[string]any
		if err := json.Unmarshal(inputJSON, &inputMap); err != nil {
			return err
		}

		outJson, err := service.ParseTemplate(inputMap, respTemplate)
		if err != nil {
			s.lg.Error("Failed to parse template", "Method", method.GetName(), "Error", err)
			return err
		}

		out = dynamic.NewMessage(outType)
		if err := out.UnmarshalJSON(outJson); err != nil {
			s.lg.Error("Failed to unmarshal json", "Method", method.GetName(), "Error", err)
			return err
		}
	}

	return stream.SendMsg(out)
}
