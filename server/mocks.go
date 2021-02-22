//Package server for mocking server definition.
package server

import (
	"errors"
	"io/ioutil"
	"strings"

	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Mocks struct wraps mocks
type Mocks struct {
	Name     string `yaml:"name"`
	Resource string `yaml:"resource"`
	Mocks    []struct {
		Mock struct {
			Request struct {
				Method  string   `yaml:"method"`
				Headers []string `yaml:"headers"`
				Body    string   `yaml:"body"`
			} `yaml:"request"`
			Response struct {
				Headers []string `yaml:"headers"`
				Status  int      `yaml:"status"`
				Body    string   `yaml:"body"`
			} `yaml:"response"`
		} `yaml:"mock"`
	} `yaml:"mocks"`
}

//ProcessMocks processes mocks available in the folde
func ProcessMocks(mocksFolder string) ([]Mocks, error) {
	var mockList []Mocks

	logger.Info("Start processing mocks folders : ", mocksFolder)

	logger.Debug("Reading mocks folders : ", mocksFolder)
	files, err := ioutil.ReadDir(mocksFolder)
	if err != nil {
		logger.Error("Error reading mocks folders", err)
		return nil, errors.New("Error reading mocks folders")
	}
	logger.Debug("Done Reading mocks folders : ", mocksFolder)

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			mocks := Mocks{}

			logger.Debug("Reading File : ", file.Name())
			fileBytes, err := ioutil.ReadFile(mocksFolder + "/" + file.Name())
			if err != nil {
				logger.Error("Error reading mock file ", err)
				return nil, errors.New("Error reading mock file")
			}
			logger.Debug("Done Reading File : ", file.Name())

			logger.Debug("Unmarshaling YAML File : ", file.Name())
			err = yaml.Unmarshal(fileBytes, &mocks)
			if err != nil {
				logger.Error("Error processing mocking YAML file ", err)
				logger.Error("Error processing mocking YAML file : ", file.Name())
				continue
			}
			logger.Debug("Done Unmarshaling YAML File : ", file.Name())
			mockList = append(mockList, mocks)
		}
	}
	logger.Debug("Got mocks ", mockList)
	return mockList, nil
}
