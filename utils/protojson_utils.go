package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ProtoJSONMarshalAndSave(m protoreflect.ProtoMessage, fileName string) error {
	bytes, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

	fileName = path.Clean(fileName)
	// fileName = filepath.Base()

	err = os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return err
	}

	fileName = setExtToJSON(fileName)

	err = os.WriteFile(fileName, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func ProtoLoadAndUnmarshal(fileName string, m protoreflect.ProtoMessage) error {
	fileName = setExtToJSON(fileName)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = protojson.Unmarshal(file, m)
	if err != nil {
		fmt.Println("das ist ein error", err)
		return err
	}

	return nil
}

// Replaces the fileextesion with json
func setExtToJSON(fileName string) string {
	return strings.TrimSuffix(fileName, path.Ext(fileName)) + ".json"
}
