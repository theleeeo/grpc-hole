package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
)

const (
	directory = "saved_services"
)

// Persist a service to disk
func Save(serv *desc.ServiceDescriptor) error {
	data := &serviceData{
		Name: serv.GetFullyQualifiedName(),
		File: serv.GetFile().GetName(),
	}

	if err := deleteDir(directory); err != nil {
		return err
	}

	path := filepath.Join(directory, serv.GetName())

	// save the service file to disk
	if err := saveFileDesc(serv.GetFile(), path); err != nil {
		return err
	}

	// save all dependencies to disk
	fileNameCache := make(map[string]struct{})
	saveDepRec(serv.GetFile(), fileNameCache, path)

	for fileName := range fileNameCache {
		data.DependentFiles = append(data.DependentFiles, fileName)
	}

	// save the service data to a file
	if err := saveDataFile(data, path); err != nil {
		return err
	}

	return nil
}

// Recursivly go through all dependencies and save it to disk
func saveDepRec(file *desc.FileDescriptor, fileNameCache map[string]struct{}, pathPrefix string) error {
	for _, dep := range file.GetDependencies() {
		// File already saved
		if _, ok := fileNameCache[dep.GetName()]; ok {
			break
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

func extractServiceName(fullName string) string {
	nameParts := strings.Split(fullName, ".")
	return nameParts[len(nameParts)-1]
}

func deleteDir(path string) error {
	return os.RemoveAll(path)
}

func saveFileDesc(file *desc.FileDescriptor, pathPrefix string) error {
	return utils.ProtoJsonMarshalAndSave(file.AsFileDescriptorProto(), filepath.Join(pathPrefix, file.GetName()))
}

func saveDataFile(data *serviceData, path string) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(path, "/data.json"), dataJson, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
