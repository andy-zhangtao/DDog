package k8s

//import (
//	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
//	"fmt"
//	"crypto/tls"
//	"github.com/andy-zhangtao/DDog/const"
//	"log"
//	"net/http"
//	"io/ioutil"
//	"encoding/json"
//)
//
//func (k *K8sMetaData) GetPods() (*k8smodel.K8s, error){
//
//	path := ""
//	switch k.Version {
//	case "1.7":
//		path = fmt.Sprintf("%s/apis/apps/v1beta1/namespaces/%s/deployments/%s/pods", k.Endpoint, k.Namespace)
//	}
//
//	if _const.DEBUG {
//		log.Printf("[GetDeployMent] path:[%s]\n", path)
//	}
//
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	}
//
//	client := &http.Client{Transport: tr}
//
//	req, err := http.NewRequest("GET", path, nil)
//	if err != nil {
//		return nil, err
//	}
//	req.SetBasicAuth("admin", k.Token)
//
//	if _const.DEBUG {
//		log.Printf("[GetDeployMent] Request Header [%v]\n", req.Header)
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//
//	data, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	defer resp.Body.Close()
//	if _const.DEBUG {
//		log.Printf("[GetDeployMent] Receive Body [%s]\n", string(data))
//	}
//
//	err = json.Unmarshal(data, &kapi)
//	if err != nil {
//		return nil, err
//	}
//
//	return &kapi, nil
//}
