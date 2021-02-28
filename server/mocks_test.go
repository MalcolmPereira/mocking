//Package server for mocking server definition.
package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestProcessMocks tests for required input
func TestProcessMocks(t *testing.T) {
	mocksFolders := []string{"../test_resources/test_mocks"}
	mocks, err := ProcessMocks(mocksFolders)
	assert.NotNil(t, mocks, "Processed mocks invalid")
	assert.Nil(t, err, "Processed mocks error is invalid")
	assert.True(t, len(mocks) > 0, "Invalid Mocks processing")
	for _, mock := range mocks {
		if strings.EqualFold(mock.Name, "Mock Demo") {
			assert.Equal(t, "Mock Demo", mock.Name, "Invalid Mock Name")
			assert.Equal(t, "/mock", mock.Resource, "Invalid Mock Resource")
			assert.True(t, (len(mock.Mocks) == 5), "Invalid Mocks")
		}
	}
}
