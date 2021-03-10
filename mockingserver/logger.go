//Package mockingserver for mocking server definition.
package mockingserver

import (
	"strings"

	logger "github.com/sirupsen/logrus"
)

//LogLevel configures the log level in the logger
func LogLevel(loglevel string) {
	if strings.EqualFold("DEBUG", loglevel) {
		logger.SetLevel(logger.DebugLevel)
	}
	if strings.EqualFold("INFO", loglevel) {
		logger.SetLevel(logger.InfoLevel)
	}
	if strings.EqualFold("WARN", loglevel) {
		logger.SetLevel(logger.WarnLevel)
	}
	if strings.EqualFold("ERROR", loglevel) {
		logger.SetLevel(logger.ErrorLevel)
	}
	if strings.EqualFold("TRACE", loglevel) {
		logger.SetLevel(logger.TraceLevel)
	}
}
