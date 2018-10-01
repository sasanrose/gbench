package bench

import (
	"bytes"
)

type mockedFSType struct {
	err  error // An arbitrary error
	file *mockedFileType
}

type mockedFileType struct {
	mockedReader *bytes.Buffer // A mocked reader to return what we want to test
}

func (m mockedFileType) Read(p []byte) (n int, err error) {
	return m.mockedReader.Read(p)
}

func (m mockedFileType) Close() error {
	return nil
}

func (m mockedFSType) Open(name string) (file, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.file, nil
}
