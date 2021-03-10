// Package main for mocking server utility
package main

import (
	"os"
	"strings"

	server "github.com/malcolmpereira/mocking/mockingserver"
	logger "github.com/sirupsen/logrus"
)

//main function for mocking utility
func main() {
	logger.Info("Start mocking..")
	var mockingYAML = "./mocking_demo/mocking.yaml"

	if len(os.Args) > 1 && len(strings.TrimSpace(os.Args[1])) > 0 {
		mockingYAML = os.Args[1]
		logger.Info("Mocking api starting with configurations:  ", mockingYAML)
	} else {
		logger.Info("No mocking.yaml specified as arguments, defaulting to mocking.yaml in current directory")
	}

	server.StartServer(mockingYAML)
}
