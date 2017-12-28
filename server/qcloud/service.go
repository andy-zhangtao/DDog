package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/server"
	"errors"
	"encoding/json"
)

func GetSampleSVCInfo(w http.ResponseWriter, r *http.Request) {
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

	allnamespace := r.Header.Get("allnamespace")
	if allnamespace == "" {
		allnamespace = "0"
	}

	q := service.Svc{
		Pub: public.Public{
			Region:   region,
			SecretId: sid,
		},
		ClusterId:    cid,
		Namespace:    namespace,
		Allnamespace: allnamespace,
		SecretKey:    key,
	}

	service, err := q.QuerySampleInfo()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(service)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
