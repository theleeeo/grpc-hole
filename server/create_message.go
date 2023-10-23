package server

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/protobuf/types/descriptorpb"
)

func CreatePopulatedMessage(f *desc.MessageDescriptor) *dynamic.Message {
	msg := dynamic.NewMessage(f)
	for _, field := range f.GetFields() {
		var value any
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			value = CreatePopulatedMessage(field.GetMessageType())
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			value = "Hello World"
		case descriptorpb.FieldDescriptorProto_TYPE_INT32,
			descriptorpb.FieldDescriptorProto_TYPE_UINT32,
			descriptorpb.FieldDescriptorProto_TYPE_SINT32,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			value = int32(123)
		case descriptorpb.FieldDescriptorProto_TYPE_INT64,
			descriptorpb.FieldDescriptorProto_TYPE_UINT64,
			descriptorpb.FieldDescriptorProto_TYPE_SINT64,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			value = int64(123)
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			value = true
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			value = float32(123.123)
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			value = float64(123.123)
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			value = []byte("Hello World")
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			value = int32(1)
		default:
			panic(fmt.Errorf("unhandled type: %s", field.GetType()))
		}
		if field.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
			msg.AddRepeatedField(field, value)
		} else {
			msg.SetField(field, value)
		}
	}
	return msg
}
