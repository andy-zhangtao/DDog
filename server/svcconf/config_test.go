package svcconf

import (
	"testing"
	"github.com/andy-zhangtao/go-unit-test-suite/http-ut"
	"net/http"
	"encoding/json"
	"github.com/andy-zhangtao/go-unit-test-suite/io-ut"
	"github.com/stretchr/testify/assert"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"net/url"
	"github.com/andy-zhangtao/DDog/model/svcconf"
)

var svc = svcconf.SvcConf{
	Name:      "SVCCONF",
	Desc:      "SVCCONF",
	Replicas:  2,
	Namespace: "SVC_NSME",
	Netconf: []svcconf.NetConfigure{
		svcconf.NetConfigure{AccessType: 2,
			InPort: 9000,
			OutPort: 8000,
			Protocol: 1,},
	},
}

var id string

func init() { mongo.DeleteSvcConfByNs(svc.Namespace) }

func TestCreateSvcConf(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	data, _ := json.Marshal(svc)

	r.Body = io_ut.GetReadCloser(string(data))

	CreateSvcConf(w, r)

	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
	assert.Equal(t, 24, len(string(w.Input)), "The svcconf id error")
	id = string(w.Input)
}

func TestGetSvcConf(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=SVC_NSME&id=" + id
	GetSvcConf(w, r)

	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
	var c svcconf.SvcConf
	json.Unmarshal(w.Input, &c)

	assert.Equal(t, svc.Name, c.Name, "The name error")
	assert.Equal(t, svc.Namespace, c.Namespace, "The name error")

	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=SVC_NSME"
	GetSvcConf(w, r)
	var cs []svcconf.SvcConf
	json.Unmarshal(w.Input, &cs)

	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
	assert.Equal(t, svc.Name, cs[0].Name, "The name error")
	assert.Equal(t, svc.Namespace, cs[0].Namespace, "The name error")
}

func TestUpgradeSvcConf(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=SVC_NSME&id=" + id

	svc.Replicas = 9
	data, _ := json.Marshal(svc)

	r.Body = io_ut.GetReadCloser(string(data))

	UpgradeSvcConf(w, r)
	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")

	GetSvcConf(w, r)
	var c svcconf.SvcConf
	json.Unmarshal(w.Input, &c)
	assert.Equal(t, 9, c.Replicas, "The repllicas should be 9")
}

func TestDeleteSvcConf(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	r.URL = new(url.URL)
	r.URL.RawQuery = "namespace=SVC_NSME&id=" + id

	DeleteSvcConf(w, r)
	assert.Equal(t, 0, w.StatusCode, "The status code should be 0")
	r.URL.RawQuery = "namespace=SVC_NSME&id=" + id
	GetSvcConf(w, r)

	assert.Equal(t, 500, w.StatusCode, "The status code should be 500")
}
