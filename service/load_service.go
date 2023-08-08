package service

import (
	"os"
	"strings"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Load a service from disk
func Load(serviceName string) (*Service, error) {
	path := "saved_services" + "/" + serviceName + "/"

	data, err := os.ReadFile(path + "data")
	if err != nil {
		return nil, err
	}

	fileContent := string(data)
	// Split the fileContent into lines
	lines := strings.Split(fileContent, "\n")
	// Remove last empty line
	lines = lines[:len(lines)-1]

	descrSet := loadDescriptorSet(lines, path)
	fileDescr, err := desc.CreateFileDescriptorFromSet(descrSet)
	if err != nil {
		return nil, err
	}

	// prf, err := protodesc.NewFiles(descriptorSet)
	// if err != nil {
	// 	return nil, err
	// }

	// descr, _ := prf.FindDescriptorByName("proto.qoutes.QouteService")
	// descr = descr.(protoreflect.ServiceDescriptor)

	// fmt.Println(fileDescr.GetServices()[0].GetMethods()) //.GetOutputType())

	return createService(fileDescr)
}

func createService(fileDescr *desc.FileDescriptor) (*Service, error) {
	serviceDescr := fileDescr.GetServices()[0]

	service := New(serviceDescr.GetFullyQualifiedName(), serviceDescr)

	return service, nil
}

func loadDescriptorSet(filenames []string, path string) *descriptorpb.FileDescriptorSet {
	descriptorSet := &descriptorpb.FileDescriptorSet{}

	// Load all dependencies
	for _, fileName := range filenames[1:] {
		descFile := &descriptorpb.FileDescriptorProto{}
		utils.ProtoLoadAndUnmarshal(path+fileName, descFile)
		descriptorSet.File = append(descriptorSet.File, descFile)
	}

	// Load the service
	descFile := &descriptorpb.FileDescriptorProto{}
	utils.ProtoLoadAndUnmarshal(path+filenames[0], descFile)
	descriptorSet.File = append(descriptorSet.File, descFile)

	return descriptorSet
}
