package namespace

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"log"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/url"
)

var debug = false

type NSpace struct {
	Pub       public.Public `json:"pub"`
	ClusterId string        `json:"cluster_id"`
	SecretKey string        `json:"secret_key"`
	sign      string
}

// http://json.golang.chinazt.cc/
// 自动生成, 使用前请校验
type NSInfo_data_namespaces struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdat"`
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

func (this NSpace) QueryNSInfo()(*NSInfo, error) {
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
		log.Printf("[获取命名空间信息]请求URL[%s]密钥[%s]签名内容[%s]生成签名[%s]",public.API_URL + reqURL,this.SecretKey,signStr,sign)
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

	return field, req
}
