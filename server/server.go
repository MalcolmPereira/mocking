//Package server for mocking server definition.
package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Current Supported Version
const currentVersion string = "0.1"

//MockingYAML
var mockingYAML *MockingYAML

//HTTP Server
var server *http.Server

//Mocks
var mocks []Mocks

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

//MockParam
type MockParam struct {
	requestHeaders  []string
	requestBody     string
	responseHeaders []string
	responseBody    string
	responseStatus  int
}

//StartServer parses mocking.yaml file to start mocking api
func StartServer(yamlFilePath string) {
	logger.Info("Start processing YAML file: ", yamlFilePath)
	var err error = nil
	mockingYAML, err = processMockingYAML(yamlFilePath)
	if err != nil {
		logger.Error("Error processing mocking YAML file", err)
		return
	}
	logger.Info("Processed YAML file: ", mockingYAML)

	mocks, err = ProcessMocks(mockingYAML.Mocks)

	if err != nil {
		logger.Error("Error processing mocks ", err)
		return
	}
	logger.Info("Processed mocks: ", mocks)

	start()
}

//start starts the Mocking API Server
func start() {
	r := mux.NewRouter()

	r.HandleFunc("/", mockingDefaultHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)

	for _, mock := range mocks {

		resource := mock.Resource

		logger.Info("Got resource: ", resource)

		for _, m := range mock.Mocks {
			method := m.Mock.Request.Method
			logger.Info("Got method: ", method)

			mockParam := &MockParam{
				requestHeaders:  m.Mock.Request.Headers,
				requestBody:     m.Mock.Request.Body,
				responseHeaders: m.Mock.Response.Headers,
				responseBody:    m.Mock.Response.Body,
				responseStatus:  m.Mock.Response.Status,
			}

			logger.Info("Creating Handler for  resource and method ", resource, method)

			r.HandleFunc(resource, mockParam.mockingHandler).Methods(method)
		}
	}

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
			http.MethodPatch,
			http.MethodTrace,
		},
		AllowCredentials: true,
		Debug:            true,
	}).Handler(r)
	server = &http.Server{
		Addr:         "0.0.0.0:" + mockingYAML.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}
	var err error = nil
	if len(strings.Trim(mockingYAML.HTTPS.CertFile, " ")) > 0 && len(strings.Trim(mockingYAML.HTTPS.CertKeyFile, " ")) > 0 {
		logger.Info("Got HTTPS Server Configuration Cert File: ", mockingYAML.HTTPS.CertFile)
		logger.Info("Got HTTPS Server Configuration Cert Key File: ", mockingYAML.HTTPS.CertKeyFile)
		logger.Info("Starting HTTPS Server on : 0.0.0.0:", mockingYAML.Port)
		err = server.ListenAndServeTLS(mockingYAML.HTTPS.CertFile, mockingYAML.HTTPS.CertKeyFile)
	} else {
		logger.Info("Starting HTTP Server on : 0.0.0.0:", mockingYAML.Port)
		err = server.ListenAndServe()
	}

	if err != nil {
		logger.Error("Server Error, failed starting Server  ", err)
		log.Fatalf("Server Error : %v", err)
	}
}

//mockingDefaultHandler is the default handler
func mockingDefaultHandler(w http.ResponseWriter, _ *http.Request) {
	mockingJSON, err := json.Marshal(mockingYAML)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Error encoding mocking YAML to json ", err)
		w.Write([]byte("Server Error Processing JSON Encoding"))
	}
	w.Write([]byte(mockingJSON))
}

//mockingHandler
func (mock *MockParam) mockingHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("In Handler for ", mock.responseBody)

	for _, headerString := range mock.responseHeaders {
		logger.Info("Got headerString ", headerString)
		headerStringVal := strings.Split(headerString, ":")
		logger.Info("Got headerStringVal ", headerStringVal)
		logger.Info("Got headerStringVal[0] ", headerStringVal[0])
		logger.Info("Got headerStringVal[1] ", headerStringVal[1])
		w.Header().Set(headerStringVal[0], headerStringVal[1])
	}
	w.WriteHeader(mock.responseStatus)
	w.Write([]byte(mock.responseBody))

	//
	// for _, headerString := range mock.responseHeaders {
	// 	headerStringVal := strings.Split(headerString, ":")
	// 	w.Header().Set(headerStringVal[0], headerStringVal[1])
	// }
}

//processMockingYAML parses and validates mocking api yaml file
func processMockingYAML(yamlFilePath string) (*MockingYAML, error) {
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
	if strings.Trim(mockingYAML.Version, " ") != currentVersion {
		logger.Error("Invalid Mocking API Version, supported version is: " + currentVersion)
		return errors.New("Invalid Mocking API Version, supported version is: " + currentVersion)
	}
	if len(strings.Trim(mockingYAML.Name, " ")) == 0 {
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
	return nil
}
