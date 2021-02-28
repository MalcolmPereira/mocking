//Package server for mocking server definition.
package server

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	logger "github.com/sirupsen/logrus"
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
	MockFolders []string `yaml:"mocksfolders"`
}

//ProcessMockingYAML parses and validates mocking api yaml file
func ProcessMockingYAML(yamlFilePath string) (*MockingYAML, error) {
	mockingYAML := MockingYAML{}
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Error("Error reading mocking YAML file ", err)
		return nil, errors.New("Error reading mocking YAML file")
	}
	err = yaml.Unmarshal(yamlFile, &mockingYAML)
	if err != nil {
		logger.Error("Error processing mocking YAML file ", err)
		return nil, errors.New("Error processing mocking YAML file")
	}

	err = validateYAML(&mockingYAML)
	if err != nil {
		logger.Error("Error validating mocking YAML file ", err)
		return nil, errors.New("Error validating mocking YAML file")
	}

	return &mockingYAML, nil
}

//Validations for Mocking YAML
func validateYAML(mockingYAML *MockingYAML) error {
	if strings.TrimSpace(mockingYAML.Version) != currentVersion {
		logger.Error("Invalid Mocking API Version, supported version is:", currentVersion)
		return errors.New("Invalid Mocking API Version, supported version is: " + currentVersion)
	}
	if len(strings.TrimSpace(mockingYAML.Name)) == 0 {
		logger.Error("Invalid Mocking Server Name, valid Mocking Server Name required")
		return errors.New("Invalid Mocking Server Name, valid Mocking Server Name required")
	}
	portNum, err := strconv.Atoi(mockingYAML.Port)
	if err != nil {
		logger.Error("Error processing mocking YAML file, Invalid Number Value for Port", err)
		return errors.New("Error processing mocking YAML file, Invalid Number Value for Port ")
	}
	if portNum < 1080 || portNum > 65535 {
		logger.Error("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535 ")
		return errors.New("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535")
	}

	if len(strings.TrimSpace(mockingYAML.Loglevel)) == 0 {
		mockingYAML.Loglevel = "INFO"

	} else if !strings.EqualFold(mockingYAML.Loglevel, "INFO") &&
		!strings.EqualFold(mockingYAML.Loglevel, "DEBUG") &&
		!strings.EqualFold(mockingYAML.Loglevel, "WARN") &&
		!strings.EqualFold(mockingYAML.Loglevel, "ERROR") &&
		!strings.EqualFold(mockingYAML.Loglevel, "TRACE") {
		mockingYAML.Loglevel = "INFO"
	}

	if len(strings.TrimSpace(mockingYAML.HTTPS.CertFile)) > 0 {
		_, err := os.Stat(mockingYAML.HTTPS.CertFile)
		if err != nil {
			logger.Error("Invalid Mocking Server HTTPS Certificatie File, valid HTTPS Certificatie File required ")
			return errors.New("Invalid Mocking Server HTTPS Certificatie File, valid HTTPS Certificatie File required")
		}
	}

	if len(strings.TrimSpace(mockingYAML.HTTPS.CertKeyFile)) > 0 {
		_, err := os.Stat(mockingYAML.HTTPS.CertKeyFile)
		if err != nil {
			logger.Error("Invalid Mocking Server HTTPS Certificatie Key File, valid HTTPS Certificatie Key File required ")
			return errors.New("Invalid Mocking Server HTTPS Certificatie Key File, valid HTTPS Certificatie Key File required")
		}
	}

	if len(mockingYAML.MockFolders) == 0 {
		logger.Error("Invalid Mocking Server Mock Folders, valid Mock Folders path required ")
		return errors.New("Invalid Mocking Server Mock Folders, valid Mock Folders path required")
	}

	for _, mockfolder := range mockingYAML.MockFolders {
		_, err := os.Stat(mockfolder)
		if err != nil {
			logger.Error("Invalid Mocking Server Mock Folder Path, valid Mock Folder Path required ")
			return errors.New("Invalid Mocking Server Mock Folder Path, valid Mock Folder Path required")
		}
	}

	return nil
}
