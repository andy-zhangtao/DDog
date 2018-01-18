package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"errors"
	"encoding/json"

	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"

	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/gogather/zsort"
	"log"
	"github.com/andy-zhangtao/DDog/model/container"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"fmt"
)

func GetSampleSVCInfo(w http.ResponseWriter, r *http.Request) {
	//clusterid := r.URL.Query().Get("clusterid")
	//if clusterid == "" {
	//	tool.ReturnError(w, errors.New(_const.ClusterNotFound))
	//	return
	//}

	var id string
	var cf *svcconf.SvcConf
	var err error
	var nsme string
	name := r.URL.Query().Get("svcname")
	if name != "" {
		//	如果上传服务名称，则直接重新部署此服务
		nsme = r.URL.Query().Get("namespace")
		if nsme == "" {
			nsme = _const.DefaultNameSpace
			if nsme == "" {
				tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
				return
			}
		}
	} else {
		id = r.URL.Query().Get("id")
		if id == "" {
			tool.ReturnError(w, errors.New(_const.IDNotFound))
			return
		}
		cf, err = svcconf.GetSvcConfByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	if name == "" {
		name = cf.Name
	}
	if nsme == "" {
		nsme = cf.Namespace
	}

	name = strings.TrimSpace(name)
	nsme = strings.TrimSpace(nsme)
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Svc{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId: md.ClusterID,
		//ServiceName: name,
		Namespace: nsme,
		SecretKey: md.Skey,
	}

	service, err := q.QuerySampleInfo()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var data []byte
	for _, svc := range service.Data.Services {
		if svc.ServiceName == name {
			data, err = json.Marshal(svc)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// RunService 在K8s集群中创建服务
// svcname 服务配置名称, 此名称应该在创建服务之前首先创建。
// namespace 命名空间名称, 如果为空则为默认值
// upgrade 是否为升级操作. 默认为false。
// 当在创建服务时，会使用以下默认参数
// 1. 默认启用健康检查和准备就绪检查。
// 2. 上述两种检查使用TCP端口检查方式
// 3. 均针对容器对外暴露的端口进行检查，如果镜像构建未对外暴露端口，则不会对此镜像启用检查
// 4. 延时30秒后启动检查
// 5. 连续三次，间隔10秒，健康均失败则检查失败
// 6. 每次检查超时时间为5秒
func RunService(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("svcname")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		if nsme == "" {
			nsme = _const.DefaultNameSpace
		}
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	up := r.URL.Query().Get("upgrade")
	isUpgrade, err := strconv.ParseBool(up)
	if err != nil {
		log.Printf("[RunService] parsebool error [%s] request value [%v]", err.Error(), up)
		isUpgrade = false
	}

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if _const.DEBUG {
		log.Printf("[RunService] Svc Conf [%v]\n", cf)
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: cf.Name,
		ServiceDesc: cf.Desc,
		Replicas:    cf.Replicas,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)
	if len(cf.Netconf) > 0 {
		var pm []service.PortMappings
		for _, n := range cf.Netconf {
			p := service.PortMappings{}
			switch n.Protocol {
			case 0:
				p.Protocol = "TCP"
			case 1:
				p.Protocol = "UDP"
			}
			p.ContainerPort = n.InPort
			p.LbPort = n.OutPort
			pm = append(pm, p)
		}
		q.PortMappings = pm
		switch cf.Netconf[0].AccessType {
		case 0:
			q.AccessType = "ClusterIP"
		case 1:
			q.AccessType = "LoadBalancer"
		case 2:
			q.AccessType = "SvcLBTypeInner"
		}
	}

	var cons []service.Containers

	containers, err := container.GetAllContainersBySvc(cf.Name, cf.Namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if len(containers) == 0 {
		tool.ReturnError(w, errors.New(fmt.Sprintf("[Find Container Error][%s]svc[%s]namespace[%s]", _const.ContainerNotFound, cf.Name, cf.Namespace)))
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

		var hk []service.HealthCheck

		for _, n := range cn.Net {
			shk := service.HealthCheck{
				Type:        service.LiveCheck,
				UnhealthNum: 5,
				DelayTime:   30,
				CheckMethod: service.CheckMethodTCP,
			}
			shk.GenerateTCPCheck(n.InPort)

			hk = append(hk, shk)
			shk.Type = service.ReadyCheck
			hk = append(hk, shk)
		}

		cons = append(cons, service.Containers{
			ContainerName: cnns.Name,
			Image:         cnns.Img,
			HealthCheck:   hk,
		})
	}

	q.Containers = cons

	if _const.DEBUG {
		log.Printf("[RunService] QCloud Request [%v] Object Deploy Type [%v] \n", q, isUpgrade)
	}
	var resp *service.SvcSMData
	if isUpgrade {
		q.Strategy = "RollingUpdate"
		resp, err = q.RedeployService()
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
	w.Header().Set("EQXC-Run-Svc", "200")
	w.Write(data)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	var cf *svcconf.SvcConf
	var err error
	id := r.URL.Query().Get("id")
	if id == "" {
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

		cf, err = svcconf.GetSvcConfByName(name, nsme)
	} else {
		cf, err = svcconf.GetSvcConfByID(id)
	}

	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	md, err := metadata.GetMetaDataByRegion("")
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
	w.Header().Set("EQXC-Run-Svc", "200")
	w.Write(data)
}

func ReinstallService(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	var id string
	var cf *svcconf.SvcConf
	var err error
	var nsme string
	name := r.URL.Query().Get("svcname")
	if name != "" {
		//	如果上传服务名称，则直接重新部署此服务
		nsme = r.URL.Query().Get("namespace")
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	} else {
		id = r.URL.Query().Get("id")
		if id == "" {
			tool.ReturnError(w, errors.New(_const.IDNotFound))
			return
		}
		cf, err = svcconf.GetSvcConfByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	if name == "" {
		name = cf.Name
	}
	if nsme == "" {
		nsme = cf.Namespace
	}

	md, err := metadata.GetMetaDataByRegion("")
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
		ServiceName: name,
		Namespace:   nsme,
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
		nsme = _const.DefaultNameSpace
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if cf == nil {
		tool.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Svc{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId: md.ClusterID,
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
		if _const.DEBUG {
			log.Printf("[DeployService] Find Svc Dist:[%s] Current:[%s]\n", cf.Name, r.ServiceName)
		}
		if strings.Compare(r.ServiceName, cf.Name) == 0 {
			isUpgrade = true
			break
		}
	}

	oldPath := r.URL.RawQuery + "&namespace=" + cf.Namespace

	if isUpgrade {
		r.URL.RawQuery = oldPath + "&upgrade=true"
	} else {
		r.URL.RawQuery = oldPath + "&upgrade=false"
	}

	if _const.DEBUG {
		log.Printf("[DeployService] Deploy Type [%v] [%s] [%s]\n", isUpgrade, r.URL.String(), oldPath)
	}

	RunService(w, r)

}

func RunSvcGroup(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svcConfGroup := r.URL.Query().Get("svcgroup")
	if svcConfGroup == "" {
		tool.ReturnError(w, errors.New(_const.SvcGroupNotFound))
		return
	}

	if _const.DEBUG {
		log.Printf("[RunSvcGroup]clusterid:[%s]namespace:[%s]svcgroup:[%s]\n", clusterid, namespace, svcConfGroup)
	}

	sg, err := mongo.GetSvcConfGroupByName(svcConfGroup, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	svcg, err := svcconf.Unmarshal(sg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if _const.DEBUG {
		log.Printf("[RunSvcGroup]svcg:[%v]\n", svcg)
	}

	svcPair := zsort.SortByValue(svcg.SvcGroup)
	rawQuery := r.URL.RawQuery
	nd := strings.Index(rawQuery, "&svcname=")
	if nd > 0 {
		//clear query path
		rawQuery = rawQuery[:nd]
	}

	for i := len(svcPair) - 1; i >= 0; i -- {

		r.URL.RawQuery = rawQuery + "&svcname=" + svcPair[i].Key

		if _const.DEBUG {
			log.Printf("[RunSvcGroup]Deploy svcname :[%s] All header:[%v] \n", svcPair[i].Key, r.URL.Query())
		}

		w.Header().Del("EQXC-Run-Svc")
		DeployService(w, r)
		if _const.DEBUG {
			log.Printf("[RunSvcGroup]Deploy svcname :[%s] Response:[%v] \n", svcPair[i].Key, w.Header())
		}

		if w.Header().Get("EQXC-Run-Svc") != "200" {
			return
		}

	}
}

func ReinstallSvcGroup(w http.ResponseWriter, r *http.Request) {
	UninstallSvcGroup(w, r)
	RunSvcGroup(w, r)
}

func UninstallSvcGroup(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svcConfGroup := r.URL.Query().Get("svcgroup")
	if svcConfGroup == "" {
		tool.ReturnError(w, errors.New(_const.SvcGroupNotFound))
		return
	}

	if _const.DEBUG {
		log.Printf("[UninstallSvcGroup]clusterid:[%s]namespace:[%s]svcgroup:[%s]\n", clusterid, namespace, svcConfGroup)
	}

	sg, err := mongo.GetSvcConfGroupByName(svcConfGroup, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	svcg, err := svcconf.Unmarshal(sg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if _const.DEBUG {
		log.Printf("[UninstallSvcGroup]svcg:[%v]\n", svcg)
	}

	svcPair := zsort.SortByValue(svcg.SvcGroup)
	rawQuery := r.URL.RawQuery

	nd := strings.Index(rawQuery, "&svcname=")
	if nd > 0 {
		//clear query path
		rawQuery = rawQuery[:nd]
	}

	for i := len(svcPair) - 1; i >= 0; i -- {
		r.URL.RawQuery = rawQuery + "&svcname=" + svcPair[i].Key

		if _const.DEBUG {
			log.Printf("[UninstallSvcGroup]Delete svcname :[%s] All header:[%v] \n", svcPair[i].Key, r.URL.Query())
		}

		w.Header().Del("EQXC-Run-Svc")
		DeleteService(w, r)
		if _const.DEBUG {
			log.Printf("[UninstallSvcGroup]Delete svcname :[%s] Response:[%v] \n", svcPair[i].Key, w.Header())
		}

		if w.Header().Get("EQXC-Run-Svc") != "200" {
			return
		}
	}
}
