package parser

import (
	"os"
)

type FileReader interface {
	Open(name string) (*os.File, error)
}

type osFileReader struct{}

func (osFileReader) Open(name string) (*os.File, error) {
	return os.Open(name)
}

var DefaultFileReader FileReader = osFileReader{}
