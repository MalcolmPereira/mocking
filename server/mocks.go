//Package server for mocking server definition.
package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var mockResourceMap map[string]string = make(map[string]string)

//MockResource struct wraps mocks
type MockResource struct {
	FileName string
	Name     string `yaml:"name"`
	Resource string `yaml:"resource"`
	Mocks    []struct {
		Mock struct {
			Path    string `yaml:"path"`
			Request struct {
				Method string `yaml:"method"`
			} `yaml:"request"`
			Responses []struct {
				Response struct {
					Headers         []string `yaml:"headers"`
					Status          int      `yaml:"status"`
					Body            string   `yaml:"body"`
					Delay           int      `yaml:"delay"`
					SkipEvery       int      `yaml:"skipevery"`
					ResponseCounter int
				} `yaml:"response"`
			} `yaml:"responses"`
		} `yaml:"mock"`
	} `yaml:"mocks"`
}

//ProcessMocks processes mocks available in the folder list
func ProcessMocks(mocksFolders []string) ([]MockResource, error) {
	var mockResourceList []MockResource
	for _, mockFolder := range mocksFolders {
		logger.Debug("Start processing mocks folder : ", mockFolder)
		mocks, err := processMockFolder(mockFolder)
		if err != nil {
			logger.Debug("Error processing mocks folder : ", mockFolder)
			logger.Error("Error processing mocks folders: ", err)
		}
		mockResourceList = append(mockResourceList, mocks...)
	}
	return mockResourceList, nil
}

//processMockFolder processing mocking yaml files under this folder
func processMockFolder(mockFolder string) ([]MockResource, error) {
	logger.Debug("Reading mocks folders : ", mockFolder)
	files, err := ioutil.ReadDir(mockFolder)
	if err != nil {
		logger.Error("Error reading mocks folders", err)
		return nil, errors.New("Error reading mocks folders")
	}
	logger.Debug("Done Reading mocks folders : ", mockFolder)
	logger.Debug("Processing files in  mocks folders : ", mockFolder)
	return processMockFiles(mockFolder, files)
}

//processMockFiles processing yaml files in the mock folder
func processMockFiles(mockFolder string, files []os.FileInfo) ([]MockResource, error) {

	var mockList []MockResource

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		logger.Debug("Reading mock yaml file : ", file.Name())
		fileBytes, err := ioutil.ReadFile(mockFolder + string(os.PathSeparator) + file.Name())
		if err != nil {
			logger.Error("Error reading mock yaml file ", err)
			return nil, errors.New("Error reading mock file")
		}
		logger.Debug("Done Reading mock yaml file  : ", file.Name())

		logger.Debug("Unmarshaling yaml file : ", file.Name())
		mockResource := MockResource{}
		err = yaml.Unmarshal(fileBytes, &mockResource)
		if err != nil {
			logger.Error("Error processing mocking YAML file : ", file.Name())
			logger.Error("Error processing mocking YAML file ", err)
			continue
		}
		logger.Debug("Done Unmarshaling yaml file : ", file.Name())

		err = validateMockResource(&mockResource, file.Name())

		if err != nil {
			logger.Error("Error validating mocking YAML file : ", file.Name())
			logger.Error("Error validating mocking YAML file ", err)
			continue
		}
		mockResource.FileName = file.Name()
		mockList = append(mockList, mockResource)
	}
	return mockList, nil
}

//validateMockResource
func validateMockResource(mockresource *MockResource, fileName string) error {
	if len(strings.TrimSpace(mockresource.Name)) == 0 {
		logger.Error("Invalid Mock Name, valid mock name required in mockResouce: " + fileName)
		return errors.New("Invalid mock name, valid mock name required in mockResouce: " + fileName)
	}
	if len(strings.TrimSpace(mockresource.Resource)) == 0 {
		logger.Error("Invalid Mock Resource, valid mock resource required in mockResouce: " + fileName)
		return errors.New("Invalid Mock Resource, valid mock resource required in mockResouce: " + fileName)
	}
	if !strings.HasPrefix(strings.TrimSpace(mockresource.Resource), "/") {
		mockresource.Resource = "/" + mockresource.Resource
	}

	for _, mock := range mockresource.Mocks {
		if len(strings.TrimSpace(mock.Mock.Request.Method)) == 0 {
			logger.Error("Invalid Mock Method, valid mock method required for request in mockResouce: " + fileName)
			return errors.New("Invalid Mock Method, valid mock method required for request in mockResouce: " + fileName)
		}
		resourcePath := mockresource.Resource + mock.Mock.Path + mock.Mock.Request.Method
		resourceFile := mockResourceMap[mockresource.Resource+mock.Mock.Path+mock.Mock.Request.Method]
		if len(strings.TrimSpace(resourceFile)) != 0 {
			logger.Error("Invalid Mock definiton, Duplicate Path: " + resourcePath + " , found in: " + resourceFile + " and " + fileName)
			return errors.New("Invalid Mock definiton, Duplicate Path: " + resourcePath + " , found in: " + resourceFile + " and " + fileName)
		}
		mockResourceMap[resourcePath] = fileName

		if len(mock.Mock.Responses) == 0 {
			logger.Error("Invalid Mock Responses, valid mock reponses required for request in mockResouce: " + fileName)
			return errors.New("Invalid Mock Responses, valid mock reponses required for request in mockResouce: " + fileName)
		}

		for _, mockresponse := range mock.Mock.Responses {
			if len(strings.TrimSpace(http.StatusText(mockresponse.Response.Status))) == 0 {
				logger.Error("Invalid Mock Responses Status Code , valid mock reponse status code  required for request  " + resourcePath + " in mockResouce: " + fileName)
				return errors.New("Invalid Mock Responses Status Code , valid mock reponse status code required for request  " + resourcePath + " in mockResouce: " + fileName)
			}
			statusString := strconv.Itoa(mockresponse.Response.Status)
			resourcePathStatus := resourcePath + statusString
			resourcePathStatusDup := mockResourceMap[resourcePathStatus]
			if len(strings.TrimSpace(resourcePathStatusDup)) != 0 {
				logger.Error("Invalid Mock definiton, Duplicate Status Code : " + statusString + " , found for : " + resourcePath + " int " + fileName)
				return errors.New("Invalid Mock definiton, Duplicate Status Code : " + statusString + " , found for : " + resourcePath + " int " + fileName)
			}
			mockresponse.Response.ResponseCounter = 0
		}
	}
	return nil
}
