package container

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/tool"
	"log"
	"github.com/andy-zhangtao/DDog/const"
	"fmt"
	"github.com/Sirupsen/logrus"
)

type Container struct {
	ID   bson.ObjectId     `json:"id,omitempty" bson:"_id,omitempty"`
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
	Svc  string            `json:"svc"`
	Nsme string            `json:"namespace"`
	Idx  int               `json:"idx"`
	Net  []NetConfigure    `json:"net"`
	Port []int             `json:"port"`
}

func (c *Container) ToString() (out string) {
	out = ""
	out += fmt.Sprintf("\n********[Container]********\n")
	out += fmt.Sprintf("ID:[%s]\n", c.ID.Hex())
	out += fmt.Sprintf("Name:[%s]\n", c.Name)
	out += fmt.Sprintf("Cmd:[%v]\n", c.Cmd)
	out += fmt.Sprintf("Env:[%v]\n", c.Env)
	out += fmt.Sprintf("Svc:[%s]\n", c.Svc)
	out += fmt.Sprintf("Nsme:[%s]\n", c.Nsme)
	out += fmt.Sprintf("Idx:[%d]\n", c.Idx)
	for _, i := range c.Net {
		out += fmt.Sprintf("Net:[%s]\n", i.ToString())
	}
	out += fmt.Sprintf("Port:[%v]\n", c.Port)
	return
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

func (net *NetConfigure) ToString() (out string) {
	out = ""
	out += fmt.Sprintf("\n********[NetConfigure]********\n")
	out += fmt.Sprintf("AccessType:[%d]\n", net.AccessType)
	out += fmt.Sprintf("InPort:[%d]\n", net.InPort)
	out += fmt.Sprintf("OutPort:[%d]\n", net.OutPort)
	out += fmt.Sprintf("Protocol:[%d]\n", net.Protocol)
	return
}

// GetContainerByName 根据名称获取容器配置信息
// conname 容器名称
// svcname 服务名称
// namespace 命名空间名称
func GetContainerByName(conname, svcname, namespace string) (con *Container, err error) {
	c, err := mongo.GetContaienrByName(conname, svcname, namespace)
	if err != nil {
		if !tool.IsNotFound(err) {
			return
		}
	}

	if _const.DEBUG {
		log.Printf("[GetContainerByName] Get Container From Mongo [%v]\n", c)
	}

	if tool.IsNotFound(err) {
		err = nil
		return
	}

	err = nil
	con, err = unmarshal(c)
	return
}

func SaveContainer(con *Container) (err error) {
	if con.ID.Hex() == "" {
		con.ID = bson.NewObjectId()
	}

	if err = mongo.SaveContainer(con); err != nil {
		return
	}
	return
}

func unmarshal(icon interface{}) (con *Container, err error) {
	if icon == nil {
		return
	}
	data, err := bson.Marshal(icon)
	if err != nil {
		return
	}

	var c Container
	err = bson.Unmarshal(data, &c)
	if err != nil {
		return
	}

	con = &c
	return
}

// DeleteContainerByName 根据名称删除容器配置信息
// conname 容器名称
// svcname 服务名称
// namespace 命名空间名称
func DeleteContainerByName(conname, svcname, namespace string) (err error) {
	c, err := mongo.GetContaienrByName(conname, svcname, namespace)
	if err != nil {
		if !tool.IsNotFound(err) {
			return
		}
	}

	if tool.IsNotFound(err) {
		err = nil
		return
	}

	con, err := unmarshal(c)
	if err != nil {
		return err
	}

	err = mongo.DeleteContainerById(con.ID.Hex())
	return
}

// UpgradeContaienrByName 升级容器配置信息
// con 容器配置指针
func UpgradeContaienrByName(con *Container) (err error) {
	err = DeleteContainerByName(con.Name, con.Svc, con.Nsme)
	if err != nil {
		return
	}

	err = SaveContainer(con)
	return
}

// DeleteAllContaienrUnderSvc 删除指定服务下面的所有容器
func DeleteAllContaienrUnderSvc(svcname, namespace string) (err error) {
	err = mongo.DeleteAllContainer(svcname, namespace)
	return
}

// UpgradeContainerNetByName 更新容器的网络配置信息
// conname 容器名称
// svcname 服务名称
// namespace 命名空间名称
// net 网络配置数据
// 如果网络数据发生变化，则isChange为true，反之为false
func UpgradeContainerNetByName(conname, svcname, namespace string, net []NetConfigure) (isChange bool, err error) {
	if _const.DEBUG {
		log.Printf("[UpgradeContainerNetByName] Compare NetConfigure name:[%s] svc:[%s] namespace:[%s] \n", conname, svcname, namespace)
	}

	con, err := GetContainerByName(conname, svcname, namespace)
	if err != nil {
		return
	}

	if _const.DEBUG {
		log.Printf("[UpgradeContainerNetByName] Readey to Compare Old:[%v] New:[%v]\n", con, net)
	}

	compare := func(net1, net2 []NetConfigure) bool {
		for i, n := range net2 {
			if net1[i].Protocol != n.Protocol {
				return true
			}

			if net1[i].InPort != n.InPort {
				return true
			}

			if net1[i].OutPort != n.OutPort {
				return true
			}

			if net1[i].AccessType != n.AccessType {
				return true
			}
		}

		return false
	}

	logrus.WithFields(logrus.Fields{"len(con.Net)": len(con.Net), "len(net)": len(net)}).Info("UpgradeContainerNetByName Compare Net Array Length")

	if len(con.Net) == 0 || len(net) > len(con.Net) {
		isChange = true
	} else {
		isChange = compare(con.Net, net)
	}

	if isChange {
		con.Net = net
		err = UpgradeContaienrByName(con)
	}

	return
}

// GetAllContainersBySvc 获取指定服务下面所有的容器数据
// svcname 服务名称
// namespace 命名空间名称
func GetAllContainersBySvc(svc, namespace string) (cons []Container, err error) {
	cns, err := mongo.GetContaienrBySvc(svc, namespace)
	if err != nil {
		return
	}

	for _, c := range cns {
		tc, _ := unmarshal(c)
		cons = append(cons, *tc)
	}

	return
}
