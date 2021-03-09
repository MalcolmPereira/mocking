//Package mockingserver for mocking server definition.
package mockingserver

import (
	"testing"

	logger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//TestLoglevel
func TestLoglevel(t *testing.T) {
	assert.NotNil(t, logger.GetLevel(), "Processed mockYAML is null")
	assert.Equal(t, "info", logger.GetLevel().String(), "Default LogLevel is invalid")
	LogLevel("INFO")
	assert.Equal(t, "info", logger.GetLevel().String(), "INFO LogLevel is invalid")
	LogLevel("info")
	assert.Equal(t, "info", logger.GetLevel().String(), "INFO LogLevel is invalid")
	LogLevel("DEBUG")
	assert.Equal(t, "debug", logger.GetLevel().String(), "DEBUG LogLevel is invalid")
	LogLevel("debug")
	assert.Equal(t, "debug", logger.GetLevel().String(), "DEBUG LogLevel is invalid")
	LogLevel("ERROR")
	assert.Equal(t, "error", logger.GetLevel().String(), "ERROR LogLevel is invalid")
	LogLevel("error")
	assert.Equal(t, "error", logger.GetLevel().String(), "ERROR LogLevel is invalid")
	LogLevel("WARN")
	assert.Equal(t, "warning", logger.GetLevel().String(), "WARN LogLevel is invalid")
	LogLevel("warn")
	assert.Equal(t, "warning", logger.GetLevel().String(), "WARN LogLevel is invalid")
	LogLevel("TRACE")
	assert.Equal(t, "trace", logger.GetLevel().String(), "TRACE LogLevel is invalid")
	LogLevel("trace")
	assert.Equal(t, "trace", logger.GetLevel().String(), "TRACE LogLevel is invalid")
}
