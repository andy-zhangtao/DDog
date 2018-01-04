package namespace

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"log"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"github.com/andy-zhangtao/qcloud_api/const"
	_constv1 "github.com/andy-zhangtao/qcloud_api/const/v1"
	"fmt"
	"strings"
)

var debug = false

type NSpace struct {
	Pub       public.Public `json:"pub"`
	ClusterId string        `json:"cluster_id"`
	SecretKey string        `json:"secret_key"`
	Name      string        `json:"name"`
	Desc      string        `json:"desc"`
	Rmname    []string      `json:"rmname"`
	sign      string
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type NSInfo_data_namespaces struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdat"`
	ClusterID   string `json:"cluster_id"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type NSInfo_data struct {
	TotalCount int                      `json:"totalcount"`
	Namespaces []NSInfo_data_namespaces `json:"namespaces"`
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type NSInfo struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	CodeDesc string      `json:"codedesc"`
	Data     NSInfo_data `json:"data"`
}

func (this NSpace) SetDebug(isDebug bool) () {
	debug = isDebug
}

func (this NSpace) QueryNSInfo() (*NSInfo, error) {
	if this.ClusterId == "" {
		return nil, errors.New("ClusterId Can not be empty!")
	}

	field, reqmap := this.queryNSInfo()
	pubMap := public.PublicParam("DescribeClusterNameSpaces", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GETccs.api.qcloud.com/v2/index.php?" + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	if debug {
		log.Printf("[获取命名空间信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]", public.API_URL+reqURL, this.SecretKey, signStr, sign)
	}

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ns NSInfo

	err = json.Unmarshal(data, &ns)
	if err != nil {
		return nil, err
	}

	return &ns, nil
}

func (this NSpace) queryNSInfo() ([]string, map[string]string) {
	var field []string
	req := make(map[string]string)

	field = append(field, "clusterId")
	req["clusterId"] = this.ClusterId

	if this.Name != "" {
		field = append(field, "name")
		req["name"] = this.Name
	}

	if this.Desc != "" {
		field = append(field, "description")
		req["description"] = this.Desc
	}

	for i, n := range this.Rmname{
		key := fmt.Sprintf("names.%d",i)
		field = append(field,key)
		req[key] = n
	}
	//if len(this.Rmname) > 0 {
	//	field = append(field, "names")
	//	for i, n := range this.Rmname {
	//		this.Rmname[i] = "\"" + n + "\""
	//	}
	//
	//	req["names"] = fmt.Sprintf("[%s]", strings.Join(this.Rmname, ","))
	//}
	return field, req
}

// CreateNamespace 创建命名空间
func (this NSpace) CreateNamespace() error {
	if this.ClusterId == "" {
		return errors.New(_const.ClusterIDEmpty)
	}

	if this.Name == "" {
		return errors.New(_const.NamespaceNameEmpty)
	}

	if this.Desc == "" {
		this.Desc = _const.NSDefaultDesc
	}

	field, reqmap := this.queryNSInfo()
	pubMap := public.PublicParam("CreateClusterNamespace", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + _constv1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)
	if debug {
		log.Printf("[创建命名空间]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]", public.API_URL+reqURL, this.SecretKey, signStr, sign)
	}

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ns NSInfo

	err = json.Unmarshal(data, &ns)
	if err != nil {
		return err
	}

	if ns.Code != 0 {
		return errors.New(ns.CodeDesc + ns.Message)
	}

	return nil
}

// DeleteNamespace 删除命名空间
func (this NSpace) DeleteNamespace() error {
	if this.ClusterId == "" {
		return errors.New(_const.ClusterIDEmpty)
	}

	if len(this.Rmname) == 0 {
		return errors.New(_const.DeleteNameLengthZero)
	}

	field, reqmap := this.queryNSInfo()
	pubMap := public.PublicParam("DeleteClusterNamespace", this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + _constv1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)
	if debug {
		log.Printf("[删除命名空间]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]", public.API_URL+reqURL, this.SecretKey, signStr, sign)
	}

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var ns NSInfo

	err = json.Unmarshal(data, &ns)
	if err != nil {
		return err
	}

	if ns.Code != 0 {
		return errors.New(ns.CodeDesc + ns.Message)
	}

	return nil
}
