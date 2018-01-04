package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/DDog/server"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"net/url"
	"strings"
)

func CreateNamespace(w http.ResponseWriter, r *http.Request) {

	region := r.URL.Query().Get("region")
	if region == "" {
		server.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		server.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		server.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	md, err := metadata.GetMetaData(region)
	if err != nil {
		server.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	desc := r.URL.Query().Get("desc")
	if desc == "" {
		desc = "create-by-ddog"
	}

	q := namespace.NSpace{
		Pub: public.Public{
			Region:   md.Region,
			SecretId: md.Sid,
		},
		SecretKey: md.Skey,
		ClusterId: clusterid,
		Name:      url.QueryEscape(name),
		Desc:      url.QueryEscape(desc),
	}

	if err = q.CreateNamespace(); err != nil {
		server.ReturnError(w, err)
		return
	}

	return
}

func Deletenamespace(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		server.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		server.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		server.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	md, err := metadata.GetMetaData(region)
	if err != nil {
		server.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	q := namespace.NSpace{
		Pub: public.Public{
			Region:   md.Region,
			SecretId: md.Sid,
		},
		SecretKey: md.Skey,
		ClusterId: clusterid,
		Rmname:    strings.Split(name, ";"),
	}

	if err = q.DeleteNamespace(); err != nil {
		server.ReturnError(w, err)
		return
	}

	return
}
