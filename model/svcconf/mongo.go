package svcconf

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/DDog/model/container"
	"log"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
)

// SvcConf 服务配置信息
// 默认情况下Replicas为1
type SvcConf struct {
	Id        bson.ObjectId            `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                   `json:"name"`
	Desc      string                   `json:"desc"`
	Replicas  int                      `json:"replicas"`
	Namespace string                   `json:"namespace"`
	Netconf   []container.NetConfigure `json:"netconf"`
	Status    int                      `json:"status"` // 0 - 处理成功 1 - 准备解析网络配置 2 - 开始解析网络配置 3 - 网络解析配置失败
	Msg       string                   `json:"msg"`
	Deploy    int                      `json:"deploy"` // 0 - 未部署 1 - 部署成功 2 - 部署中 3 - 蓝绿部署中 4 - 部署失败
	Instance  []SvcInstance            `json:"instance"`
}

// SvcInstance 服务实例信息
// 用于蓝绿发布/金丝雀发布
type SvcInstance struct {
	Name   string `json:"name"`   //服务名称. 此名称对应的是K8s中的服务实例名
	Status int    `json:"status"` //服务当前状态. 0 - 未部署 1 - 部署成功 2 - 部署中 3 - 部署失败
	Msg    string `json:"msg"`
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

func Conver(conf interface{}) (c *SvcConf, err error) {
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

func GetSvcConfByName(svcname, namespace string) (scf *SvcConf, err error) {
	sv, err := mongo.GetSvcConfByName(svcname, namespace)
	if err != nil {
		if !tool.IsNotFound(err) {
			return nil, err
		}
		return nil, nil
	}

	nscf, err := Conver(sv)
	if err != nil {
		return nil, err
	}

	return nscf, nil
}

func GetSvcConfByID(id string) (*SvcConf, error) {
	conf, err := mongo.GetSvcConfByID(id)
	if err != nil {
		return nil, err
	}

	data, err := bson.Marshal(conf)
	if err != nil {
		return nil, err
	}

	var cf SvcConf
	err = bson.Unmarshal(data, &cf)
	if err != nil {
		return nil, err
	}

	return &cf, nil
}

func SaveSvcConf(scf *SvcConf) error {
	if _const.DEBUG {
		log.Printf("[SaveSvcConf] Save SvcConf [%v]\n", scf)
	}
	return mongo.SaveSvcConfig(scf)
}

func UpdateSvcConf(scf *SvcConf) error {
	if _const.DEBUG {
		log.Printf("[UpdateSvcConf] DeleteSvcConfById [%s]\n", scf.Id.Hex())
	}
	err := mongo.DeleteSvcConfById(scf.Id.Hex())
	if err != nil {
		return err
	}

	return SaveSvcConf(scf)
}

// GenerateNetconifg 重建服务的网络配置信息
func GenerateNetconifg(scf *SvcConf) (err error) {

	cons, err := container.GetAllContainersBySvc(scf.Name, scf.Namespace)
	if err != nil {
		return
	}

	if len(cons) == 0 {
		return errors.New(fmt.Sprintf("This SVC contains 0 container. Name:[%s]Namespace:[%s]", scf.Name, scf.Namespace))
	}
	var net []container.NetConfigure

	pm := make(map[int]int)
	for _, cn := range cons {
		for _, n := range cn.Net {
			if pm[n.InPort] > 0 {
				pm[n.InPort] += 1
			} else {
				pm[n.InPort] = n.InPort
			}

			net = append(net, container.NetConfigure{
				AccessType: n.AccessType,
				InPort:     n.InPort,
				OutPort:    pm[n.InPort],
				Protocol:   n.Protocol,
			})
		}
	}

	scf.Netconf = net
	log.Printf("[GenerateNetconifg] SvcConf [%v]\n", scf)
	err = UpdateSvcConf(scf)
	return
}
