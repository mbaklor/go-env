package env_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mbaklor/go-env"
	"github.com/stretchr/testify/assert"
)

func setEnvs(t *testing.T) {
	t.Setenv("FIRST", "a string")
	t.Setenv("SECOND", "1234")
	t.Setenv("THIRD", "more string")
	t.Setenv("FOURTH", "true")
	t.Setenv("FIFTH", "-987")

	t.Setenv("STRING", "test")
	t.Setenv("INT", "345")
	t.Setenv("BOOL", "true")
	t.Setenv("SHOULD_ZERO", "false")
	t.Setenv("NEST_STRING", "nested")
	t.Setenv("NEST_STRING_PTR", "nested pointer")

}

func concatDelim[T any](data []T, delim string) string {
	del := delim
	if delim == "" {
		del = string(os.PathListSeparator)
	}
	s := make([]string, 0, len(data))
	for _, d := range data {
		s = append(s, fmt.Sprintf("%v", d))
	}
	return strings.Join(s, del)
}

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
		Fourth bool   `env:"-"`
		Nested nestedStruct
	}

	type unsetVars struct {
		UnsetString string `env:"UNSET_STRING"`
		UnsetInt    int    `env:"UNSET_INT"`
		UnsetBool   bool   `env:"UNSET_BOOL"`
	}

	setEnvs(t)

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
	assert.Equal(t, false, tt.Fourth)
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

	setEnvs(t)

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

	t.Setenv("NO_STRING", "test")
	err = env.Load(&s)
	assert.NoError(t, err)
	assert.NotNil(t, s.NilStruct)
}

func TestSliceLoad(t *testing.T) {
	type sliceVars struct {
		StrSlice   []string `env:"STR_SLICE"`
		IntSlice   []int    `env:"INT_SLICE"`
		DelimSlice []string `env:"DELIM_SLICE,delim=,"`
	}

	var s sliceVars
	err := env.Load(&s)
	assert.NoError(t, err)
	assert.Nil(t, s.StrSlice)
	assert.Nil(t, s.IntSlice)
	assert.Nil(t, s.DelimSlice)

	var strSlice = []string{
		"test1",
		"test2",
		"test3",
	}
	var intSlice = []int{
		1,
		2,
		3,
	}
	var delimSlice = []string{
		"test3",
		"test2",
		"test1",
	}

	t.Setenv("STR_SLICE", concatDelim(strSlice, ""))
	t.Setenv("INT_SLICE", concatDelim(intSlice, ""))
	t.Setenv("DELIM_SLICE", concatDelim(delimSlice, ","))

	err = env.Load(&s)
	assert.NoError(t, err)
	assert.Equal(t, strSlice, s.StrSlice)
	assert.Equal(t, intSlice, s.IntSlice)
	assert.Equal(t, delimSlice, s.DelimSlice)
}

type MyTime struct {
	time.Time
}

func (m *MyTime) UnmarshalText(b []byte) error {
	t, err := time.Parse("15:04:05", string(b))
	if err != nil {
		return fmt.Errorf("UnmarshalText for %s: %w", string(b), err)
	}
	m.Time = t
	return nil
}

type MyDate struct {
	time.Time
}

func (m *MyDate) UnmarshalEnv(s string) error {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("UnmarshalText for %s: %w", s, err)
	}
	m.Time = t
	return nil
}

func TestTimeLoad(t *testing.T) {
	type timeVars struct {
		String  string        `env:"STRING"`
		Now     time.Time     `env:"NOW"`
		TimeVar MyTime        `env:"TIME"`
		DateVar MyDate        `env:"DATE"`
		Since   time.Duration `env:"SINCE"`
	}

	now := time.Now()

	t.Setenv("STRING", "test")
	t.Setenv("NOW", now.Format(time.RFC3339Nano))
	t.Setenv("TIME", "21:30:00")
	t.Setenv("DATE", "2026-09-17")
	t.Setenv("SINCE", "2h30m")

	var v timeVars
	err := env.Load(&v)
	assert.NoError(t, err)
	assert.Equal(t, "test", v.String)
	assert.Equal(t, now.UTC(), v.Now.UTC())
	assert.Equal(t, "09:30:00 PM", v.TimeVar.Format("03:04:05 PM"))
	assert.Equal(t, "17/09/26", v.DateVar.Format("02/01/06"))
	assert.Equal(t, 150.0, v.Since.Minutes())
}
