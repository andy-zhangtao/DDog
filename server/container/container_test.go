package container

import (
	"testing"
	"github.com/andy-zhangtao/go-unit-test-suite/http-ut"
	"net/http"
	"github.com/andy-zhangtao/go-unit-test-suite/io-ut"
	"github.com/andy-zhangtao/DDog/model/container"
	"encoding/json"
	"net/url"
	"github.com/stretchr/testify/assert"
	"github.com/andy-zhangtao/DDog/server/mongo"
)

var con = container.Container{
	Name: "UT",
	Img:  "UT_IMG",
	Svc:  "UT_SVC",
	Nsme: "UT_NSME",
	Idx:  1,
}

func TestCreateContainer(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	data, _ := json.Marshal(con)

	r.Body = io_ut.GetReadCloser(string(data))

	svc := make(map[string]string)
	svc["name"] = con.Svc
	svc["namespace"] = con.Nsme
	mongo.SaveSvcConfig(svc)
	CreateContainer(w, r)
	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
}

func TestGetContainer(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)
	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=UT_NSME&svc=UT_SVC"

	GetContainer(w, r)
	var cns []container.Container

	json.Unmarshal(w.Input, &cns)
	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
	assert.Equal(t, "UT", cns[0].Name, "The container name error")
}

func TestDeleteContainer(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=UT_NSME&svc=UT_SVC"

	DeleteContainer(w, r)

	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
}
