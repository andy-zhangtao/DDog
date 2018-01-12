package metadata

import (
	"testing"
	"net/http"
	"github.com/andy-zhangtao/go-unit-test-suite/http-ut"
	"github.com/andy-zhangtao/go-unit-test-suite/io-ut"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"net/url"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/const"
)

func TestStartup(t *testing.T) {
	go func() {
		ret := <-bridge.GetMetaChan()
		assert.Equal(t, 1, ret, "The status code should be 1")
	}()

	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	md := metadata.MetaData{
		Sid:    "123",
		Skey:   "456",
		Region: "local",
	}

	data, _ := json.Marshal(&md)

	r.Body = io_ut.GetReadCloser(string(data))

	metadata.DeleteMetaDataByRegion(md.Region)
	Startup(w, r)

	assert.Equal(t, 0, w.StatusCode, "The status code should be 200")
}

func TestGetMetaData(t *testing.T) {
	metadata.SaveMetaData(metadata.MetaData{
		Sid:    "123",
		Skey:   "456",
		Region: "local",})

	md, err := GetMetaData("local")
	assert.Empty(t, err)

	assert.Equal(t, "123", md.Sid, "The Sid should be 123")
	assert.Equal(t, "456", md.Skey, "The Sid should be 456")

	md, err = GetMetaData("local1")
	assert.EqualError(t, err, "local1 Metadata 获取为空")
}

func TestGetMetaDataWithHttp(t *testing.T) {
	r := new(http.Request)
	w := new(http_ut.TestResponseWriter)

	r.URL = new(url.URL)
	r.URL.RawQuery = "region=local"
	GetMetaDataWithHttp(w, r)

	var md metadata.MetaData

	json.Unmarshal(w.Input, &md)

	assert.Equal(t, "123", md.Sid, "Get respones error!")
	assert.Equal(t, "456", md.Skey, "Get respones error!")
	assert.Equal(t, "local", md.Region, "Get respones error!")
}

func TestGetMdByClusterID(t *testing.T) {
	_const.RegionMap["l"] = "local"
	c := make(map[string]string)
	c["clusterid"] = "cid"
	c["region"] = "l"
	mongo.SaveCluster(c)

	md, err := GetMdByClusterID("cid")

	assert.Empty(t, err)

	assert.Equal(t, "123", md.Sid, "Get respones error!")
	assert.Equal(t, "456", md.Skey, "Get respones error!")
	assert.Equal(t, "local", md.Region, "Get respones error!")

}
