// Package main for ckureit utility
package main

import (
	"os"

	"github.com/malcolmpereira/mocking/server"
	logger "github.com/sirupsen/logrus"
)

//main functin for mocking utility
func main() {
	logger.Info("Start mocking")

	var mockingYAML = "mocking.yaml"
	if len(os.Args) > 1 {
		mockingYAML = os.Args[1]
	} else {
		logger.Info("No mocking.yaml specified as arguments, defaulting to mocking.yaml in current directory")
	}
	logger.Info("Mocking api starting with configurations:  ", mockingYAML)

	server.StartServer(mockingYAML)
}
