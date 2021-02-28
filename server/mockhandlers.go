//Package server for mocking server definition.
package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

//MockResponses to return
type MockResponses struct {
	Path         string
	Method       string
	NextResponse int
	Responses    []Response
}

//AllMock to return
type AllMocks struct {
	Responses []MockResponses
}

//SetMocks set up mock handlers
func SetMocks(router *mux.Router, resources []MockResource) {

	var allMocks []MockResponses

	for _, resource := range resources {
		logger.Info("Got resource: ", resource.Resource)

		for _, mock := range resource.Mocks {

			logger.Info("Got resource: ", resource.Resource)

			path := mock.Mock.Request.Path

			logger.Info("Got path: ", path)

			if len(strings.TrimSpace(path)) > 0 && !strings.HasPrefix(strings.TrimSpace(path), "/") {
				path = resource.Resource + "/" + path
			} else {
				path = resource.Resource + path
			}
			mock.Mock.Request.Path = path

			method := mock.Mock.Request.Method
			logger.Debug("Got Resource Path: ", path)
			logger.Debug("Got mock.Mock.Path: ", mock.Mock.Request.Path)
			logger.Debug("Got Resource Method: ", method)

			var responseMocks []Response
			for _, mockResponse := range mock.Mock.Responses {
				responseMock := &Response{
					Headers:         mockResponse.Response.Headers,
					Status:          mockResponse.Response.Status,
					Body:            mockResponse.Response.Body,
					Delay:           mockResponse.Response.Delay,
					SkipEvery:       mockResponse.Response.SkipEvery,
					ResponseCounter: mockResponse.Response.ResponseCounter,
				}
				responseMocks = append(responseMocks, *responseMock)
			}
			mockResponse := &MockResponses{
				Path:         mock.Mock.Request.Path,
				Method:       mock.Mock.Request.Method,
				NextResponse: 0,
				Responses:    responseMocks,
			}

			logger.Info("Creating Handler for resource and method ", mock.Mock.Request.Path, " :  ", method)
			router.HandleFunc(path, mockResponse.mockingHandler).Methods(method)

			allMocks = append(allMocks, *mockResponse)
		}
	}
	mockDefault := &AllMocks{
		Responses: allMocks,
	}
	router.HandleFunc("/", mockDefault.mockingDefaultHandler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
}

//mockingHandler
func (mockResponse *MockResponses) mockingHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug("In Handler for Response Path: ", mockResponse.Path)
	logger.Debug("In Handler for Response Method: ", mockResponse.Method)

	var response Response
	for {
		response = mockResponse.Responses[mockResponse.NextResponse]
		if response.SkipEvery > 0 && response.ResponseCounter < response.SkipEvery {
			response.ResponseCounter++
			continue
		}
		mockResponse.NextResponse++
		if mockResponse.NextResponse > len(mockResponse.Responses) {
			mockResponse.NextResponse = 0
		}
		break
	}

	logger.Debug("In Handler for Response Status: ", response.Status)
	if response.Delay > 1 {
		logger.Debug("In Handler for Delay Now Sleeping (ms): ", response.Delay)
		time.Sleep(time.Duration(response.Delay) * time.Millisecond)
	}

	for _, headerString := range response.Headers {
		logger.Debug("Got headerString ", headerString)
		headerStringVal := strings.Split(headerString, ":")
		if len(headerStringVal) == 2 {
			w.Header().Set(headerStringVal[0], headerStringVal[1])
		} else {
			logger.Warn("Skipping Invalid Response Header : ", headerString)
		}
	}
	w.WriteHeader(response.Status)
	w.Write([]byte(response.Body))
}

//mockingDefaultHandler is the default handler
func (mocks *AllMocks) mockingDefaultHandler(w http.ResponseWriter, _ *http.Request) {
	mockingJSON, err := json.Marshal(mocks)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		logger.Error("Error encoding mocking YAML to json ", err)
		w.Write([]byte("Server Error Processing JSON Encoding"))
	}
	w.Write([]byte(mockingJSON))
}
