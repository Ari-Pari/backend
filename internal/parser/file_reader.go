package parser

import (
	"io"
	"os"
)

type FileReader interface {
	Open(name string) (io.ReadCloser, error)
}

type osFileReader struct{}

func (osFileReader) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

var DefaultFileReader FileReader = osFileReader{}
