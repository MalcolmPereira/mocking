//Package server for mocking server definition.
package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

//MockParam
type MockParam struct {
	requestHeaders  []string
	requestBody     string
	responseHeaders []string
	responseBody    string
	responseStatus  int
}

//SetMocks set up mock handlers
func SetMocks(router *mux.Router, mocks []Mocks) {

	router.HandleFunc("/", mockingDefaultHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)

	for _, mock := range mocks {

		resource := mock.Resource

		logger.Debug("Got resource: ", resource)

		for _, m := range mock.Mocks {
			method := m.Mock.Request.Method
			logger.Debug("Got method: ", method)

			mockParam := &MockParam{
				requestHeaders:  m.Mock.Request.Headers,
				requestBody:     m.Mock.Request.Body,
				responseHeaders: m.Mock.Response.Headers,
				responseBody:    m.Mock.Response.Body,
				responseStatus:  m.Mock.Response.Status,
			}
			logger.Info("Creating Handler for resource and method ", resource, " :  ", method)
			router.HandleFunc(resource, mockParam.mockingHandler).Methods(method)
		}
	}

}

//mockingHandler
func (mock *MockParam) mockingHandler(w http.ResponseWriter, r *http.Request) {

	logger.Debug("In Handler for ", mock.responseBody)

	for _, headerString := range mock.responseHeaders {
		logger.Debug("Got headerString ", headerString)
		headerStringVal := strings.Split(headerString, ":")
		logger.Debug("Got headerStringVal ", headerStringVal)
		logger.Debug("Got headerStringVal[0] ", headerStringVal[0])
		logger.Debug("Got headerStringVal[1] ", headerStringVal[1])
		w.Header().Set(headerStringVal[0], headerStringVal[1])
	}
	w.WriteHeader(mock.responseStatus)
	w.Write([]byte(mock.responseBody))
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
