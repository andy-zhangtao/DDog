// 维护集群配置信息
package handler

import (
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"net/http"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server"
	"encoding/json"
	"strconv"
	"github.com/andy-zhangtao/DDog/server/etcd"
)

type ClusterHandler struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
}

const (
	CloudEtcdRootPath = "/cloud/meta"
	CloudEtcdClusterInfo = "/cluster/info"
)

// saveClusterInfo 查询集群信息并持久化到etcd
// 如果选择持久化，则会覆盖旧值
func (this ClusterHandler) saveClusterInfo(save bool) (*cvm.ClusterInfo, error) {

	c := cvm.Cluster{
		Pub: public.Public{
			Region:   this.Region,
			SecretId: this.SecretId,
		},
		SecretKey: this.SecretKey,
	}

	c.SetDebug(true)

	cinfo, err := c.QueryClusters()
	if err != nil {
		return nil, err
	}

	if save {
		data, err := json.Marshal(cinfo.Data.Clusters)
		if err != nil{
			return nil, err
		}
		err = etcd.Put(CloudEtcdRootPath+"/"+c.Pub.Region+CloudEtcdClusterInfo,string(data))
		if err != nil{
			return nil, err
		}
	}

	return cinfo, nil
}

// QueryClusterInfo 查询集群列表信息
// 列出此账户下所有的集群信息
func QueryClusterInfo(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
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

	var ch ClusterHandler

	err = json.Unmarshal(data, &ch)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	cinfo, err := ch.saveClusterInfo(save)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err = json.Marshal(cinfo)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
