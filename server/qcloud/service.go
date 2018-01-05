package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/server"
	"errors"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/metadata"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"github.com/andy-zhangtao/DDog/server/container"
	"gopkg.in/mgo.v2/bson"
	"log"
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
	if err != nil {
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

func RunService(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("svcid")
	if id == "" {
		server.ReturnError(w, errors.New(_const.SvcIDNotFound))
		return
	}

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		server.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	conf, err := mongo.GetSvcConfByID(id)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var cf svcconf.SvcConf
	//var ok bool
	//if cf, ok = conf.(svcconf.SvcConf); !ok {
	//	server.ReturnError(w, errors.New("Get Svc Conf Error "+reflect.TypeOf(conf).String()+fmt.Sprintf("[%v]", conf)))
	//	return
	//}

	data, err := bson.Marshal(conf)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	err = bson.Unmarshal(data, &cf)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var cluster cvm.ClusterInfo_data_clusters

	cs, err := mongo.GetClusterById(clusterid)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err = bson.Marshal(cs)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	err = bson.Unmarshal(data, &cluster)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	log.Println(cluster)
	md, err := metadata.GetMetaData(_const.RegionMap[cluster.Region])
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   cluster.ClusterId,
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
		server.ReturnError(w, err)
		return
	}

	for _, cn := range containers {
		var cnns container.Container
		data, err = bson.Marshal(cn)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		err = bson.Unmarshal(data, &cnns)
		if err != nil {
			server.ReturnError(w, err)
			return
		}
		//
		//if cnns, ok = cn.(container.Container); !ok {
		//	server.ReturnError(w, errors.New("Get Container Info Error "+reflect.TypeOf(conf).String()))
		//	return
		//}

		cons = append(cons, service.Containers{
			ContainerName: cnns.Name,
			Image:         cnns.Img,
		})
	}

	q.Containers = cons

	resp, err := q.CreateNewSerivce()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
