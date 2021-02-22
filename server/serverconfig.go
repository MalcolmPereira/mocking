//Package server for mocking server definition.
package server

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

//Current Supported Version
const currentVersion string = "0.1"

//MockingYAML struct wraps mocking api configuration
type MockingYAML struct {
	Version  string `yaml:"version"`
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	Loglevel string `yaml:"loglevel"`
	HTTPS    struct {
		CertFile    string `yaml:"certfile"`
		CertKeyFile string `yaml:"certkeyfile"`
	} `yaml:"https"`
	Mocks string `yaml:"mocksfolder"`
}

//ProcessMockingYAML parses and validates mocking api yaml file
func ProcessMockingYAML(yamlFilePath string) (*MockingYAML, error) {
	mockingYAML := MockingYAML{}
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		log.Printf("Error reading mocking YAML file: %v", err)
		return nil, errors.New("Error reading mocking YAML file")
	}
	err = yaml.Unmarshal(yamlFile, &mockingYAML)
	if err != nil {
		log.Printf("Error processing mocking YAML file: %v", err)
		return nil, errors.New("Error processing mocking YAML file")
	}

	err = validateYAML(&mockingYAML)
	if err != nil {
		log.Printf("Error validating mocking YAML file: %v", err)
		return nil, errors.New("Error validating mocking YAML file")
	}

	return &mockingYAML, nil
}

//Validations for Mocking YAML
func validateYAML(mockingYAML *MockingYAML) error {
	if strings.Trim(mockingYAML.Version, " ") != currentVersion {
		log.Println("Invalid Mocking API Version, supported version is:", currentVersion)
		return errors.New("Invalid Mocking API Version, supported version is: " + currentVersion)
	}
	if len(strings.Trim(mockingYAML.Name, " ")) == 0 {
		log.Println("Invalid Mocking Server Name, valid Mocking Server Name required")
		return errors.New("Invalid Mocking Server Name, valid Mocking Server Name required")
	}
	portNum, err := strconv.Atoi(mockingYAML.Port)
	if err != nil {
		log.Println("Error processing mocking YAML file, Invalid Number Value for Port", err)
		return errors.New("Error processing mocking YAML file, Invalid Number Value for Port ")
	}
	if portNum < 1080 || portNum > 65535 {
		log.Println("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535 ")
		return errors.New("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535")
	}
	return nil
}
