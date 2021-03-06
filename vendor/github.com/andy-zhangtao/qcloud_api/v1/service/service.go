package service

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/url"
	"fmt"
	"github.com/andy-zhangtao/qcloud_api/const/v1"
)

var debug = false

const (
	CREATE_SVC           = iota
	UPGRADE_SVC
	DELETE_SVC
	REDEPLOY_SVC
	QUERYSVCINSTANCE
	QUERYSVCINFO
	DELETESVCINSTANCE
	MODIFYSVCINSTANCE
	DESCRIBESERVICEEVENT
)

type Svc struct {
	Pub          public.Public `json:"pub"`
	ClusterId    string        `json:"cluster_id"`
	Namespace    string        `json:"namespace"`
	Allnamespace string        `json:"allnamespace"`
	SecretKey    string        `json:"secret_key"`
	sign         string
}

type Service struct {
	Pub          public.Public  `json:"pub"`
	ClusterId    string         `json:"cluster_id"`
	ServiceName  string         `json:"service_name"`
	ServiceDesc  string         `json:"service_desc"`
	Replicas     int            `json:"replicas"`
	AccessType   string         `json:"access_type"`
	Namespace    string         `json:"namespace"`
	Containers   []Containers   `json:"containers"`
	PortMappings []PortMappings `json:"port_mappings"`
	Strategy     string         `json:"strategy"`
	Instance     []string       `json:"instance"`
	ScaleTo      int            `json:"scale_to"`
	SecretKey    string
	sign         string
}

type Containers struct {
	ContainerName string            `json:"container_name"`
	Image         string            `json:"image"`
	Envs          map[string]string `json:"envs"`
	Command       string            `json:"command"`
	HealthCheck   []HealthCheck     `json:"health_check"`
	Cpu           int               `json:"cpu"`
	CpuLimits     int               `json:"cpuLimits"`
	Memory        int               `json:"memory"`
	MemoryLimits  int               `json:"memoryLimits"`
}

type PortMappings struct {
	LbPort        int    `json:"lb_port"`
	ContainerPort int    `json:"container_port"`
	NodePort      int    `json:"node_port"`
	Protocol      string `json:"protocol"`
}

type SvcData_data_services struct {
	ServiceName     string `json:"servicename"`
	Status          string `json:"status"`
	ServiceIp       string `json:"serviceip"`
	ExternalIp      string `json:"externalip"`
	LbId            string `json:"lbid"`
	LbStatus        string `json:"lbstatus"`
	AccessType      string `json:"accesstype"`
	DesiredReplicas int    `json:"desiredreplicas"`
	CurrentReplicas int    `json:"currentreplicas"`
	CreatedAt       string `json:"createdat"`
	Namespace       string `json:"namespace"`
}

type SvcData_data struct {
	TotalCount  int                     `json:"totalcount"`
	Services    []SvcData_data_services `json:"services"`
	Instance    []Instance              `json:"instances"`
	ServiceInfo ServiceInfo             `json:"service"`
	EventList   []Event_data_eventList  `json:"eventList"`
}

type SvcSMData struct {
	Code     int          `json:"code"`
	Message  string       `json:"message"`
	CodeDesc string       `json:"codedesc"`
	Url      string       `json:"request"`
	Data     SvcData_data `json:"data"`
}

// 实例数据
type Instance_data struct {
	TotalCount int        `json:"totalcount"`
	Instance   []Instance `json:"instanaces"`
}

// 实例数据
type SvcInstance struct {
	Code     int           `json:"code"`
	Message  string        `json:"message"`
	CodeDesc string        `json:"codedesc"`
	Url      string        `json:"request"`
	Data     Instance_data `json:"data"`
}

// 事件数据
type Event_data_eventList struct {
	FirstSeen string `json:"firstseen"`
	LastSeen  string `json:"lastseen"`
	Count     int    `json:"count"`
	Level     string `json:"level"`
	ObjType   string `json:"objtype"`
	ObjName   string `json:"objname"`
	Reason    string `json:"reason"`
	SrcReason string `json:"srcreason"`
	Message   string `json:"message"`
}

func (this Svc) querySampleInfo() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.ClusterId != "" {
		field = append(field, "clusterId")
		req["clusterId"] = this.ClusterId
	}

	if this.Namespace != "" {
		field = append(field, "namespace")
		req["namespace"] = this.Namespace
	}
	an, _ := strconv.Atoi(this.Allnamespace)
	if an != 0 {
		field = append(field, "allnamespace")
		req["allnamespace"] = this.Allnamespace
	}

	return field, req
}

// QueryClusters 查询集群信息
func (this Svc) QuerySampleInfo() (*SvcSMData, error) {
	field, reqmap := this.querySampleInfo()
	pubMap := public.PublicParam("DescribeClusterService", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[获取服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", public.API_URL+reqURL, this.SecretKey, signStr, sign)
	}

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ssmd SvcSMData

	err = json.Unmarshal(data, &ssmd)
	if err != nil {
		return nil, err
	}

	return &ssmd, nil
}

func (this Svc) SetDebug(isDebug bool) {
	debug = isDebug
}

func (this Service) SetDebug(isDebug bool) {
	debug = isDebug
}

func (this Service) CreateNewSerivce() (*SvcSMData, error) {
	// 新建服务,不需要填写升级策略
	this.Strategy = ""
	return this.generateRequest(CREATE_SVC)
}

func (this Service) createSvc() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.ClusterId != "" {
		field = append(field, "clusterId")
		req["clusterId"] = this.ClusterId
	}

	if this.Namespace != "" {
		field = append(field, "namespace")
		req["namespace"] = this.Namespace
	}

	if this.ServiceName != "" {
		field = append(field, "serviceName")
		req["serviceName"] = this.ServiceName
	}

	if this.ServiceDesc != "" {
		field = append(field, "serviceDesc")
		req["serviceDesc"] = this.ServiceDesc
	}

	if this.Replicas > 0 {
		field = append(field, "replicas")
		req["replicas"] = strconv.Itoa(this.Replicas)
	}

	if this.AccessType != "" {
		field = append(field, "accessType")
		req["accessType"] = this.AccessType
	}

	for i, c := range this.Containers {
		if c.ContainerName != "" {
			key := fmt.Sprintf("containers.%d.containerName", i)
			field = append(field, key)
			req[key] = c.ContainerName
		}

		if c.Image != "" {
			key := fmt.Sprintf("containers.%d.image", i)
			field = append(field, key)
			req[key] = c.Image
		}

		n := 0
		for k := range c.Envs {
			key := fmt.Sprintf("containers.%d.envs.%d.name", i, n)
			field = append(field, key)
			req[key] = k
			key = fmt.Sprintf("containers.%d.envs.%d.value", i, n)
			field = append(field, key)
			req[key] = c.Envs[k]
			n++
		}

		if c.Command != "" {
			key := fmt.Sprintf("containers.%d.command", i)
			field = append(field, key)
			req[key] = url.QueryEscape(c.Command)
		}

		for n, hk := range c.HealthCheck {
			if hk.Type != "" {
				key := fmt.Sprintf("containers.%d.healthCheck.%d.type", i, n)
				field = append(field, key)
				req[key] = hk.Type
			}

			key := fmt.Sprintf("containers.%d.healthCheck.%d.healthNum", i, n)
			field = append(field, key)
			req[key] = strconv.Itoa(hk.HealthNum)

			key = fmt.Sprintf("containers.%d.healthCheck.%d.unhealthNum", i, n)
			field = append(field, key)
			req[key] = strconv.Itoa(hk.UnhealthNum)

			key = fmt.Sprintf("containers.%d.healthCheck.%d.intervalTime", i, n)
			field = append(field, key)
			req[key] = strconv.Itoa(hk.IntervalTime)

			key = fmt.Sprintf("containers.%d.healthCheck.%d.timeOut", i, n)
			field = append(field, key)
			req[key] = strconv.Itoa(hk.TimeOut)

			key = fmt.Sprintf("containers.%d.healthCheck.%d.delayTime", i, n)
			field = append(field, key)
			req[key] = strconv.Itoa(hk.DelayTime)

			key = fmt.Sprintf("containers.%d.healthCheck.%d.checkMethod", i, n)
			field = append(field, key)
			req[key] = hk.CheckMethod

			switch hk.CheckMethod {
			case CheckMethodHTTP:
				key = fmt.Sprintf("containers.%d.healthCheck.%d.port", i, n)
				field = append(field, key)
				req[key] = strconv.Itoa(hk.Port)

				key = fmt.Sprintf("containers.%d.healthCheck.%d.protocol", i, n)
				field = append(field, key)
				req[key] = hk.Protocol

				key = fmt.Sprintf("containers.%d.healthCheck.%d.path", i, n)
				field = append(field, key)
				req[key] = hk.Path
			case CheckMethodCmd:
				key = fmt.Sprintf("containers.%d.healthCheck.%d.cmd", i, n)
				field = append(field, key)
				req[key] = hk.Cmd
			case CheckMethodTCP:
				key = fmt.Sprintf("containers.%d.healthCheck.%d.port", i, n)
				field = append(field, key)
				req[key] = strconv.Itoa(hk.Port)
			}

		}
		key := fmt.Sprintf("containers.%d.cpu", i)
		field = append(field, key)
		req[key] = strconv.Itoa(c.Cpu)

		key = fmt.Sprintf("containers.%d.cpuLimits", i)
		field = append(field, key)
		req[key] = strconv.Itoa(c.CpuLimits)

		key = fmt.Sprintf("containers.%d.memory", i)
		field = append(field, key)
		req[key] = strconv.Itoa(c.Memory)

		key = fmt.Sprintf("containers.%d.memoryLimits", i)
		field = append(field, key)
		req[key] = strconv.Itoa(c.MemoryLimits)
	}

	for i, p := range this.PortMappings {
		if p.LbPort > 0 {
			key := fmt.Sprintf("portMappings.%d.lbPort", i)
			field = append(field, key)
			req[key] = strconv.Itoa(p.LbPort)
		}

		if p.ContainerPort > 0 {
			key := fmt.Sprintf("portMappings.%d.containerPort", i)
			field = append(field, key)
			req[key] = strconv.Itoa(p.ContainerPort)
		}

		if p.NodePort > 0 {
			key := fmt.Sprintf("portMappings.%d.nodePort", i)
			field = append(field, key)
			req[key] = strconv.Itoa(p.NodePort)
		}

		if p.Protocol != "" {
			key := fmt.Sprintf("portMappings.%d.protocol", i)
			field = append(field, key)
			req[key] = p.Protocol
		}
	}

	if this.Strategy != "" {
		field = append(field, "strategy")
		req["strategy"] = this.Strategy
	}

	for i, n := range this.Instance {
		key := fmt.Sprintf("instances.%d", i)
		field = append(field, key)
		req[key] = n
	}

	field = append(field, "scaleTo")
	req["scaleTo"] = strconv.Itoa(this.ScaleTo)
	return field, req
}

func (this Service) UpgradeService() (*SvcSMData, error) {
	//field, reqmap := this.createSvc()
	//pubMap := public.PublicParam("ModifyClusterService", this.Pub.Region, this.Pub.SecretId)
	//this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	//signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	//sign := public.GenerateSignature(this.SecretKey, signStr)
	//reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)
	//
	//if debug {
	//	log.Printf("[升级服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", public.API_URL+reqURL, this.SecretKey, signStr, sign)
	//}
	//
	//resp, err := http.Get(public.API_URL + reqURL)
	//if err != nil {
	//	return nil, err
	//}
	//
	//data, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var ssmd SvcSMData
	//
	//err = json.Unmarshal(data, &ssmd)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &ssmd, nil
	return this.generateRequest(1)
}

func (this Service) DeleteService() (*SvcSMData, error) {
	return this.generateRequest(DELETE_SVC)
}

func (this Service) RedeployService() (*SvcSMData, error) {
	return this.generateRequest(REDEPLOY_SVC)
}

func (this Service) QueryInstance() (*SvcSMData, error) {
	return this.generateRequest(QUERYSVCINSTANCE)
}

func (this Service) QuerySvcInfo() (*SvcSMData, error) {
	return this.generateRequest(QUERYSVCINFO)
}

func (this Service) DestoryInstance() (*SvcSMData, error) {
	return this.generateRequest(DELETESVCINSTANCE)
}

func (this Service) ModeifyInstance() (*SvcSMData, error) {
	return this.generateRequest(MODIFYSVCINSTANCE)
}

func (this Service) DescribeServiceEvent() (*SvcSMData, error) {
	return this.generateRequest(DESCRIBESERVICEEVENT)
}

// generateRequest 生成操作请求
// 每个请求中都存在公共部分,因此在这里只需要处理特殊操作对应的数据即可
// 0 - 创建服务
// 1 - 升级服务
// 2 - 删除服务
// 3 - 重新部署服务
// 4 - 查询实例
// 5 - 查询服务详情
// 6 - 删除实例
// 7 - 修改实例个数
func (this Service) generateRequest(kind int) (*SvcSMData, error) {
	var svcKind string
	var debugStr string
	switch kind {
	case CREATE_SVC:
		//	创建
		svcKind = "CreateClusterService"
		debugStr = "创建"
	case UPGRADE_SVC:
		//	升级
		svcKind = "ModifyClusterService"
		debugStr = "升级"
	case DELETE_SVC:
		//	删除
		svcKind = "DeleteClusterService"
		debugStr = "删除"
	case REDEPLOY_SVC:
		//	重新部署
		svcKind = "RedeployClusterService"
		debugStr = "重新部署"
	case QUERYSVCINSTANCE:
		//	查询服务实例状态
		svcKind = "DescribeServiceInstance"
		debugStr = "查询服务实例"
	case QUERYSVCINFO:
		//	查询服务详情
		svcKind = "DescribeClusterServiceInfo"
		debugStr = "查询服务详情"
	case DELETESVCINSTANCE:
		//	删除实例
		svcKind = "DeleteInstances"
		debugStr = "删除实例"
	case MODIFYSVCINSTANCE:
		//	修改服务实例
		svcKind = "ModifyServiceReplicas"
		debugStr = "修改实例数量"
	case DESCRIBESERVICEEVENT:
		//	获取服务实例集群事件
		svcKind = "DescribeServiceEvent"
		debugStr = "获取服务事件"
	}

	field, reqmap := this.createSvc()
	pubMap := public.PublicParam(svcKind, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[%s 服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", debugStr, public.API_URL+reqURL, this.SecretKey, signStr, sign)
	}

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ssmd SvcSMData

	err = json.Unmarshal(data, &ssmd)
	if err != nil {
		return nil, err
	}

	if ssmd.Code != 0 {
		ssmd.Url = public.API_URL + reqURL
	}

	return &ssmd, nil

}
