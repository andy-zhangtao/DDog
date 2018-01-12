package tool

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
)

type TestResponseWriter struct {
	headerMap   http.Header
	inputLength int
	input       []byte
	statusCode  int
}

func (trw *TestResponseWriter) Header() http.Header {
	if trw.headerMap == nil {
		trw.headerMap = http.Header{}
	}

	return trw.headerMap
}

func (trw *TestResponseWriter) Write(input []byte) (int, error) {
	trw.input = input
	trw.inputLength = len(trw.input)
	return trw.inputLength, nil
}

func (trw *TestResponseWriter) WriteHeader(code int) {
	trw.statusCode = code
	return
}

func TestIsNotFound(t *testing.T) {
	err := errors.New("Here is a error")
	found := IsNotFound(err)

	assert.Equal(t, false, found, "This should be False")

	err = errors.New("  not found")
	found = IsNotFound(err)
	assert.Equal(t, true, found, "This should be True")
}

func TestReturnError(t *testing.T) {
	w := new(TestResponseWriter)
	err := errors.New("This is a err")
	ReturnError(w, err)

	assert.Equal(t, "application/json", w.headerMap.Get("Content-Type"), "The header Content-Type wrong!")
	assert.Equal(t, "500", w.headerMap.Get("EQXC-Run-Svc"), "The header Content-Type wrong!")
	assert.Equal(t, 500, w.statusCode, "The status code should be 500")

	assert.Equal(t,"{\"code\":500,\"msg\":\"This is a err\"}",string(w.input),"Return data error")
}

func TestReturnResp(t *testing.T) {
	w := new(TestResponseWriter)
	ReturnResp(w, []byte("null"))

	assert.Equal(t,"{\"code\":0,\"msg\":\"Operation Succ! But donot find any data!\",\"data\":\"null\"}",string(w.input),"Return data error")
	ReturnResp(w, []byte("OK"))
	assert.Equal(t,"{\"code\":1000,\"msg\":\"Operation Succ!\",\"data\":\"OK\"}",string(w.input),"Return data error")
}
