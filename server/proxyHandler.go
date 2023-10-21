package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	encoding.RegisterCodec(proxyCodec{})
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
	out := createPopulatedMessage(outType)

	return stream.SendMsg(out)
}

func createPopulatedMessage(f *desc.MessageDescriptor) *dynamic.Message {
	msg := dynamic.NewMessage(f)
	for _, field := range f.GetFields() {
		var value any
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			value = createPopulatedMessage(field.GetMessageType())
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			value = "Hello World"
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
			value = 123
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			value = true
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
			descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			value = 123.123
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			value = []byte("Hello World")
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			value = 1
		}
		if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			msg.AddRepeatedField(field, value)
		} else {
			msg.SetField(field, value)
		}
	}
	return msg
}
