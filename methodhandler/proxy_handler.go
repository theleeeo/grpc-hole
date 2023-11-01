package methodhandler

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type proxyHandler struct {
	method *desc.MethodDescriptor
	lg     hclog.Logger

	targetChannel *grpc.ClientConn
}

func NewProxyHandler(method *desc.MethodDescriptor, logger hclog.Logger, proxyTarget string) (Handler, error) {
	if method.IsClientStreaming() || method.IsServerStreaming() {
		return nil, fmt.Errorf("only unary methods is supported; %q is %s", method.GetFullyQualifiedName(), methodType(method))
	}

	conn, err := createServiceConn(context.Background(), proxyTarget)
	if err != nil {
		logger.Error("Failed to create stub", "Error", err)
		return nil, err
	}

	return &proxyHandler{
		method:        method,
		lg:            logger,
		targetChannel: conn,
	}, nil
}

func (h *proxyHandler) Name() string {
	return h.method.GetName()
}

func (h *proxyHandler) Handle(stream grpc.ServerStream) error {
	inputMsg := dynamic.NewMessage(h.method.GetInputType())
	if err := stream.RecvMsg(inputMsg); err != nil {
		return err
	}

	h.lg.Info("Received request", "Method", h.method.GetName())
	if h.lg.IsDebug() {
		h.lg.Debug("Request message")
		fmt.Println(inputMsg.String())
	}

	resp := dynamic.NewMessage(h.method.GetOutputType())
	if err := h.targetChannel.Invoke(context.Background(), requestMethod(h.method), inputMsg, resp); err != nil {
		return err
	}

	if h.lg.IsDebug() {
		h.lg.Debug("Response message")
		fmt.Println(resp.String())
	}

	return stream.SendMsg(resp)
}

func requestMethod(md *desc.MethodDescriptor) string {
	return fmt.Sprintf("/%s/%s", md.GetService().GetFullyQualifiedName(), md.GetName())
}

func methodType(md *desc.MethodDescriptor) string {
	if md.IsClientStreaming() && md.IsServerStreaming() {
		return "bidi-streaming"
	} else if md.IsClientStreaming() {
		return "client-streaming"
	} else if md.IsServerStreaming() {
		return "server-streaming"
	}
	return "unary"

}

func createServiceConn(ctx context.Context, target string) (*grpc.ClientConn, error) {
	// Create a gRPC client connection to the target server
	conn, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
