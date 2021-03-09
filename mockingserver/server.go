//Package mockingserver for mocking server definition.
package mockingserver

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	logger "github.com/sirupsen/logrus"
)

//MockingYAML
var mockingYAML *MockingYAML

//HTTP Server
var server *http.Server

//StartServer parses mocking.yaml file to start mocking api
func StartServer(yamlFilePath string) {
	logger.SetLevel(logger.DebugLevel)

	logger.Debug("Starting Mock Server with YAML ", yamlFilePath)
	var err error = nil
	mockingYAML, err = ProcessMockingYAML(yamlFilePath)

	if err != nil {
		log.Fatalln("Error processing mocking YAML file", err)
		return
	}

	logger.Debug("Processed Mock Server YAML File ", yamlFilePath)

	LogLevel(mockingYAML.Loglevel)

	logger.Debug("Start Processing Mocks")
	mocks, err := ProcessMocks(mockingYAML.MockFolders)
	if err != nil {
		logger.Error("Error processing mocks ", err)
		return
	}
	logger.Debug("End Processing Mocks, Mocks: ", mocks)

	logger.Debug("Start Get HTTP Server: ")
	server := getHTTPServer(mocks)
	logger.Debug("End Get HTTP Server: ", server)

	start(server)
}

//start starts the Mocking API Server
func getHTTPServer(mocks []MockResource) *http.Server {
	router := mux.NewRouter()

	SetMocks(router, mocks)

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
	}).Handler(router)

	server = &http.Server{
		Addr:         "0.0.0.0:" + mockingYAML.Port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	return server
}

//start starts the Mocking API Server
func start(server *http.Server) {
	logger.Debug("Starting HTTP Server")

	var err error = nil

	if len(strings.Trim(mockingYAML.HTTPS.CertFile, " ")) > 0 && len(strings.Trim(mockingYAML.HTTPS.CertKeyFile, " ")) > 0 {
		logger.Debug("Got HTTPS Server Configuration Cert File: ", mockingYAML.HTTPS.CertFile)
		logger.Debug("Got HTTPS Server Configuration Cert Key File: ", mockingYAML.HTTPS.CertKeyFile)
		logger.Info("Starting HTTPS Server on : 0.0.0.0:", mockingYAML.Port)
		err = server.ListenAndServeTLS(mockingYAML.HTTPS.CertFile, mockingYAML.HTTPS.CertKeyFile)
	} else {
		logger.Info("Starting HTTP Server on : 0.0.0.0:", mockingYAML.Port)
		err = server.ListenAndServe()
	}
	logger.Debug("Done Starting HTTP Server")

	if err != nil {
		logger.Error("Server Error, failed starting Server  ", err)
		log.Fatalf("Server Error : %v", err)
	}
}
