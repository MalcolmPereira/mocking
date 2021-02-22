// Package main for ckureit utility
package main

import (
	"os"
	"strings"

	"github.com/malcolmpereira/mocking/server"
	logger "github.com/sirupsen/logrus"
)

//main functin for mocking utility
func main() {
	logger.Info("Start mocking..")
	var mockingYAML = "mocking.yaml"

	if len(os.Args) > 1 && len(strings.TrimSpace(os.Args[1])) > 0 {
		mockingYAML = os.Args[1]
		logger.Info("Mocking api starting with configurations:  ", mockingYAML)
	} else {
		logger.Info("No mocking.yaml specified as arguments, defaulting to mocking.yaml in current directory")
	}

	server.StartServer(mockingYAML)
}
