// 维护命名空间信息
package handler

import (
	"net/http"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server"
	"encoding/json"
	"strconv"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	ns "github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"github.com/andy-zhangtao/DDog/server/etcd"
)

type NameSpace struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
	ClusterID string `json:"cluster_id"`
}

const (
	CloudEtcdNameSpaceInfo = "/namespace/info"
)

func QueryNameSpace(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var ns NameSpace

	err = json.Unmarshal(data, &ns)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	isSave := r.URL.Query().Get("save")
	if isSave == "" || isSave != "true" {
		isSave = "false"
	}

	save, err := strconv.ParseBool(isSave)
	if err != nil {
		save = false
	}

	nsinfo, err := ns.saveNSInfo(save)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err = json.Marshal(nsinfo)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (this NameSpace) saveNSInfo(save bool) (*ns.NSInfo, error) {
	c := ns.NSpace{
		Pub: public.Public{
			Region:   this.Region,
			SecretId: this.SecretId,
		},
		SecretKey: this.SecretKey,
		ClusterId: this.ClusterID,
	}

	c.SetDebug(true)

	ns, err := c.QueryNSInfo()
	if err != nil {
		return nil, err
	}

	if save {
		data, err := json.Marshal(ns.Data.Namespaces)
		if err != nil {
			return nil, err
		}
		err = etcd.Put(CloudEtcdRootPath+"/"+c.Pub.Region+CloudEtcdNameSpaceInfo, string(data))
		if err != nil {
			return nil, err
		}
	}

	return ns, nil
}
