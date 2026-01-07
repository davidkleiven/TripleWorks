package testutils

import "net/http"

type FailingResponseWriter struct {
	http.ResponseWriter
	WriteErr    error
	WantToWrite []byte
}

func (f *FailingResponseWriter) Write(content []byte) (int, error) {
	f.WantToWrite = content
	return 0, f.WriteErr
}
