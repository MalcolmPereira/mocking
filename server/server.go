//Package server for mocking server definition.
package server

import (
	"context"
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

//MockingYAML struct wraps mocking api configuration
type MockingYAML struct {
	Version  string `yaml:"version"`
	Name     string `yaml:"name"`
	Port     string `yaml:"port"`
	Loglevel string `yaml:"loglevel"`
	HTTPS    struct {
		CertFile    string `yaml:"certfile"`
		CertKeyFile string `yaml:"certkeyfile"`
	}
}

//StartServer parses mocking.yaml file to start mocking api
func StartServer(yamlFilePath string) {
	logger.Info("Start processing YAML file: ", yamlFilePath)

	var err error = nil
	mockingYAML, err = processMockingYAML(yamlFilePath)
	if err != nil {
		logger.Error("Error processing mocking YAML file: %v", err)
		return
	}
	logger.Info("Processed YAML file: ", mockingYAML)

	logger.Info("Starting Mocking API Server: ")
	start()
	logger.Info("Done Starting Mocking API Server: ")

	loopForMock()
}

func loopForMock() {
	var isFirstTime bool = true
	for {
		logger.Info("Watching for Changes: ")
		time.Sleep(15 * time.Second)

		if !isFirstTime {
			logger.Info("Stoping API Server: ")
			stopServer()
			logger.Info("Done Stoping API Server: ")

			logger.Info("Starting Mocking API Server: ")
			start()
			logger.Info("Done Starting Mocking API Server: ")
		} else {
			logger.Info("First Time Skip: ")
			isFirstTime = false
		}

	}
}

//processMockingYAML parses and validates mocking api yaml file
func processMockingYAML(yamlFilePath string) (*MockingYAML, error) {
	mockingYAML := MockingYAML{}

	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Error("Error reading mocking YAML file: %v", err)
		return nil, errors.New("Error reading mocking YAML file")
	}

	err = yaml.Unmarshal(yamlFile, &mockingYAML)
	if err != nil {
		logger.Error("Error processing mocking YAML file: %v", err)
		return nil, errors.New("Error processing mocking YAML file")
	}

	err = validateYAML(&mockingYAML)
	if err != nil {
		logger.Error("Error validating mocking YAML file: %v", err)
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
		logger.Error("Error processing mocking YAML file, Invalid Number Value for Port: %v", err)
		return errors.New("Error processing mocking YAML file, Invalid Number Value for Port ")
	}
	if portNum < 1080 || portNum > 65535 {
		logger.Error("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535 ")
		return errors.New("Invalid Mocking Server Port, valid Mocking Server Port required: 1081 - 65535")
	}

	return nil
}

//start starts the Mocking API Server
func start() {

	r := mux.NewRouter()

	r.HandleFunc("/", mockingDefaultHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)

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

	go func() {

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

	}()
}

//stopServer stops the HTTP Server
func stopServer() error {
	logger.Info("Stopping HTTP Server on : 0.0.0.0:", mockingYAML.Port)
	err := server.Shutdown(context.TODO())
	if err != nil {
		logger.Error("Stopping HTTP Server on : 0.0.0.0:", mockingYAML.Port)
		return errors.New("Error Stopping HTTP Server")
	}
	return nil
}

//mockingDefaultHandler is the default handler
func mockingDefaultHandler(w http.ResponseWriter, _ *http.Request) {
	mockingJSON, err := json.Marshal(mockingYAML)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Error encoding mocking YAML to json:  %v ", err)
		w.Write([]byte("Server Error Processing JSON Encoding"))
	}
	w.Write([]byte(mockingJSON))
}
