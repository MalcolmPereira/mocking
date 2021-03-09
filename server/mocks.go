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

//MockResource struct wraps a mock resource
type MockResource struct {
	FileName string
	Name     string  `yaml:"name"`
	Resource string  `yaml:"resource"`
	Mocks    []Mocks `yaml:"mocks"`
}

//Mocks struct wraps mocks
type Mocks struct {
	Mock Mock `yaml:"mock"`
}

//Mock struct wraps the request and reponses for mocking
type Mock struct {
	Request   Request     `yaml:"request"`
	Responses []Responses `yaml:"responses"`
}

//Request warps the mock request method and path
type Request struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
}

//Responses wraps the responses available for the request
type Responses struct {
	Response Response `yaml:"response"`
}

//Response wraps the response
type Response struct {
	Headers         []string `yaml:"headers"`
	Status          int      `yaml:"status"`
	Body            string   `yaml:"body"`
	File            string   `yaml:"file"`
	Delay           int      `yaml:"delay"`
	SkipEvery       int      `yaml:"skipevery"`
	ResponseCounter int
}

var mockResourceMap map[string]string = make(map[string]string)

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

	return validateMockRequest(mockresource.Mocks, mockresource.Resource, fileName)
}

func validateMockRequest(mocks []Mocks, resource string, fileName string) error {

	for _, mock := range mocks {
		if len(strings.TrimSpace(mock.Mock.Request.Method)) == 0 {
			logger.Error("Invalid Mock Method, valid mock method required for request in mockResouce: " + fileName)
			return errors.New("Invalid Mock Method, valid mock method required for request in mockResouce: " + fileName)
		}

		resourcePath := resource + mock.Mock.Request.Path + mock.Mock.Request.Method
		resourceFile := mockResourceMap[resource+mock.Mock.Request.Path+mock.Mock.Request.Method]
		if len(strings.TrimSpace(resourceFile)) != 0 {
			logger.Error("Invalid Mock definiton, Duplicate Path: " + resourcePath + " , found in: " + resourceFile + " and " + fileName)
			return errors.New("Invalid Mock definiton, Duplicate Path: " + resourcePath + " , found in: " + resourceFile + " and " + fileName)
		}
		mockResourceMap[resourcePath] = fileName

		if len(mock.Mock.Responses) == 0 {
			logger.Error("Invalid Mock Responses, valid mock reponses required for request in mockResouce: " + fileName)
			return errors.New("Invalid Mock Responses, valid mock reponses required for request in mockResouce: " + fileName)
		}

		err := validateMockResponse(mock.Mock.Responses, resourcePath, fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateMockResponse(responses []Responses, resourcePath string, fileName string) error {
	for _, mockresponse := range responses {
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
	return nil
}
