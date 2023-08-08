package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func v1() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// s := grpc.NewServer()
	s := grpc.NewServer(grpc.UnaryInterceptor(testInterceptor))

	// pb.RegisterMyAppServiceServer(s, &server{})

	fmt.Println("Server started")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return

	conn, err := grpc.Dial(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	refClient := grpcreflect.NewClientAuto(context.Background(), conn)
	defer refClient.Reset()

	serviceNames, err := refClient.ListServices()
	if err != nil {
		panic(err)
	}

	service, _ := refClient.ResolveService(serviceNames[1])
	servFile := service.GetFile()

	fileNames := SaveDeps(servFile, "test1")

	descriptorSet := &descriptorpb.FileDescriptorSet{}
	for _, fileName := range fileNames {
		descFile := &descriptorpb.FileDescriptorProto{}
		utils.ProtoLoadAndUnmarshal(fileName, descFile)
		descriptorSet.File = append(descriptorSet.File, descFile)
	}

	prf, err := protodesc.NewFiles(descriptorSet)
	if err != nil {
		panic(err)
	}

	desc, err := prf.FindDescriptorByName(protoreflect.FullName("proto.qoutes.Qoute"))
	if err != nil {
		panic(err)
	}

	br, err := os.ReadFile("testmsg.json")
	if err != nil {
		panic(err)
	}
	msg := dynamicpb.NewMessage(desc.(protoreflect.MessageDescriptor))
	protojson.Unmarshal(br, msg)

	fmt.Println("MESSAGE:", msg)
}

// Saves all of the files dependencies
func SaveDeps(file *desc.FileDescriptor, serviceName string) []string {
	fileNameCache := make(map[string]struct{})
	saveDepRec(file, fileNameCache)

	fileNames := make([]string, 0)
	for fileName := range fileNameCache {
		fileNames = append(fileNames, fileName)
	}

	return fileNames
}

func saveDepRec(file *desc.FileDescriptor, fileNameCache map[string]struct{}) {
	// File already saved
	if _, ok := fileNameCache[file.GetName()]; ok {
		return
	}
	fileNameCache[file.GetName()] = struct{}{}

	// Save the file
	utils.ProtoJsonMarshalAndSave(file.AsFileDescriptorProto(), file.GetName())

	for _, dep := range file.GetDependencies() {
		saveDepRec(dep, fileNameCache)
	}
}
