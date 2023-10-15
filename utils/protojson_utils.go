package utils

import (
	"os"
	"path"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func ProtoJsonMarshalAndSave(m protoreflect.ProtoMessage, fileName string) error {
	bytes, err := protojson.Marshal(m)
	if err != nil {
		return err
	}

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
		return err
	}

	return nil
}

// Replaces the fileextesion with json
func setExtToJSON(fileName string) string {
	return strings.TrimSuffix(fileName, path.Ext(fileName)) + ".json"
}
