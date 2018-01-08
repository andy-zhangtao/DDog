package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"errors"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/metadata"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/andy-zhangtao/DDog/server/container"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"github.com/andy-zhangtao/DDog/server/tool"
)

func GetSampleSVCInfo(w http.ResponseWriter, r *http.Request) {
	//sid := r.Header.Get("secretId")
	//if sid == "" {
	//	tool.ReturnError(w, errors.New("SecretId Can not be empty"))
	//	return
	//}
	//
	//key := r.Header.Get("secretKey")
	//if key == "" {
	//	tool.ReturnError(w, errors.New("SecretKey Can not be empty"))
	//	return
	//}

	region := r.URL.Query().Get("region")
	if region == "" {
		tool.ReturnError(w, errors.New("Region Can not be empty"))
		return
	}

	cid := r.URL.Query().Get("clusterid")
	if cid == "" {
		tool.ReturnError(w, errors.New("Clusterid Can not be empty"))
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
	if err != nil {
		tool.ReturnError(w, err)
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
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(service)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func RunService(w http.ResponseWriter, r *http.Request) {
	//id := r.URL.Query().Get("svcid")
	//if id == "" {
	//	tool.ReturnError(w, errors.New(_const.SvcIDNotFound))
	//	return
	//}

	name := r.URL.Query().Get("svcname")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	up := r.URL.Query().Get("upgrade")
	isUpgrade, err := strconv.ParseBool(up)
	if err != nil {
		isUpgrade = false
	}

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	//cf, err := svcconf.GetSvcConfByID(id)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//conf, err := mongo.GetSvcConfByID(id)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//var cf svcconf.SvcConf
	//
	//data, err := bson.Marshal(conf)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//err = bson.Unmarshal(data, &cf)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}

	//var cluster cvm.ClusterInfo_data_clusters
	//
	//cs, err := mongo.GetClusterById(clusterid)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//data, err = bson.Marshal(cs)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//err = bson.Unmarshal(data, &cluster)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	////log.Println(cluster)
	//md, err := metadata.GetMetaData(_const.RegionMap[cluster.Region])
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}

	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   clusterid,
		ServiceName: cf.Name,
		ServiceDesc: cf.Desc,
		Replicas:    cf.Replicas,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
		PortMappings: service.PortMappings{
			LbPort:        cf.Netconf.OutPort,
			ContainerPort: cf.Netconf.InPort,
		},
	}

	switch cf.Netconf.Protocol {
	case 0:
		q.PortMappings.Protocol = "TCP"
	case 1:
		q.PortMappings.Protocol = "UDP"
	}

	q.SetDebug(true)
	switch cf.Netconf.AccessType {
	case 0:
		q.AccessType = "ClusterIP"
	case 1:
		q.AccessType = "LoadBalancer"
	case 2:
		q.AccessType = "SvcLBTypeInner"
	}

	var cons []service.Containers

	containers, err := mongo.GetContaienrBySvc(cf.Name, cf.Namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	for _, cn := range containers {
		var cnns container.Container
		data, err := bson.Marshal(cn)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		err = bson.Unmarshal(data, &cnns)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		cons = append(cons, service.Containers{
			ContainerName: cnns.Name,
			Image:         cnns.Img,
		})
	}

	q.Containers = cons

	var resp *service.SvcSMData
	if isUpgrade {
		q.Strategy = "RollingUpdate"
		resp, err = q.UpgradeService()
	} else {
		resp, err = q.CreateNewSerivce()
	}

	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		tool.ReturnError(w, errors.New(_const.IDNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	cf, err := svcconf.GetSvcConfByID(id)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	//conf, err := mongo.GetSvcConfByID(id)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//var cf svcconf.SvcConf
	//
	//data, err := bson.Marshal(conf)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//err = bson.Unmarshal(data, &cf)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}

	//var cluster cvm.ClusterInfo_data_clusters
	//
	//cs, err := mongo.GetClusterById(clusterid)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//data, err = bson.Marshal(cs)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//err = bson.Unmarshal(data, &cluster)
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//md, err := metadata.GetMetaData(_const.RegionMap[cluster.Region])
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   clusterid,
		ServiceName: cf.Name,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
	}

	resp, err := q.DeleteService()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ReinstallService(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		tool.ReturnError(w, errors.New(_const.IDNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	cf, err := svcconf.GetSvcConfByID(id)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   clusterid,
		ServiceName: cf.Name,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
	}

	resp, err := q.RedeployService()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DeployService(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("svcname")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	md, err := metadata.GetMdByClusterID(clusterid)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Svc{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId: clusterid,
		Namespace: cf.Namespace,
		SecretKey: md.Skey,
	}
	q.SetDebug(true)
	resp, err := q.QuerySampleInfo()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	isUpgrade := false
	for _, r := range resp.Data.Services {
		if strings.Compare(r.ServiceName, cf.Name) == 0 {
			isUpgrade = true
			break
		}
	}

	if isUpgrade {
		r.URL.Query().Set("isupgrade", "true")
	} else {
		r.URL.Query().Set("isupgrade", "false")
	}

	RunService(w, r)

}
