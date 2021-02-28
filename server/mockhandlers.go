//Package server for mocking server definition.
package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

//AllMocks to return
type AllMocks struct {
	Responses []MockResponses
}

//SetMocks set up mock handlers
func SetMocks(router *mux.Router, resources []MockResource) {

	var allMocks []MockResponses

	for _, resource := range resources {
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
			logger.Debug("Got Resource Path: ", mock.Mock.Request.Path)
			logger.Debug("Got Resource Method: ", mock.Mock.Request.Method)

			var responseMocks []Response
			for _, mockResponse := range mock.Mock.Responses {
				responseMock := &Response{
					Headers:         mockResponse.Response.Headers,
					Status:          mockResponse.Response.Status,
					Body:            mockResponse.Response.Body,
					File:            mockResponse.Response.File,
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

	err := r.ParseMultipartForm(28)
	if err != nil {
		logger.Info("Multipart Parsing Ignore", err)
	}

	var response Response
	for {
		response = mockResponse.Responses[mockResponse.NextResponse]

		if response.SkipEvery > 0 && response.ResponseCounter < response.SkipEvery {
			response.ResponseCounter++
			continue
		}

		mockResponse.NextResponse++

		if mockResponse.NextResponse >= len(mockResponse.Responses) {
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

	if len(strings.TrimSpace(response.Body)) > 0 {
		w.WriteHeader(response.Status)
		w.Write([]byte(response.Body))

	} else if len(strings.TrimSpace(response.File)) > 0 {
		logger.Info("Now Reading File ", response.File)
		filedata, fileErr := ioutil.ReadFile(response.File)

		if fileErr != nil {
			logger.Info("Error Reading File ", fileErr)
			w.WriteHeader(500)
			w.Write([]byte("Error Reading File: " + response.File))
		}

		w.WriteHeader(200)
		fileBuffer := bytes.NewBuffer(filedata)
		if _, err := fileBuffer.WriteTo(w); err != nil {
			fmt.Fprintf(w, "%s", err)
			logger.Error("Error writing file to http ", err)
		}

	} else {
		w.WriteHeader(response.Status)
		w.Write([]byte(""))
	}
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
