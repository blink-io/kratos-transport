package http3

import "net/http"

type mockResponseWriter struct {
	StatusCode int
	Data       []byte
	header     http.Header
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.Data = b
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}
