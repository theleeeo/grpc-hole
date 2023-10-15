package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/TheLeeeo/grpc-hole/service"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
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
	// if err := scanning.ScanService(":5001"); err != nil {
	// 	panic(err)
	// }
	// return
	s, err := service.Load("QouteService")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer(grpc.UnknownServiceHandler(createProxyHandler(s)))

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

func createProxyHandler(service *desc.ServiceDescriptor) func(srv interface{}, stream grpc.ServerStream) error {
	return func(srv interface{}, stream grpc.ServerStream) error {
		// Extract the fullMethodName name from the stream's context.
		fullMethodName, ok := grpc.Method(stream.Context())
		if !ok {
			log.Println("Failed to get method from context")
			return fmt.Errorf("unknown method")
		}

		fmt.Printf("Received request for method: %s\n", fullMethodName)

		method := service.FindMethodByName(getMethodName(fullMethodName))

		err := handleRequest(stream, method)
		if err != nil {
			return err
		}

		// return stream.SendMsg(&dyn) // You may want to customize the response.
		return nil
	}
}

func getMethodName(fullName string) string {
	nameParts := strings.Split(fullName, "/")
	return nameParts[len(nameParts)-1]
}

func handleRequest(stream grpc.ServerStream, method *desc.MethodDescriptor) error {
	msg := dynamic.NewMessage(method.GetInputType())
	if err := stream.RecvMsg(msg); err != nil {
		return err
	}

	m, _ := msg.MarshalJSON()
	fmt.Println(string(m))

	outType := method.GetOutputType()
	// out := dynamic.NewMessage(outType)
	out := randomizeMessage(outType)

	// qouteType := outType.GetFields()[0].GetMessageType()
	// qoute := dynamic.NewMessage(qouteType)

	//load the file
	// b, err := os.ReadFile("testmsg.json")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// err = qoute.UnmarshalJSON(b)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// populate the field with random data

	return stream.SendMsg(out)
}

func randomizeMessage(f *desc.MessageDescriptor) *dynamic.Message {
	msg := dynamic.NewMessage(f)
	for _, field := range f.GetFields() {
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, randomizeMessage(field.GetMessageType()))
			} else {
				msg.SetField(field, randomizeMessage(field.GetMessageType()))
			}
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, "Hello World")
			} else {
				msg.SetField(field, "Hello World")
			}
		case descriptorpb.FieldDescriptorProto_TYPE_INT32,
			descriptorpb.FieldDescriptorProto_TYPE_INT64,
			descriptorpb.FieldDescriptorProto_TYPE_UINT32,
			descriptorpb.FieldDescriptorProto_TYPE_UINT64,
			descriptorpb.FieldDescriptorProto_TYPE_SINT32,
			descriptorpb.FieldDescriptorProto_TYPE_SINT64,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, 123)
			} else {
				msg.SetField(field, 123)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, true)
			} else {
				msg.SetField(field, true)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
			descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, 123.123)
			} else {
				msg.SetField(field, 123.123)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, []byte("Hello World"))
			} else {
				msg.SetField(field, []byte("Hello World"))
			}
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
				msg.AddRepeatedField(field, 1)
			} else {
				msg.SetField(field, 1)
			}
		}
	}
	return msg
}
