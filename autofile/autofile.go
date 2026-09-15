package autofile

import "github.com/mbaklor/go-env"

var AutoLoadFileErr error

func init() {
	AutoLoadFileErr = env.LoadFiles()
}
