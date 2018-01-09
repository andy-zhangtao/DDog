package svcconf

import (
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server/tool"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"strings"
	"strconv"
	"log"
)

// SvcConf 服务配置信息
// 默认情况下Replicas为1
type SvcConf struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"name"`
	Desc      string        `json:"desc"`
	Replicas  int           `json:"replicas"`
	Namespace string        `json:"namespace"`
	Netconf   NetConfigure  `json:"netconf"`
}

// NetConfigure 服务配置信息
// accessType 默认为ClusterIP:
//     0 - ClusterIP
//     1 - LoadBalancer
//     2 - SvcLBTypeInner
// Inport 容器监听端口
// Outport 负载监听端口
// protocol 协议类型 默认为TCP
//     0 - TCP
//     1 - UDP
type NetConfigure struct {
	AccessType int `json:"access_type"`
	InPort     int `json:"in_port"`
	OutPort    int `json:"out_port"`
	Protocol   int `json:"protocol"`
}

// SvcConfGroup 服务群组配置信息
// 作为自己的软服务编排(以业务场景为主,进行的服务编排.不依赖于k8s的服务编排)
type SvcConfGroup struct {
	Id        bson.ObjectId  `json:"id,omitempty" bson:"_id,omitempty"`
	SvcGroup  map[string]int `json:"svc_group"`
	Namespace string         `json:"namespace"`
	Clusterid string         `json:"clusterid"`
	Name      string         `json:"name"`
}

func CreateSvcConf(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var conf SvcConf

	err = json.Unmarshal(data, &conf)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkConf(conf); err != nil {
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

func DeleteSvcConf(w http.ResponseWriter, r *http.Request) {

	//w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		nsme := r.URL.Query().Get("namespace")
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
		err := mongo.DeleteSvcConfByNs(nsme)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	} else {
		err := mongo.DeleteSvcConfById(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}
}
func checkConf(conf SvcConf) error {
	if conf.Name == "" {
		return errors.New(_const.NameNotFound)
	}
	if conf.Namespace == "" {
		return errors.New(_const.NamespaceNotFound)
	}

	conf.Name = strings.Replace(strings.ToLower(conf.Name), " ", "-", -1)
	conf.Namespace = strings.Replace(strings.ToLower(conf.Namespace), " ", "-", -1)
	if conf.Netconf.Protocol != 0 && conf.Netconf.Protocol != 1 {
		return errors.New(_const.LbProtocolError)
	}

	if conf.Netconf.InPort == 0 || conf.Netconf.OutPort == 0 {
		return errors.New(_const.LbPortError)
	}

	if conf.Netconf.AccessType != 0 && conf.Netconf.AccessType != 1 && conf.Netconf.AccessType != 2 {
		return errors.New(_const.AccessTypeError)
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

	conf, err := conver(c)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var nc SvcConf

	err = json.Unmarshal(data, &nc)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if nc.Replicas > 0 {
		conf.Replicas = nc.Replicas
	}

	if nc.Netconf.Protocol != 0 && nc.Netconf.Protocol != 1 {
		tool.ReturnError(w, errors.New(_const.LbProtocolError))
		return
	} else {
		conf.Netconf.Protocol = nc.Netconf.Protocol
	}

	if nc.Netconf.InPort == 0 || nc.Netconf.OutPort == 0 {
		tool.ReturnError(w, errors.New(_const.LbPortError))
		return
	} else if nc.Netconf.InPort > 0 {
		conf.Netconf.InPort = nc.Netconf.InPort
	} else if nc.Netconf.OutPort > 0 {
		conf.Netconf.OutPort = nc.Netconf.OutPort
	}

	if nc.Netconf.AccessType != 0 && nc.Netconf.AccessType != 1 && nc.Netconf.AccessType != 2 {
		tool.ReturnError(w, errors.New(_const.AccessTypeError))
		return
	} else {
		conf.Netconf.AccessType = nc.Netconf.AccessType
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

func conver(conf interface{}) (c *SvcConf, err error) {
	data, err := bson.Marshal(conf)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &c)
	if err != nil {
		return
	}

	return
}

func GetSvcConfByID(id string) (cf SvcConf, err error) {
	conf, err := mongo.GetSvcConfByID(id)
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

func GetSvcConfByName(name, ns string) (cf SvcConf, err error) {
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

	var conf SvcConf

	err = json.Unmarshal(data, &conf)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkConf(conf); err != nil {
		tool.ReturnError(w, err)
		return
	}

	rm := _const.RespMsg{
		Code: 1000,
		Msg:  "SvcConfig Upgrade",
	}

	cf, err := mongo.GetSvcConfByName(conf.Name, conf.Namespace)
	if cf == nil {
		if conf.Replicas == 0 {
			conf.Replicas = 1
		}

		conf.Id = bson.NewObjectId()
		if err = mongo.SaveSvcConfig(conf); err != nil {
			tool.ReturnError(w, err)
			return
		}
		rm.Code = 1001
		rm.Msg = "Create New SvcConfig"
	} else {
		nc, err := conver(cf)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		err = mongo.DeleteSvcConfById(nc.Id.Hex())
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		conf.Id = nc.Id
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

	nscg, err := Unmarshal(scg)
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

	nscg, err := Unmarshal(scg)
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
func Unmarshal(scg interface{}) (nscf SvcConfGroup, err error) {
	if scg == nil {
		return
	}
	data, err := bson.Marshal(scg)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &nscf)
	if err != nil {
		return
	}

	return
}
