package env_test

import (
	"os"
	"testing"

	"github.com/mbaklor/env"
	"github.com/stretchr/testify/assert"
)

type testType struct {
	First  string `env:"FIRST"`
	Second int    `env:"SECOND"`
	Third  string `env:"THIRD"`
	Nested nestedStruct
}

type nestedStruct struct {
	Fourth bool `env:"FOURTH"`
	Fifth  int
	Sixth  string
}

func setEnvVars() {
	os.Setenv("FIRST", "a string")
	os.Setenv("SECOND", "1234")
	os.Setenv("THIRD", "more string")
	os.Setenv("FOURTH", "true")
	os.Setenv("FIFTH", "-987")
}

func TestLoad(t *testing.T) {
	setEnvVars()
	fail := "should fail"
	err := env.Load(fail)
	assert.Error(t, err)

	var pt *testType = nil
	err = env.Load(pt)
	assert.Error(t, err)

	var tt testType
	err = env.Load(&tt)
	assert.NoError(t, err)
	assert.Equal(t, "a string", tt.First)
	assert.Equal(t, 1234, tt.Second)
	assert.Equal(t, "more string", tt.Third)
	assert.Equal(t, true, tt.Nested.Fourth)
	assert.Equal(t, -987, tt.Nested.Fifth)
	assert.Equal(t, "", tt.Nested.Sixth)

}
