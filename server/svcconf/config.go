package svcconf

import (
	"encoding/json"
	"errors"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/container"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Operation struct{}

// CPort 容器端口数据
type CPort struct {
	Name string                   `json:"name"`
	Img  string                   `json:"img"`
	Net  []container.NetConfigure `json:"net"`
}

const (
	ModuleName = "svcconf"
)

func CreateSvcConf(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var conf svcconf.SvcConf

	err = json.Unmarshal(data, &conf)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkConf(&conf); err != nil {
		tool.ReturnError(w, err)
		return
	}

	cf, err := mongo.GetSvcConfByName(conf.Name, conf.Namespace)
	if cf != nil {
		tool.ReturnError(w, errors.New(_const.SvcConfExist))
		return
	}

	if conf.Replicas == 0 {
		conf.Replicas = 1
	}

	conf.Id = bson.NewObjectId()
	if err = mongo.SaveSvcConfig(conf); err != nil {
		tool.ReturnError(w, err)
		return
	}
	w.Write([]byte(conf.Id.Hex()))
	return
}

func GetSvcConf(w http.ResponseWriter, r *http.Request) {
	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		nsme = strings.Replace(strings.ToLower(nsme), " ", "-", -1)
		conf, err := mongo.GetSvcConfNs(nsme)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(conf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		w.Write(data)
	} else {
		conf, err := mongo.GetSvcConfByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(conf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		w.Write(data)
	}

}
func QuerySvcConf(w http.ResponseWriter, r *http.Request) {
	type scs struct {
		Status int           `json:"status"`
		Port   []map[int]int `json:"port"`
	}
	svc := r.URL.Query().Get("svc")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
	}

	scf, err := svcconf.GetSvcConfByName(svc, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var s scs
	var p []map[int]int

	s.Status = scf.Status

	for _, f := range scf.Netconf {
		port := make(map[int]int)
		port[f.InPort] = f.OutPort
		p = append(p, port)
	}

	s.Port = p
	data, err := json.Marshal(&s)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}

// DeleteSvcConf 销毁服务配置
// 销毁服务所有资源,包括数据库资源, 服务实例资源
func DeleteSvcConf(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svcname")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
		if namespace == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	data, err := json.Marshal(_const.DestoryMsg{
		Svcname:   svc,
		Namespace: namespace,
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{"Marshal DestoryMsg Error": err,}).Error(ModuleName)
		tool.ReturnError(w, errors.New(err.Error()))
		return
	}

	err = bridge.SendDestoryMsg(string(data))
	if err != nil {
		logrus.WithFields(logrus.Fields{"DestorySvc Error": err,}).Error(ModuleName)
		tool.ReturnError(w, errors.New(err.Error()))
		return
	}

	//oper := Operation{}
	//err = oper.DeleteSvcConf(_const.DestoryMsg{
	//	Svcname:   svc,
	//	Namespace: namespace,
	//})
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	tool.ReturnResp(w, []byte("Delete Succ!"))
	return

}
func checkConf(conf *svcconf.SvcConf) error {
	if conf.Name == "" {
		return errors.New(_const.NameNotFound)
	}

	if conf.Namespace == "" {
		conf.Namespace = _const.DefaultNameSpace
	}

	/*默认实例数为2，要不然没法做蓝绿发布*/
	if conf.Replicas == 0 {
		conf.Replicas = 2
	}

	conf.Name = strings.Replace(strings.ToLower(conf.Name), " ", "-", -1)
	conf.Namespace = strings.Replace(strings.ToLower(conf.Namespace), " ", "-", -1)

	if len(conf.Netconf) > 0 {
		for _, n := range conf.Netconf {
			if n.Protocol != 0 && n.Protocol != 1 {
				return errors.New(_const.LbProtocolError)
			}

			if n.InPort == 0 || n.OutPort == 0 {
				return errors.New(_const.LbPortError)
			}

			if n.AccessType != 0 && n.AccessType != 1 && n.AccessType != 2 {
				return errors.New(_const.AccessTypeError)
			}
		}

	}

	return nil
}

func UpgradeSvcConf(w http.ResponseWriter, r *http.Request) {
	cid := r.URL.Query().Get("id")
	if cid == "" {
		tool.ReturnError(w, errors.New(_const.IDNotFound))
		return
	}

	c, err := mongo.GetSvcConfByID(cid)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	conf, err := svcconf.Conver(c)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var nc svcconf.SvcConf

	err = json.Unmarshal(data, &nc)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if nc.Replicas > 0 {
		conf.Replicas = nc.Replicas
	}

	if len(nc.Netconf) > 0 {
		for i, n := range nc.Netconf {
			if n.Protocol != 0 && n.Protocol != 1 {
				tool.ReturnError(w, errors.New(_const.LbProtocolError))
				return
			} else {
				conf.Netconf[i].Protocol = n.Protocol
			}

			if n.InPort == 0 || n.OutPort == 0 {
				tool.ReturnError(w, errors.New(_const.LbPortError))
				return
			} else if n.InPort > 0 {
				conf.Netconf[i].InPort = n.InPort
			} else if n.OutPort > 0 {
				conf.Netconf[i].OutPort = n.OutPort
			}

			if n.AccessType != 0 && n.AccessType != 1 && n.AccessType != 2 {
				tool.ReturnError(w, errors.New(_const.AccessTypeError))
				return
			} else {
				conf.Netconf[i].AccessType = n.AccessType
			}
		}
	}

	err = mongo.DeleteSvcConfById(conf.Id.Hex())
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	err = mongo.SaveSvcConfig(conf)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	w.Write([]byte(conf.Id.Hex()))
	return
}

// UpdateNetPort 更新服务端口数据
// *svc 服务名称
// *namespace 命名空间
// *port 端口数据(base64编码)
// {
//	"name":string,
// 	"img":string,
//	"net":[
//		{NetConfigure}
//	]
// }
func UpdateNetPort(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svc")
	if svc == "" {
		tool.ReturnError(w, errors.New("svc empty!"))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		tool.ReturnError(w, errors.New("namespace empty!"))
		return
	}

	port := r.URL.Query().Get("net")
	if port == "" {
		tool.ReturnError(w, errors.New("port empty!"))
		return
	}

	pb, err := url.QueryUnescape(port)
	//pb, err := base64.StdEncoding.DecodeString(port)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if _const.DEBUG {
		log.Printf("[UpdateNetPort] NetConfigure [%s]\n", pb)
	}

	var cp CPort

	err = json.Unmarshal([]byte(pb), &cp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	isChange, err := container.UpgradeContainerNetByName(cp.Name, svc, nsme, cp.Net)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	if _const.DEBUG {
		log.Printf("[UpdateNetPort] Compare NetConfigure Is change? [%v]!\n", isChange)
	}

	if isChange {
		scf, err := svcconf.GetSvcConfByName(svc, nsme)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		err = svcconf.GenerateNetconifg(scf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	return
}

func GetSvcConfByName(name, ns string) (cf svcconf.SvcConf, err error) {
	conf, err := mongo.GetSvcConfByName(name, ns)
	if err != nil {
		return
	}

	data, err := bson.Marshal(conf)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &cf)
	if err != nil {
		return
	}

	return
}

func CheckSvcConf(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var conf svcconf.SvcConf

	err = json.Unmarshal(data, &conf)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkConf(&conf); err != nil {
		tool.ReturnError(w, err)
		return
	}

	rm := _const.RespMsg{
		Code: 1000,
		Msg:  "SvcConfig Upgrade",
	}

	logrus.WithFields(logrus.Fields{"New Svc Conf": conf}).Info(ModuleName)
	cf, err := mongo.GetSvcConfByName(conf.Name, conf.Namespace)
	if cf == nil {
		err = checkConf(&conf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		conf.Id = bson.NewObjectId()
		if err = mongo.SaveSvcConfig(conf); err != nil {
			tool.ReturnError(w, err)
			return
		}
		rm.Code = 1001
		rm.Msg = "Create New SvcConfig"
	} else {
		nc, err := svcconf.Conver(cf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		svcconf.MergerSvc(nc, &conf)
		err = mongo.DeleteSvcConfById(nc.Id.Hex())
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		conf.Id = nc.Id
		logrus.WithFields(logrus.Fields{"Replace Svc Conf": conf, "Old Svc Conf": nc}).Info(ModuleName)
		if err = mongo.SaveSvcConfig(conf); err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	data, err = json.Marshal(&rm)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Write(data)
	return
}

func AddSvcConfGroup(w http.ResponseWriter, r *http.Request) {
	svcname := r.URL.Query().Get("svcname")
	if svcname == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	clustid := r.URL.Query().Get("clusterid")
	if clustid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	idx := r.URL.Query().Get("idx")
	if idx == "" {
		tool.ReturnError(w, errors.New(_const.IdxNotFound))
		return
	}

	dx, err := strconv.Atoi(idx)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if dx == 0 {
		tool.ReturnError(w, errors.New(_const.IdxVlaueError))
	}

	_, err = mongo.GetSvcConfByName(svcname, namespace)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tool.ReturnError(w, errors.New(_const.SvcNotFound))
			return
		}
		tool.ReturnError(w, err)
		return
	}

	scg, err := mongo.GetSvcConfGroupByName(name, namespace)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			tool.ReturnError(w, err)
			return
		}
	}

	nscg, err := svcconf.Unmarshal(scg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if nscg.Id == "" {
		nscg.Id = bson.NewObjectId()
	}

	nscg.Name = name
	nscg.Namespace = namespace
	nscg.Clusterid = clustid
	sfMap := nscg.SvcGroup
	if len(sfMap) == 0 {
		sfMap = make(map[string]int)
	}

	if sfMap[svcname] > 0 {
		tool.ReturnError(w, errors.New(_const.SvcHasExist))
		return
	} else {
		sfMap[svcname] = dx
	}

	nscg.SvcGroup = sfMap
	mongo.DeleteSvcConfGroup(nscg.Id.Hex())
	if err = mongo.SaveSvcConfGroup(nscg); err != nil {
		tool.ReturnError(w, err)
		return
	}

	return
}

func GetSvcConfGroup(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	name := r.URL.Query().Get("name")
	if _const.DEBUG {
		log.Printf("[GetSvcConfGroup] namespace:[%s] name:[%s] \n", namespace, name)
	}

	w.Header().Set("Content-Type", "application/json")
	if name == "" {
		scg, err := mongo.GetAllSvcConfGroupByNs(namespace)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(&scg)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		if _const.DEBUG {
			log.Printf("[GetSvcConfGroup] Http Response [%s] \n", string(data))
		}

		tool.ReturnResp(w, data)
	} else {
		scg, err := mongo.GetSvcConfGroupByName(name, namespace)
		if err != nil && !tool.IsNotFound(err) {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(&scg)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		if _const.DEBUG {
			log.Printf("[GetSvcConfGroup] Http Response [%s] \n", string(data))
		}

		tool.ReturnResp(w, data)
	}

}

func DeleteSvcConfGroup(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.NameNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	scg, err := mongo.GetSvcConfGroupByName(name, namespace)
	if namespace == "" {
		tool.ReturnError(w, err)
		return
	}

	nscg, err := svcconf.Unmarshal(scg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	err = mongo.DeleteSvcConfGroup(nscg.Id.Hex())
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	svcname := r.URL.Query().Get("svcname")
	if svcname != "" {
		delete(nscg.SvcGroup, svcname)
		if err = mongo.SaveSvcConfGroup(nscg); err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	return

}

func (this *Operation) DeleteSvcConf(msg _const.DestoryMsg) error {
	// 当只需要删除单个服务时(例如整体服务升级，不需要通过蓝绿发布时),此时服务名为 SVC-xxx格式
	// 在数据库中肯定无法找到相对应的记录，因此数据库记录不作为必须条件
	// 当数据库有记录时，删除其名下所有服务
	// 当数据库没有记录时，直接删除当前服务
	var md *metadata.MetaData
	var err error
	needDeleteService := true //在预发布环境中不需要删除服务

	switch msg.Namespace {
	case "proenv":
		fallthrough
	case "release":
		needDeleteService = false
		md, err = metadata.GetMetaDataByRegion("", msg.Namespace)
	case "testenv":
		md, err = metadata.GetMetaDataByRegion("", "testenv")
	default:
		md, err = metadata.GetMetaDataByRegion("")
	}
	//if msg.Namespace == "proenv" {
	//	needDeleteService = false
	//	md, err = metadata.GetMetaDataByRegion("", msg.Namespace)
	//} else {
	//	md, err = metadata.GetMetaDataByRegion("")
	//}

	if err != nil {
		return err
	}
	scf, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Query SvcConf": scf}).Info(ModuleName)
	if scf == nil {
		q := service.Service{
			Pub: public.Public{
				SecretId: md.Sid,
				Region:   md.Region,
			},
			ClusterId:   md.ClusterID,
			Namespace:   msg.Namespace,
			ServiceName: msg.Svcname,
			SecretKey:   md.Skey,
		}

		q.SetDebug(true)

		resp, err := q.DeleteService()
		if err != nil {
			return err
		}
		if resp.Code != 0 {
			if strings.Contains(resp.CodeDesc, "KubeResourceNotFound") {
				logrus.WithFields(logrus.Fields{"DeleteService Error": resp}).Error(ModuleName)
				return nil
			}
			return errors.New(resp.Message)
		}
		return nil
	}

	if scf.SvcName == "" {
		logrus.WithFields(logrus.Fields{"SvcName-Empty": scf.Name, "Not-Delete": true}).Error(ModuleName)
		return nil
	}

	if len(scf.SvcNameBak) > 0 {
		needDeleteService = false
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		Namespace:   scf.Namespace,
		ServiceName: scf.SvcName,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)

	resp, err := q.DeleteService()
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		if strings.Contains(resp.CodeDesc, "KubeResourceNotFound") {
			logrus.WithFields(logrus.Fields{"DeleteService Error": resp}).Error(ModuleName)
			logrus.WithFields(logrus.Fields{"Need Delete Service And Container": needDeleteService}).Info(ModuleName)
			if needDeleteService {
				return scf.DeleteMySelf()
			}
		}
		return errors.New(resp.Message)
	}
	/*删除可能存在的升级服务*/
	for k, _ := range scf.SvcNameBak {
		q := service.Service{
			Pub: public.Public{
				SecretId: md.Sid,
				Region:   md.Region,
			},
			ClusterId:   md.ClusterID,
			Namespace:   scf.Namespace,
			ServiceName: k,
			SecretKey:   md.Skey,
		}

		q.SetDebug(true)

		resp, err := q.DeleteService()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":   err.Error(),
				"svcname": k,
			}).Error("Delete Service Error")
		}

		if resp.Code != 0 {
			logrus.WithFields(logrus.Fields{
				"resp_code": resp.Code,
				"svcname":   k,
				"msg":       resp.Message,
			}).Error("Delete Service Failed")
		}
	}

	logrus.WithFields(logrus.Fields{"Need Delete Service And Container": needDeleteService}).Info(ModuleName)
	if needDeleteService {
		err = container.DeleteAllContaienrUnderSvc(scf.Name, scf.Namespace)
		if err != nil {
			return err
		}
		return scf.DeleteMySelf()
	}

	return nil
}
