package public

import (
	gg "github.com/andy-zhangtao/gogather/time"
	ga "github.com/andy-zhangtao/gogather/random"
	gz "github.com/andy-zhangtao/gogather/zsort"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

const API_URL = "https://ccs.api.qcloud.com/v2/index.php?"
const Repostiory_API_URL = "https://ccr.api.qcloud.com/v2/index.php?"

type Public struct {
	//Action   string `json:"action"`
	SecretId string `json:"secret_id"`
	Region   string `json:"region"`
}

var (
	PubilcField = []string{"Action", "SecretId", "Region", "Timestamp", "Nonce"}
)

// PublicParam生成公共请求数据
// 包括 Action/SecretId/Region/Timestamp/Nonce
func PublicParam(action, region, secretId string) map[string]string {
	req := make(map[string]string)
	req["Action"] = action
	req["SecretId"] = secretId
	req["Region"] = region
	req["Timestamp"] = gg.GetTimeStamp(10)
	req["Nonce"] = ga.GetRandom(6)
	return req
}

func Generate() {

}

// generateSignature 生成请求签名字符串
// field 请求字段集合
// reqmap 待计算的请求map
// publicmap 公共请求map,调用public.PublicParam生成
func GenerateSignatureString(field []string, reqmap, publicMap map[string]string) string {
	field = append(field, PubilcField...)
	field = gz.DictSort(field)
	//log.Println(field)

	req := ""
	for k, v := range reqmap {
		publicMap[k] = v
	}
	for i, key := range field {
		if i == 0 {
			req = key + "=" + publicMap[key]
		} else {
			req += "&" + key + "=" + publicMap[key]
		}
	}
	return req
}

// generateSignature 生成最终的请求签名,使用HMAC-SHA1
// key加密key，req请求字符串
func GenerateSignature(key, req string) string {
	k := []byte(key)
	mac := hmac.New(sha1.New, k)
	mac.Write([]byte(req))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
