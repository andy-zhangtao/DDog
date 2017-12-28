package qcloud

import (
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"net/http"
	"github.com/andy-zhangtao/DDog/server"
	"errors"
	"encoding/json"
)

func GetClusterNodes(w http.ResponseWriter, r *http.Request) {

	sid := r.Header.Get("secretId")
	if sid == "" {
		server.ReturnError(w, errors.New("SecretId Can not be empty"))
		return
	}

	key := r.Header.Get("secretKey")
	if key == "" {
		server.ReturnError(w, errors.New("SecretKey Can not be empty"))
		return
	}

	region := r.Header.Get("region")
	if region == "" {
		server.ReturnError(w, errors.New("Region Can not be empty"))
		return
	}

	cid := r.Header.Get("clusterid")
	if cid == "" {
		server.ReturnError(w, errors.New("Clusterid Can not be empty"))
		return
	}

	namespace := r.Header.Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	q := cvm.Cluster{
		Pub: public.Public{
			Region:   region,
			SecretId: sid,
		},
		Cid:       cid,
		Namespace: namespace,
		Offset:    0,
		Limit:     20,
		SecretKey: key,
	}

	nodes, err := q.QueryClusterNodes()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
