package bench

import (
	"bytes"
	"time"
)

type mockedFSType struct {
	err  error // An arbitary error
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

type mockedRenderer struct{}

func (m *mockedRenderer) Render() error {
	return nil
}
func (m *mockedRenderer) AddReceivedDataLength(url string, contentLength int64)         {}
func (m *mockedRenderer) SetTotalDuration(duration time.Duration)                       {}
func (m *mockedRenderer) AddResponseTime(url string, time time.Duration)                {}
func (m *mockedRenderer) AddResponseStatusCode(url string, statusCode int, failed bool) {}
func (m *mockedRenderer) AddTimedoutResponse(url string)                                {}
