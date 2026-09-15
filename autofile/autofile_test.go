package autofile_test

import (
	"os"
	"testing"

	"github.com/mbaklor/go-env/autofile"
	"github.com/stretchr/testify/assert"
)

func TestAutoLoadFile(t *testing.T) {
	assert.ErrorIs(t, autofile.AutoLoadFileErr, os.ErrNotExist)
}
