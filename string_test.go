package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToEnvVar(t *testing.T) {
	st := stringToEnvVar("TestVar")
	assert.Equal(t, "TEST_VAR", st)

	st = stringToEnvVar("TestAPIVar")
	assert.Equal(t, "TEST_API_VAR", st)

	st = stringToEnvVar("TestAPIVarHTTP")
	assert.Equal(t, "TEST_API_VAR_HTTP", st)

	st = stringToEnvVar("HTTPTestAPIVarHTTP")
	assert.Equal(t, "HTTP_TEST_API_VAR_HTTP", st)
}
