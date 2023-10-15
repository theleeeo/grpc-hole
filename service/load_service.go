package service

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Load a service from disk
func Load(serviceName string) (*desc.ServiceDescriptor, error) {
	path := "saved_services" + "/" + serviceName + "/"

	data, err := os.ReadFile(path + "data.json")
	if err != nil {
		return nil, err
	}

	sd := &serviceData{}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return nil, err
	}

	descrSet, err := loadDescriptorSet(sd.File, sd.DependentFiles, path)
	if err != nil {
		return nil, err
	}

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

	service := fileDescr.FindService(sd.Name)
	if service == nil {
		return nil, errors.New("Service not found")
	}

	return service, nil
}

func loadDescriptorSet(mainFile string, deps []string, path string) (*descriptorpb.FileDescriptorSet, error) {
	descriptorSet := &descriptorpb.FileDescriptorSet{}

	// Load all dependencies
	for _, fileName := range deps {
		descFile := &descriptorpb.FileDescriptorProto{}
		if err := utils.ProtoLoadAndUnmarshal(path+fileName, descFile); err != nil {
			return nil, err
		}
		descriptorSet.File = append(descriptorSet.File, descFile)
	}

	// Load the service
	descFile := &descriptorpb.FileDescriptorProto{}
	if err := utils.ProtoLoadAndUnmarshal(path+mainFile, descFile); err != nil {
		return nil, err
	}
	descriptorSet.File = append(descriptorSet.File, descFile)

	return descriptorSet, nil
}
