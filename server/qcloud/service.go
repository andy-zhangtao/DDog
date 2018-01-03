package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/server"
	"errors"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/metadata"
)

func GetSampleSVCInfo(w http.ResponseWriter, r *http.Request) {
	//sid := r.Header.Get("secretId")
	//if sid == "" {
	//	server.ReturnError(w, errors.New("SecretId Can not be empty"))
	//	return
	//}
	//
	//key := r.Header.Get("secretKey")
	//if key == "" {
	//	server.ReturnError(w, errors.New("SecretKey Can not be empty"))
	//	return
	//}


	region := r.URL.Query().Get("region")
	if region == "" {
		server.ReturnError(w, errors.New("Region Can not be empty"))
		return
	}

	cid := r.URL.Query().Get("clusterid")
	if cid == "" {
		server.ReturnError(w, errors.New("Clusterid Can not be empty"))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = "default"
	}

	allnamespace := r.URL.Query().Get("allnamespace")
	if allnamespace == "" {
		allnamespace = "0"
	}

	md, err := metadata.GetMetaData(region)
	if err != nil{
		server.ReturnError(w, err)
		return
	}

	q := service.Svc{
		Pub: public.Public{
			Region:   region,
			SecretId: md.Sid,
		},
		ClusterId:    cid,
		Namespace:    namespace,
		Allnamespace: allnamespace,
		SecretKey:    md.Skey,
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
