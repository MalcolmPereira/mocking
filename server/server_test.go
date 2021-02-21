//Package server for mocking server definition.
package server

import (
	"fmt"
	"testing"
)

//TestValidateRequired tests for required input
func TestProcessYAML(t *testing.T) {
	mockingYAML, err := processMockingYAML("../mocking.yaml")
	if err != nil {
		t.Fatalf("FAIL TestProcessYAML test for valid input")
	}
	fmt.Println("mockingYAML ", mockingYAML)
}