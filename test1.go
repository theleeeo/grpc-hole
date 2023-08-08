package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
)

func testInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	fmt.Println("testInterceptor")
	fmt.Println(req)
	fmt.Println(info.FullMethod)
	fmt.Println(info.Server)
	return
	// return handler(ctx, req)
}

func Test() {}

type TestServer struct{}

func _TestService_Test_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(proto.Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	// if interceptor == nil {
	// 	return srv.(QouteServiceServer).GetDailyQoute(ctx, in)
	// }
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "TestTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return Test, nil
	}
	return interceptor(ctx, in, info, handler)
}

var TestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.qoutes.QouteService",
	HandlerType: (new(any)),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDailyQoute",
			Handler:    _TestService_Test_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "qoutes/qoute.proto",
}

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("Method:", info.FullMethod)
	return handler(ctx, req)
}

type proxyCodec struct{}

func (proxyCodec) Marshal(v interface{}) ([]byte, error) {
	return *(v.(*[]byte)), nil
}

func (proxyCodec) Unmarshal(data []byte, v interface{}) error {
	vv := v.(*[]byte)
	*vv = data
	return nil
}

func (proxyCodec) Name() string {
	return "proxy"
}

func main() {
	// scanning.ScanService(":5001")
	s, err := service.Load("QouteService")
	if err != nil {
		panic(err)
	}

	meth := s.GetMethod("GetRandomQoute")
	fmt.Println(meth)

	srv := grpc.NewServer(grpc.UnknownServiceHandler(createProxyHandler(meth.GetInputType(), meth.GetOutputType())))

	// ts := &TestServer{}
	// s.RegisterService(&TestService_ServiceDesc, ts)

	// reflection.Register(s)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Server started")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	encoding.RegisterCodec(proxyCodec{})
}

func createProxyHandler(in, out *desc.MessageDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		// Extract the method name from the stream's context.
		method, ok := grpc.Method(stream.Context())
		if !ok {
			log.Println("Failed to get method from context")
			return fmt.Errorf("unknown method")
		}

		fmt.Printf("Received request for method: %s\n", method)

		msg := dynamic.NewMessage(in)
		if err := stream.RecvMsg(msg); err != nil {
			return err
		}

		fmt.Println(msg)
		fs := msg.GetKnownFields()
		fmt.Println(fs)
		fmt.Println(msg.GetField(fs[0]))

		// return stream.SendMsg(&dyn) // You may want to customize the response.
		return nil
	}
}
