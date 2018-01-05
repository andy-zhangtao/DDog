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

type Svc struct {
	Pub          public.Public `json:"pub"`
	ClusterId    string        `json:"cluster_id"`
	Namespace    string        `json:"namespace"`
	Allnamespace string        `json:"allnamespace"`
	SecretKey    string        `json:"secret_key"`
	sign         string
}

type Service struct {
	Pub          public.Public `json:"pub"`
	ClusterId    string        `json:"cluster_id"`
	ServiceName  string        `json:"service_name"`
	ServiceDesc  string        `json:"service_desc"`
	Replicas     int           `json:"replicas"`
	AccessType   string        `json:"access_type"`
	Namespace    string        `json:"namespace"`
	Containers   []Containers  `json:"containers"`
	PortMappings PortMappings  `json:"port_mappings"`
	Strategy     string        `json:"strategy"`
	SecretKey    string
	sign         string
}

type Containers struct {
	ContainerName string            `json:"container_name"`
	Image         string            `json:"image"`
	Envs          map[string]string `json:"envs"`
	Command       string            `json:"command"`
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
	TotalCount int                     `json:"totalcount"`
	Services   []SvcData_data_services `json:"services"`
}
type SvcSMData struct {
	Code     int          `json:"code"`
	Message  string       `json:"message"`
	CodeDesc string       `json:"codedesc"`
	Data     SvcData_data `json:"data"`
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
	field, reqmap := this.createSvc()
	pubMap := public.PublicParam("CreateClusterService", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[创建服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", public.API_URL+reqURL, this.SecretKey, signStr, sign)
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
			req[key] = url.QueryEscape(k)
			key = fmt.Sprintf("containers.%d.envs.%d.value", i, n)
			field = append(field, key)
			req[key] = url.QueryEscape(c.Envs[k])
			n++
		}

		if c.Command != "" {
			key := fmt.Sprintf("containers.%d.command", i)
			field = append(field, key)
			req[key] = url.QueryEscape(c.Command)
		}
	}

	if this.PortMappings.LbPort > 0 {
		field = append(field, "portMappings.0.lbPort")
		req["portMappings.0.lbPort"] = strconv.Itoa(this.PortMappings.LbPort)
	}

	if this.PortMappings.ContainerPort > 0 {
		field = append(field, "portMappings.0.containerPort")
		req["portMappings.0.containerPort"] = strconv.Itoa(this.PortMappings.ContainerPort)
	}

	if this.PortMappings.NodePort > 0 {
		field = append(field, "portMappings.0.nodePort")
		req["portMappings.0.nodePort"] = strconv.Itoa(this.PortMappings.NodePort)
	}

	if this.PortMappings.Protocol != "" {
		field = append(field, "portMappings.0.protocol")
		req["portMappings.0.protocol"] = this.PortMappings.Protocol
	}

	if this.Strategy != "" {
		field = append(field, "strategy")
		req["strategy"] = this.Strategy
	}
	return field, req
}

func (this Service) UpgradeService() (*SvcSMData, error) {
	field, reqmap := this.createSvc()
	pubMap := public.PublicParam("ModifyClusterService", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[升级服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", public.API_URL+reqURL, this.SecretKey, signStr, sign)
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

func (this Service) DeleteService() (*SvcSMData, error) {
	return this.generateRequest(2)
}

func (this Service) generateRequest(kind int) (*SvcSMData, error) {
	var svcKind string
	var debugStr string
	switch kind {
	case 0:
		//	创建
		svcKind = "CreateClusterService"
		debugStr = "创建"
	case 1:
		//	升级
		svcKind = "ModifyClusterService"
		debugStr = "升级"
	case 2:
		//	删除
		svcKind = "DeleteClusterService"
		debugStr = "删除"
	}

	field, reqmap := this.createSvc()
	pubMap := public.PublicParam(svcKind, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[%s服务信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]\n", debugStr, public.API_URL+reqURL, this.SecretKey, signStr, sign)
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
