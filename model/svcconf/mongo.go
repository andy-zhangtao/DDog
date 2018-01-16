package svcconf

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/tool"
)

// SvcConf 服务配置信息
// 默认情况下Replicas为1
type SvcConf struct {
	Id        bson.ObjectId  `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string         `json:"name"`
	Desc      string         `json:"desc"`
	Replicas  int            `json:"replicas"`
	Namespace string         `json:"namespace"`
	Netconf   []NetConfigure `json:"netconf"`
	Status    int            `json:"status"` // 0 - 处理成功 1 - 准备解析网络配置 2 - 开始解析网络配置 3 - 网络解析配置失败
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
	return mongo.SaveSvcConfig(scf)
}

func UpdateSvcConf(scf *SvcConf) error {
	err := mongo.DeleteSvcConfById(scf.Id.Hex())
	if err != nil{
		return err
	}
	return SaveSvcConf(scf)
}
