//Package mockingserver for mocking server definition.
package mockingserver

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

//TestSetMocks tests for Router
func TestSetMocks(t *testing.T) {
	mocksFolders := []string{"../mocking_demo/test_resources/test_mocks"}
	mocks, err := ProcessMocks(mocksFolders)
	assert.NotNil(t, mocks, "Processed mocks invalid")
	assert.Nil(t, err, "Processed mocks error is invalid")
	assert.True(t, len(mocks) > 0, "Invalid Mocks processing")
	router := mux.NewRouter()
	SetMocks(router, mocks)
}
