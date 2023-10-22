package service

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	ErrServiceNotFound = errors.New("Service not found")
)

// Load a service from disk
func Load(serviceDir string) (*desc.ServiceDescriptor, error) {
	data, err := LoadDataFile(serviceDir)
	if err != nil {
		return nil, err
	}

	descrSet, err := loadDescriptorSet(data.File, data.DependentFiles, serviceDir)
	if err != nil {
		return nil, err
	}

	fileDescr, err := desc.CreateFileDescriptorFromSet(descrSet)
	if err != nil {
		return nil, err
	}

	service := fileDescr.FindService(data.Name)
	if service == nil {
		return nil, ErrServiceNotFound
	}

	return service, nil
}

func loadDescriptorSet(mainFile string, deps []string, path string) (*descriptorpb.FileDescriptorSet, error) {
	descriptorSet := &descriptorpb.FileDescriptorSet{}
	depPath := filepath.Join(path, dependencyDirName)
	// Load all dependencies
	for _, fileName := range deps {
		descFile := &descriptorpb.FileDescriptorProto{}
		if err := utils.ProtoLoadAndUnmarshal(filepath.Join(depPath, fileName), descFile); err != nil {
			return nil, err
		}
		descriptorSet.File = append(descriptorSet.File, descFile)
	}

	// Load the service
	descFile := &descriptorpb.FileDescriptorProto{}
	if err := utils.ProtoLoadAndUnmarshal(filepath.Join(depPath, mainFile), descFile); err != nil {
		return nil, err
	}
	descriptorSet.File = append(descriptorSet.File, descFile)

	return descriptorSet, nil
}

func LoadDataFile(servicePath string) (*serviceData, error) {
	data, err := os.ReadFile(filepath.Join(servicePath, ServiceDataFileName))
	if err != nil {
		return nil, err
	}

	sd := &serviceData{}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return nil, err
	}

	return sd, nil
}
