package service

import (
	"os"
	"strings"

	"github.com/TheLeeeo/grpc-hole/utils"
	"github.com/jhump/protoreflect/desc"
)

// Persist a service to disk
func (s *Service) Save() {
	fileNameCache := make(map[string]struct{})

	path := "saved_services" + "/" + extractServiceName(s.Name) + "/"

	saveDepRec(s.GetDescriptorFile(), fileNameCache, path)

	fileNames := make([]byte, 0)
	for fileName := range fileNameCache {
		bytes := []byte(fileName + "\n")
		fileNames = append(fileNames, bytes...)
	}

	// save the filenames to a file
	err := os.WriteFile(path+"/data", fileNames, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

// Recursivly go through all dependencies and save it to disk
func saveDepRec(file *desc.FileDescriptor, fileNameCache map[string]struct{}, pathPrefix string) {
	// File already saved
	if _, ok := fileNameCache[file.GetName()]; ok {
		return
	}
	fileNameCache[file.GetName()] = struct{}{}

	// Save the file
	utils.ProtoJsonMarshalAndSave(file.AsFileDescriptorProto(), pathPrefix+file.GetName())

	for _, dep := range file.GetDependencies() {
		saveDepRec(dep, fileNameCache, pathPrefix)
	}
}

func extractServiceName(fullName string) string {
	nameParts := strings.Split(fullName, ".")
	return nameParts[len(nameParts)-1]
}
