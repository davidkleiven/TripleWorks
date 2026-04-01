package pkg

import (
	"errors"
	"io"
	"os"
)

type Opener interface {
	Open(name string) (io.ReadCloser, error)
}

type FsOpener struct{}

func (f *FsOpener) Open(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

type FailingReader struct{}

func (f *FailingReader) Read(p []byte) (int, error) {
	return 0, errors.New("could not read")
}

type FailingReadOpener struct{}

func (f *FailingReadOpener) Open(name string) (io.ReadCloser, error) {
	reader := &FailingReader{}
	return io.NopCloser(reader), nil
}
