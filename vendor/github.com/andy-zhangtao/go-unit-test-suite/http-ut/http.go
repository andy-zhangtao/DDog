package http_ut

import "net/http"

type TestResponseWriter struct {
	HeaderMap   http.Header
	InputLength int
	Input       []byte
	StatusCode  int
}

func (trw *TestResponseWriter) Header() http.Header {
	if trw.HeaderMap == nil {
		trw.HeaderMap = http.Header{}
	}

	return trw.HeaderMap
}

func (trw *TestResponseWriter) Write(input []byte) (int, error) {
	trw.Input = input
	trw.InputLength = len(trw.Input)
	return trw.InputLength, nil
}

func (trw *TestResponseWriter) WriteHeader(code int) {
	trw.StatusCode = code
	return
}
