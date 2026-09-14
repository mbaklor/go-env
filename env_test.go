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

func TestPtrLoad(t *testing.T) {
	type nestedPtr struct {
		NestString    string
		NestStringPtr *string
	}

	type shouldNilPtr struct {
		NoString  string  `env:"NO_STRING"`
		NoInt     int     `env:"NO_INT"`
		NoBool    bool    `env:"NO_BOOL"`
		NoPointer *string `env:"NO_POINTER"`
	}

	type prtVars struct {
		String     *string `env:"STRING"`
		Int        *int    `env:"INT"`
		Bool       *bool   `env:"BOOL"`
		ShouldZero *bool   `env:"SHOULD_ZERO"`
		ShouldNil  *bool
		Struct     *nestedPtr
		NilStruct  *shouldNilPtr
	}

	var s prtVars
	err := env.Load(&s)
	assert.NoError(t, err)
	assert.Nil(t, s.String)
	assert.Nil(t, s.Int)
	assert.Nil(t, s.Bool)
	assert.Nil(t, s.ShouldZero)
	assert.Nil(t, s.ShouldNil)
	assert.Nil(t, s.Struct)
	assert.Nil(t, s.NilStruct)

	os.Setenv("STRING", "test")
	os.Setenv("INT", "345")
	os.Setenv("BOOL", "true")
	os.Setenv("SHOULD_ZERO", "false")
	os.Setenv("NEST_STRING", "nested")
	os.Setenv("NEST_STRING_PTR", "nested pointer")

	err = env.Load(&s)
	assert.NoError(t, err)
	assert.Equal(t, "test", *s.String)
	assert.Equal(t, 345, *s.Int)
	assert.Equal(t, true, *s.Bool)
	assert.Equal(t, false, *s.ShouldZero)
	assert.Nil(t, s.ShouldNil)
	assert.Nil(t, s.NilStruct)
	assert.Equal(t, "nested", s.Struct.NestString)
	assert.Equal(t, "nested pointer", *s.Struct.NestStringPtr)

	os.Setenv("NO_STRING", "test")
	err = env.Load(&s)
	assert.NoError(t, err)
	assert.NotNil(t, s.NilStruct)
}
