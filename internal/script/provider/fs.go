package provider

import (
	"io"
	"os"
)

type Fs interface {
	Open(path string) (io.ReadCloser, error)
}

var DefaultFs Fs = OsFs{}

type OsFs struct{}

func (fs OsFs) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}
