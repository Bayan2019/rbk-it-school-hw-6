package spy

import (
	"io"
	"net/http"
)

type SpyReadCloser struct {
	io.ReadCloser
	bytesRead int
}

func (r *SpyReadCloser) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	r.bytesRead += n
	return n, err
}

type SpyResponseWriter struct {
	http.ResponseWriter
	bytesWritten int
	StatusCode   int
}

func (w *SpyResponseWriter) Write(p []byte) (int, error) {
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(p)
	w.bytesWritten += n
	return n, err
}

func (w *SpyResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *SpyResponseWriter) WriteHeaderKeyValue(key, value string) {
	// w.statusCode = statusCode
	// w.ResponseWriter.WriteHeader(statusCode)

	w.Header().Set(key, value)
}
