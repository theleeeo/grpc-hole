package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
)

const (
	dependencyDirName = "deps"
)

// Persist a service to disk
func Save(dir string, serv *desc.ServiceDescriptor) error {
	data := &serviceData{
		Name:    serv.GetFullyQualifiedName(),
		File:    serv.GetFile().GetName(),
		SavedAt: time.Now(),
	}

	servicePath := filepath.Join(dir, serv.GetName())

	if err := deleteDir(servicePath); err != nil {
		return err
	}

	depPath := filepath.Join(servicePath, dependencyDirName)
	// save the service file to disk
	if err := saveFileDesc(serv.GetFile(), depPath); err != nil {
		return err
	}

	// save all dependencies to disk
	fileNameCache := make(map[string]struct{})
	if err := saveDepRec(serv.GetFile(), fileNameCache, depPath); err != nil {
		return err
	}

	for fileName := range fileNameCache {
		data.DependentFiles = append(data.DependentFiles, fileName)
	}

	// save the service data to a file
	if err := saveDataFile(data, servicePath); err != nil {
		return err
	}

	return nil
}

// Recursively go through all dependencies and save it to disk
func saveDepRec(file *desc.FileDescriptor, fileNameCache map[string]struct{}, pathPrefix string) error {
	for _, dep := range file.GetDependencies() {
		// File already saved
		if _, ok := fileNameCache[dep.GetName()]; ok {
			continue
		}
		fileNameCache[dep.GetName()] = struct{}{}

		if err := saveFileDesc(dep, pathPrefix); err != nil {
			return err
		}

		if err := saveDepRec(dep, fileNameCache, pathPrefix); err != nil {
			return err
		}
	}

	return nil
}

func deleteDir(path string) error {
	return os.RemoveAll(path)
}

func saveFileDesc(file *desc.FileDescriptor, pathPrefix string) error {
	return utils.ProtoJSONMarshalAndSave(file.AsFileDescriptorProto(), filepath.Join(pathPrefix, file.GetName()))
}

func saveDataFile(data *serviceData, path string) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, ServiceDataFileName), dataJSON, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
