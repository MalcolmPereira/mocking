//Package mockingserver for mocking server definition.
package mockingserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestProcessMockingYAML valid test
func TestProcessMockingYAML(t *testing.T) {
	mockingYAML, err := ProcessMockingYAML("../mocking_demo/test_resources/mocking_test.yaml")
	if err != nil {
		t.Fatalf("FAIL TestProcessMockingYAML test for valid input")
	}
	assert.Nil(t, err, "Processed mockYAML error ")
	assert.NotNil(t, mockingYAML, "Processed mockYAML is null")
	assert.Equal(t, "0.1", mockingYAML.Version, "Mocking server version is invalid")
	assert.Equal(t, "MockingServer", mockingYAML.Name, "Mocking server name is invalid")
	assert.Equal(t, "2021", mockingYAML.Port, "Mocking server port is invalid")
	assert.Equal(t, "DEBUG", mockingYAML.Loglevel, "Mocking server log level is invalid")
	assert.Equal(t, "../mocking_demo/certificates/mockingServer.crt", mockingYAML.HTTPS.CertFile, "Mocking server cert file is invalid")
	assert.Equal(t, "../mocking_demo/certificates/mockingServer.key", mockingYAML.HTTPS.CertKeyFile, "Mocking server cert key file is invalid")
	assert.Equal(t, 4, len(mockingYAML.MockFolders), "Mocking server mock folders is invalid")
}

//TestProcessMockingYAML_Invalid invalid test
func TestProcessMockingYAML_Invalid(t *testing.T) {
	_, err := ProcessMockingYAML("../mocking_demo/test_resources/mocking_test_invalid.yaml")
	if err == nil {
		t.Fatalf("FAIL TestProcessMockingYAML_Invalid test for invalidInput")
	}
}
