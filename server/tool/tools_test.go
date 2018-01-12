package tool

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/andy-zhangtao/go-unit-test-suite/http-ut"
)

func TestIsNotFound(t *testing.T) {
	err := errors.New("Here is a error")
	found := IsNotFound(err)

	assert.Equal(t, false, found, "This should be False")

	err = errors.New("  not found")
	found = IsNotFound(err)
	assert.Equal(t, true, found, "This should be True")
}

func TestReturnError(t *testing.T) {
	w := new(http_ut.TestResponseWriter)
	err := errors.New("This is a err")
	ReturnError(w, err)

	assert.Equal(t, "application/json", w.HeaderMap.Get("Content-Type"), "The header Content-Type wrong!")
	assert.Equal(t, "500", w.HeaderMap.Get("EQXC-Run-Svc"), "The header Content-Type wrong!")
	assert.Equal(t, 500, w.StatusCode, "The status code should be 500")

	assert.Equal(t, "{\"code\":500,\"msg\":\"This is a err\"}", string(w.Input), "Return data error")
}

func TestReturnResp(t *testing.T) {
	w := new(http_ut.TestResponseWriter)
	ReturnResp(w, []byte("null"))

	assert.Equal(t, "{\"code\":0,\"msg\":\"Operation Succ! But donot find any data!\",\"data\":\"null\"}", string(w.Input), "Return data error")
	ReturnResp(w, []byte("OK"))
	assert.Equal(t, "{\"code\":1000,\"msg\":\"Operation Succ!\",\"data\":\"OK\"}", string(w.Input), "Return data error")
}
