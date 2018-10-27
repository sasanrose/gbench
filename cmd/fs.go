package cmd

import (
	"io"
	"os"
)

var fs fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (file, error)
}

type file interface {
	io.Closer
	io.Reader
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error) { return os.Open(name) }
