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
	"math"
	"github.com/sirupsen/logrus"
)

// SvcConf 服务配置信息
// 默认情况下Replicas为2
type SvcConf struct {
	Id            bson.ObjectId            `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string                   `json:"name"`
	Desc          string                   `json:"desc"`
	SvcName       string                   `json:"svc_name"`     // 在K8s中的服务名
	SvcNameBak    map[string]LoadBlance    `json:"svc_name_bak"` // 升级过程中的备份服务名
	Replicas      int                      `json:"replicas"`
	Namespace     string                   `json:"namespace"`
	Netconf       []container.NetConfigure `json:"netconf"`
	Status        int                      `json:"status"` // 0 - 处理成功 1 - 准备解析网络配置 2 - 开始解析网络配置 3 - 网络解析配置失败
	Msg           string                   `json:"msg"`
	Deploy        int                      `json:"deploy"` // 0 - 未部署 1 - 部署成功 2 - 部署中 3 - 蓝绿部署中 4 - 部署失败 5-滚动部署部分完成 6 - 数据同步 7-回滚中 8-升级确认中
	Instance      []SvcInstance            `json:"instance"`
	LbConfig      LoadBlance               `json:"lb_config"`
	BackID        string                   `json:"back_id"`
	BackContainer []container.Container    `json:"back_container,omitempty"`
}

// SvcInstance 服务实例信息
// 用于蓝绿发布/金丝雀发布
type SvcInstance struct {
	Name   string `json:"name"`   //服务名称. 此名称对应的是K8s中的服务实例名
	Status int    `json:"status"` //服务当前状态. 0 - 未部署 1 - 部署成功 2 - 部署中 3 - 部署失败 4 - 滚动部署中
	Msg    string `json:"msg"`
}

// LoadBlance 负载均衡数据
type LoadBlance struct {
	IP   string `json:"ip"`
	Port []int  `json:"port"`
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

// ToString 格式化输出结构体内容
func (scf *SvcConf) ToString() (out string) {
	out += fmt.Sprintf("\n********[SvcConf]********\n")
	out += fmt.Sprintf("ID:[%s]\n", scf.Id.Hex())
	out += fmt.Sprintf("Name:[%s]\n", scf.Name)
	out += fmt.Sprintf("Desc:[%s]\n", scf.Desc)
	out += fmt.Sprintf("SvcName:[%s]\n", scf.SvcName)
	out += fmt.Sprintf("SvcNameBak:[%v]\n", scf.SvcNameBak)
	out += fmt.Sprintf("Replicas:[%d]\n", scf.Replicas)
	out += fmt.Sprintf("Namespace:[%s]\n", scf.Namespace)
	for _, n := range scf.Netconf {
		out += fmt.Sprintf("\t\t\t\tNetConf:[%s]\n", n.ToString())
	}
	out += fmt.Sprintf("Status:[%d]\n", scf.Status)
	out += fmt.Sprintf("Msg:[%s]\n", scf.Msg)
	out += fmt.Sprintf("Deploy:[%d]\n", scf.Deploy)
	for _, s := range scf.Instance {
		out += fmt.Sprintf("\t\t\t\tInstance:[%s]\n", s.ToString())
	}
	out += fmt.Sprintf("LbConfig:[%s]\n", scf.LbConfig.ToString())
	out += fmt.Sprintf("BackID:[%s]\n", scf.BackID)
	for _, b := range scf.BackContainer {
		out += fmt.Sprintf("\t\t\t\tBackContainer:[%s]\n", b.ToString())
	}
	return
}

// 将当前结构体的数据复制给tcp
func (scf *SvcConf) Copy(tcp *SvcConf) {
	tcp.Id = scf.Id
	tcp.Name = scf.Name
	tcp.Desc = scf.Desc
	tcp.SvcName = scf.SvcName
	tcp.SvcNameBak = scf.SvcNameBak
	tcp.Replicas = scf.Replicas
	tcp.Namespace = scf.Namespace
	tcp.Netconf = scf.Netconf
	tcp.Status = scf.Status
	tcp.Msg = scf.Msg
	tcp.Deploy = scf.Deploy
	tcp.Instance = scf.Instance
	tcp.LbConfig = scf.LbConfig
	tcp.BackID = scf.BackID
	tcp.BackContainer = scf.BackContainer
}

func (sis *SvcInstance) ToString() (out string) {
	out = ""
	out += fmt.Sprintf("\n********[SvcInstance]********\n")
	out += fmt.Sprintf("Name:[%s]\n", sis.Name)
	out += fmt.Sprintf("Status:[%d]\n", sis.Status)
	out += fmt.Sprintf("Msg:[%s]\n", sis.Msg)
	return
}

func (lb *LoadBlance) ToString() (out string) {
	out = ""
	out += fmt.Sprintf("\n********[LoadBlance]********\n")
	out += fmt.Sprintf("IP:[%s]\n", lb.IP)
	for _, i := range lb.Port {
		out += fmt.Sprintf("\t\t\t\tPort:[%d]\n", i)
	}
	return
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

func GetSvcConfByDeployStatus(deploy int) (scf *SvcConf, err error) {
	sv, err := mongo.GetSvcConfByDeployStatus(deploy)
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

func GetSvcConfByNamespace(namespace string) (scf []SvcConf, err error) {
	svcs, err := mongo.GetSvcConfNs(namespace)
	if err != nil {
		if !tool.IsNotFound(err) {
			return nil, err
		}
		return nil, nil
	}

	for _, s := range svcs {
		st, err := Conver(s)
		if err != nil {
			return nil, err
		}

		scf = append(scf, *st)
	}

	return
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
		log.Printf("[SaveSvcConf mongo.go] Save SvcConf [%v]\n", scf)
	}
	return mongo.SaveSvcConfig(scf)
}

func UpdateSvcConf(scf *SvcConf) error {
	//if _const.DEBUG {
	//	log.Printf("[UpdateSvcConf mongo.go] DeleteSvcConfById [%s]\n", scf.Id.Hex())
	//}
	logrus.WithFields(logrus.Fields{"svc_conf_id": scf.Id.Hex(), "data": scf}).Info("DeleteSvcConfById")
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

// CountInstances 计算本次需要升级实例名称
// 返回挑选出来需要升级的实例名称,同时返回剩余可以用于升级的实例个数
// 在更新Service—Config状态时需要同时知道哪些Instance是本次升级的，哪些Instance是本次挑选时落选的,因此把落选Instance一并返回
func (sc *SvcConf) CountInstances(scope float64) ([]string, []string, int) {

	/*先检索当前还没有升级的实例个数*/
	var instances []SvcInstance
	for _, is := range sc.Instance {
		if is.Status != 4 {
			instances = append(instances, is)
		}
	}

	var name []string
	var leftName []string

	maxNumber := math.Ceil(scope * float64(len(instances)))
	log.Printf("[CountInstance] Ready Roll Up [%v] Service, In Face I Can Support [%v] Services \n", scope, maxNumber)
	number := int(maxNumber)
	if number > 0 {
		for i := 0; i < number; i ++ {
			name = append(name, instances[i].Name)
		}
	}

	for i := number; i < len(instances); i ++ {
		leftName = append(leftName, instances[i].Name)
	}

	return name, leftName, len(instances) - number
}

// BackupSvcConf 备份服务配置
func (sc *SvcConf) BackupSvcConf() error {
	nsc := SvcConf{
		Name:      sc.Name + "_bak",
		Replicas:  sc.Replicas,
		Namespace: sc.Namespace,
		Netconf:   sc.Netconf,
		LbConfig:  sc.LbConfig,
	}

	cons, err := container.GetAllContainersBySvc(sc.Name, sc.Namespace)
	if err != nil {
		return err
	}
	nsc.Id = bson.NewObjectId()
	nsc.BackContainer = cons
	sc.BackID = nsc.Id.Hex()

	err = UpdateSvcConf(sc)
	if err != nil {
		return err
	}

	return SaveSvcConf(&nsc)
}

// GetBackSvcConf 取回备份配置
func (sc *SvcConf) GetBackSvcConf() (bsc *SvcConf, err error) {
	return GetSvcConfByID(sc.BackID)
}

func (sc *SvcConf) DeleteMySelf() (err error) {

	if sc.BackID != "" {
		/*删除备份*/
		mongo.DeleteSvcConfById(sc.BackID)
	}

	logrus.WithFields(logrus.Fields{"Delete SvcConf ID": sc.Id.Hex()}).Info("DeleteMySelf")
	err = mongo.DeleteSvcConfById(sc.Id.Hex())
	return
}

func MergerSvc(oldSvc, newSvc *SvcConf) {
	if newSvc.Name == "" {
		newSvc.Name = oldSvc.Name
	}

	if newSvc.Desc == "" {
		newSvc.Desc = oldSvc.Desc
	}

	if newSvc.SvcName == "" {
		newSvc.SvcName = oldSvc.SvcName
	}

	if len(newSvc.SvcNameBak) == 0 {
		newSvc.SvcNameBak = oldSvc.SvcNameBak
	}

	if newSvc.Replicas == 0 {
		newSvc.Replicas = oldSvc.Replicas
	}

	if newSvc.Namespace == "" {
		newSvc.Namespace = oldSvc.Namespace
	}

	if len(newSvc.Netconf) == 0 {
		newSvc.Netconf = oldSvc.Netconf
	}

	if len(newSvc.Instance) == 0 {
		newSvc.Instance = oldSvc.Instance
	}

	if newSvc.LbConfig.IP == "" {
		newSvc.LbConfig = oldSvc.LbConfig
	}

	if newSvc.BackID == "" {
		newSvc.BackID = oldSvc.BackID
	}

	if len(newSvc.BackContainer) == 0 {
		newSvc.BackContainer = oldSvc.BackContainer
	}
}
