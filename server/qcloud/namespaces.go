package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/DDog/server/tool"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"net/url"
	"strings"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

func CreateNamespace(w http.ResponseWriter, r *http.Request) {

	//region := r.URL.Query().Get("region")
	//if region == "" {
	//	tool.ReturnError(w, errors.New(_const.RegionNotFound))
	//	return
	//}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	desc := r.URL.Query().Get("desc")
	if desc == "" {
		desc = "create-by-ddog"
	}

	name = strings.Replace(strings.ToLower(name), " ", "-", -1)
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

	q.SetDebug(_const.DEBUG)

	if err = q.CreateNamespace(); err != nil {
		tool.ReturnError(w, err)
		return
	}

	err = mongo.SaveNamespace(q)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	return
}

func Deletenamespace(w http.ResponseWriter, r *http.Request) {
	//region := r.URL.Query().Get("region")
	//if region == "" {
	//	tool.ReturnError(w, errors.New(_const.RegionNotFound))
	//	return
	//}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, errors.New(_const.RegionNotFound))
		return
	}

	name = strings.Replace(strings.ToLower(name), " ", "-", -1)
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
		tool.ReturnError(w, err)
		return
	}

	err = mongo.DeleteNamespaceByName(q.ClusterId, name)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	return
}

func CheckNamespace(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	rm := _const.RespMsg{
		Code: 1000,
		Msg:  "Namespace Exist",
	}

	name = strings.Replace(strings.ToLower(name), " ", "-", -1)
	_, err := mongo.GetNamespaceByName(clusterid, name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			md, err := metadata.GetMdByClusterID(clusterid)
			if err != nil {
				tool.ReturnError(w, errors.New(_const.RegionNotFound))
				return
			}

			q := namespace.NSpace{
				Pub: public.Public{
					Region:   md.Region,
					SecretId: md.Sid,
				},
				SecretKey: md.Skey,
				ClusterId: clusterid,
				Name:      url.QueryEscape(name),
				Desc:      url.QueryEscape("create-by-ddog"),
			}

			q.SetDebug(_const.DEBUG)
			if err = q.CreateNamespace(); err != nil {
				tool.ReturnError(w, err)
				return
			}

			err = mongo.SaveNamespace(q)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}
			rm.Code = 1001
			rm.Msg = "Create New Namespace"
		} else {
			tool.ReturnError(w, err)
			return
		}
	}

	//ns, err := UnMarshalNamespace(nsi)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//name = strings.Replace(strings.ToLower(name), " ", "-", -1)
	//if ns.ClusterID == "" {
	//
	//}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(&rm)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Write(data)

}

func UnMarshalNamespace(ns interface{}) (namespace namespace.NSInfo_data_namespaces, err error) {
	data, err := bson.Marshal(ns)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &namespace)
	if err != nil {
		return
	}

	return
}
