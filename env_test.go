package env_test

import (
	"os"
	"testing"

	"github.com/mbaklor/go-env"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type nestedStruct struct {
		Fourth bool `env:"FOURTH"`
		Fifth  int
		Sixth  string
	}

	type testType struct {
		First  string `env:"FIRST"`
		Second int    `env:"SECOND"`
		Third  string `env:"THIRD"`
		Nested nestedStruct
	}

	type unsetVars struct {
		UnsetString string `env:"UNSET_STRING"`
		UnsetInt    int    `env:"UNSET_INT"`
		UnsetBool   bool   `env:"UNSET_BOOL"`
	}

	os.Setenv("FIRST", "a string")
	os.Setenv("SECOND", "1234")
	os.Setenv("THIRD", "more string")
	os.Setenv("FOURTH", "true")
	os.Setenv("FIFTH", "-987")

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

	var unset unsetVars
	err = env.Load(&unset)
	assert.NoError(t, err)
	assert.Equal(t, "", unset.UnsetString)
	assert.Equal(t, 0, unset.UnsetInt)
	assert.Equal(t, false, unset.UnsetBool)
}
