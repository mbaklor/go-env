package env_test

import (
	"os"
	"testing"

	"github.com/mbaklor/go-env"
	"github.com/stretchr/testify/assert"
)

func writeEnvFiles() func() {
	f, err := os.Create(".env")
	if err != nil {
		return func() {}
	}
	defer f.Close()
	f.WriteString(`HOST=127.0.0.1
PORT=8080`)

	c, err := os.Create(".env.local")
	if err != nil {
		return func() { os.Remove(".env") }
	}
	defer c.Close()
	c.WriteString(`# this is a comment
LOG_LEVEL=DEBUG # this is another comment

# and a blank line
CACHE=true
`)

	return func() {
		os.Remove(".env")
		os.Remove(".env.local")
	}
}

func TestLoadFile(t *testing.T) {
	t.Cleanup(writeEnvFiles())

	type fileEnv struct {
		Host     string `env:"HOST"`
		Port     int    `env:"PORT"`
		LogLevel string `env:"LOG_LEVEL"`
		Cache    bool   `env:"CACHE"`
	}

	err := env.LoadFiles(".env.local", ".env.develop")
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	var s fileEnv
	err = env.Load(&s)
	assert.NoError(t, err)
	assert.Zero(t, s)

	err = env.LoadFiles(".env.local")
	assert.NoError(t, err)
	err = env.Load(&s)
	assert.NoError(t, err)
	assert.Equal(t, "DEBUG", s.LogLevel)
	assert.Equal(t, true, s.Cache)

	err = env.LoadFiles()
	assert.NoError(t, err)
}
