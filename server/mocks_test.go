//Package server for mocking server definition.
package server

import (
	"testing"
)

//TestValidateRequired tests for required input
func TestProcessMocks(t *testing.T) {
	ProcessMocks("../mocks")
}
