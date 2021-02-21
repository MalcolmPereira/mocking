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
	files, err := ioutil.ReadDir(mocksFolder)

	if err != nil {
		logger.Error("Error reading mocks folders", err)
		return nil, errors.New("Error reading mocks folders")
	}

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			logger.Info("Got File Name : ", file.Name())
			mocks := Mocks{}
			logger.Info("Got mocks : ", mocks)

			fileBytes, err := ioutil.ReadFile(mocksFolder + "/" + file.Name())
			if err != nil {
				logger.Error("Error reading mock file ", err)
				return nil, errors.New("Error reading mock file")
			}

			err = yaml.Unmarshal(fileBytes, &mocks)
			if err != nil {
				logger.Error("Error processing mocking YAML file ", err)
				return nil, errors.New("Error processing mocking YAML file")
			}

			mockList = append(mockList, mocks)
		}
	}

	logger.Info("Got mocks ", mockList)

	return mockList, nil
}
