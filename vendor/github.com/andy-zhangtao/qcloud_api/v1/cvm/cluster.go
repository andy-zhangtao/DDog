package cvm

import (
	"strconv"
	"github.com/andy-zhangtao/qcloud_api/v1/public"

	"log"

	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
)

var debug = false

type Cluster struct {
	Pub        public.Public `json:"pub"`
	Cid        string        `json:"cid"`
	Cname      string        `json:"cname"`
	Status     string        `json:"status"`
	OrderField string        `json:"order_field"`
	OrderType  string        `json:"order_type"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	SecretKey  string        `json:"secret_key"`
	Namespace  string        `json:"namespace"`
	sign       string
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterNode_data_nodes struct {
	InstanceId           string `json:"instanceid"`
	InstanceName         string `json:"instancename"`
	InstanceType         string `json:"instancetype"`
	ZoneId               int    `json:"zoneid"`
	WanIp                string `json:"wanip"`
	LanIp                string `json:"lanip"`
	Cpu                  int    `json:"cpu"`
	Mem                  int    `json:"mem"`
	KernelVersion        string `json:"kernelversion"`
	OsImage              string `json:"osimage"`
	PodCidr              string `json:"podcidr"`
	IsNormal             int    `json:"isnormal"`
	AbnormalReason       string `json:"abnormalreason"`
	CvmState             int    `json:"cvmstate"`
	CvmPayMode           int    `json:"cvmpaymode"`
	NetworkPayMode       int    `json:"networkpaymode"`
	CreatedAt            string `json:"createdat"`
	InstanceCreateTime   string `json:"instancecreatetime"`
	InstanceDeadlineTime string `json:"instancedeadlinetime"`
	Unschedulable        bool   `json:"unschedulable"`
	Zone                 string `json:"zone"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterNode_data struct {
	TotalCount int                      `json:"totalcount"`
	Nodes      []ClusterNode_data_nodes `json:"nodes"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterNode struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	CodeDesc string           `json:"codedesc"`
	Data     ClusterNode_data `json:"data"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterInfo_data_clusters struct {
	ClusterId        string `json:"clusterid"`
	ClusterName      string `json:"clustername"`
	Description      string `json:"description"`
	Status           string `json:"status"`
	UnVpcId          string `json:"unvpcid"`
	VpcId            int    `json:"vpcid"`
	ClusterCIDR      string `json:"clustercidr"`
	CreatedAt        string `json:"createdat"`
	UpdatedAt        string `json:"updatedat"`
	NodeStatus       string `json:"nodestatus"`
	NodeNum          int    `json:"nodenum"`
	Os               string `json:"os"`
	TotalCpu         int    `json:"totalcpu"`
	TotalMem         int    `json:"totalmem"`
	RegionId         int    `json:"regionid"`
	K8sVersion       string `json:"k8sversion"`
	OpenHttps        int    `json:"openhttps"`
	MasterLbSubnetId string `json:"masterlbsubnetid"`
	ProjectId        int    `json:"projectid"`
	Region           string `json:"region"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterInfo_data struct {
	TotalCount int                         `json:"totalcount"`
	Clusters   []ClusterInfo_data_clusters `json:"clusters"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type ClusterInfo struct {
	Code     int              `json:"code"`
	Message  string           `json:"message"`
	CodeDesc string           `json:"codedesc"`
	Data     ClusterInfo_data `json:"data"`
}

// queryCluster 查询集群数据API
func (this Cluster) queryCluster() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.Cid != "" {
		field = append(field, "clusterIds.n")
		req["clusterIds.n"] = this.Cid
	}

	if this.Cname != "" {
		field = append(field, "clusterName")
		req["clusterName"] = this.Cname
	}

	if this.Status != "" {
		field = append(field, "status")
		req["status"] = this.Status
	}

	if this.OrderField != "" {
		field = append(field, "orderField")
		req["orderField"] = this.OrderField
	}

	if this.OrderType != "" {
		field = append(field, "orderType")
		req["orderType"] = this.OrderType
	}

	if this.Offset > 0 {
		field = append(field, "offset")
		req["offset"] = strconv.Itoa(this.Offset)
	}

	if this.Limit > 0 {
		field = append(field, "limit")
		req["limit"] = strconv.Itoa(this.Limit)
	}

	return field, req
}

func (this Cluster) queryClusterNode() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	if this.Cid != "" {
		field = append(field, "clusterId")
		req["clusterId"] = this.Cid
	}

	if this.Offset > 0 {
		field = append(field, "offset")
		req["offset"] = strconv.Itoa(this.Offset)
	}

	if this.Limit > 0 {
		field = append(field, "limit")
		req["limit"] = strconv.Itoa(this.Limit)
	}

	if this.Namespace != "" {
		field = append(field, "namespace")
		req["namespace"] = this.Namespace
	}
	return field, req
}

// QueryClusters 查询集群信息
func (this Cluster) QueryClusters() (*ClusterInfo, error) {

	if this.SecretKey == "" {
		return nil, errors.New("SecretKey Can not be empty!")
	}

	field, reqmap := this.queryCluster()
	pubMap := public.PublicParam("DescribeCluster", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + sign

	if debug {
		log.Println(public.API_URL + reqURL)
		log.Println(this.SecretKey)
		log.Println(signStr)
		log.Println(sign)
	}
	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cn ClusterInfo

	err = json.Unmarshal(data, &cn)
	if err != nil {
		return nil, err
	}

	return &cn, nil
}

func (this Cluster) QueryClusterNodes() (*ClusterNode, error) {
	field, reqmap := this.queryClusterNode()
	pubMap := public.PublicParam("DescribeClusterInstances", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + sign

	if debug {
		log.Println(public.API_URL + reqURL)
	}
	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cn ClusterNode

	err = json.Unmarshal(data, &cn)
	if err != nil {
		return nil, err
	}

	return &cn, nil
}

func (this Cluster) SetDebug(isDebug bool) {
	debug = isDebug
}
