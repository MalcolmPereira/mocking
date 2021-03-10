# Mocking API Server

## Configuration Files

1. mocking.yaml

    This configuration file configures mocking api server.

    The file can be named anything as long as it is valid yaml with configurations listed below.

        version: (Required) This should be set to 0.1 for now.

        name: (Required) Name for mocking api server e.g. MockingServer

        port: (Required) Port that will be used for running mocking api server, e.g. 2021

        loglevel: (Required) Logging level, e.g. INFO, other values include DEBUG, WARN, ERROR

        https: (Optional) If using https, specify https option and certificate files. If not using https this option can be ignored.

        certfile: (Required when https) This is the TLS certificate file,  e.g. ./mocking_demo/certificates/mockingServer.crt
  
        certkeyfile: (Required) This is TLS certificate key file,  e.g. ./mocking_demo/certificates/mockingServer.key

        mockfoldersfiles: (Required) List of folder to find mocks that mock server will use to serve mock responses.
        both folder containing mock.yaml files or a mock.yaml file can be specified.
           
            e.g. 
                 - ./mocking_demo/demo_mocks1
                 - ./mocking_demo/demo_mocks1
                 - ./mocking_demo/demo_mocks3
                 - ./mocking_demo/demo_mock.yaml 

2. mocks.yaml

    This is mock that will be served by mocking api server.

    Multiple mock files can be specified for the mocking api server. The combination of Resource + Path + Method must be unique, any duplicates will be ignored by  mocking api server. Multiple responses can be specified for Resource + Path + Method combination which will be iterated through sequentially.

    The mock configuration is specified in yaml and can be named anything as long as it is valid yaml and consists of following settings.

        mockversion: (Required) This should be set to 0.1

        name: (Required) Named for the mock , can be any string example:- Mock Demo

        resource: (Required) The resource that is being mocked example: /mock , so the mock URI would consists of http(s)://<SERVER>:<PORT>/mock

        mocks: (Required) This consists of array of mock, each mock serves mock for the resource configured with respect to path, methods and various responses. mocks contains collections of mock

        mock: (Required) This is a mock for a unique combination of Resource + Path + Method. Each mock is made up of request and collection of responses.

        request: This is the mock request which will be handled by the mocking api server which is made up of path and method.
        
        path: This is any path under resource defined for mock at the global level, if the path is not specified the global resource path will be used else this path will be appended to the global resource path to make up resource + "/" + path

        method: This is the request method that will be handled by the mocking api server.


