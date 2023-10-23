package service

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/TheLeeeo/grpc-hole/cli/vars"
	"github.com/TheLeeeo/grpc-hole/templateparse"
	"github.com/spf13/viper"
)

const (
	responsesDirName = "responses"
)

func SaveResponseFile(serviceDir string, method string, data []byte) error {
	responseDir := filepath.Join(serviceDir, responsesDirName)

	if err := createDir(responseDir); err != nil {
		return err
	}

	path := filepath.Join(responseDir, method+".json")

	if err := os.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}

	return nil
}

func LoadResponse(serviceName, methodName string) ([]byte, error) {
	baseDir := viper.GetString(vars.SerivceDirKey)
	responseDir := filepath.Join(baseDir, serviceName, responsesDirName)
	path := filepath.Join(responseDir, methodName+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ParseTemplate(input map[string]any, template []byte) ([]byte, error) {
	var templateMap map[string]any
	if err := json.Unmarshal(template, &templateMap); err != nil {
		return nil, err
	}

	out, err := templateparse.ParseTemplate(input, templateMap)
	if err != nil {
		return nil, err
	}

	return json.Marshal(out)
}
