package metadata

import (
	"net/http"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/etcd"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/bridge"
)

type metaData struct {
	Sid    string `json:"secret_id"`
	Skey   string `json:"secret_key"`
	Region string `json:"region"`
}

func Startup(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil{
		server.ReturnError(w, err)
		return
	}

	var md metaData
	err = json.Unmarshal(data, &md)
	if err != nil{
		server.ReturnError(w, err)
		return
	}

	if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo, md.Sid); err != nil{
		server.ReturnError(w, err)
		return
	}

	if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo, md.Skey); err != nil{
		server.ReturnError(w, err)
		return
	}

	if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdRegionInfo, md.Region); err != nil{
		server.ReturnError(w, err)
		return
	}

	bridge.GetMetaChan() <- 1
	return
}
