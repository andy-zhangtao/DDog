package k8s

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	GetAllDeployMent = iota
	GetSpecDeployMent
	GetAllService
	GetSpecService
	GetAllPods
	GetSpecPods
)

const (
	ModuleName = "K8s-API-Agent"
)

func (k *K8sMetaData) invokeK8sAPI(kind int) ([]byte, error) {

	path := ""
	switch k.Version {
	case "1.7":
		switch kind {
		case GetAllDeployMent:
			path = fmt.Sprintf("%s/apis/apps/v1beta1/namespaces/%s/deployments", k.Endpoint, k.Namespace)
		case GetSpecDeployMent:
			path = fmt.Sprintf("%s/apis/apps/v1beta1/namespaces/%s/deployments/%s", k.Endpoint, k.Namespace, k.Svcname)
		case GetAllService:
			path = fmt.Sprintf("%s/api/v1/namespaces/%s/services", k.Endpoint, k.Namespace)
		case GetSpecService:
			path = fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s", k.Endpoint, k.Namespace, k.Svcname)
		case GetAllPods:
			path = fmt.Sprintf("%s/api/v1/namespaces/%s/pods", k.Endpoint, k.Namespace)
		case GetSpecPods:
			query := fmt.Sprintf("qcloud-app=%s",k.Svcname)
			query = url.QueryEscape(query)
			path = fmt.Sprintf("%s/api/v1/namespaces/%s/pods?labelSelector=%s", k.Endpoint, k.Namespace, query)
		}
	}

	logrus.WithFields(logrus.Fields{"Path": path}).Info(ModuleName)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("admin", k.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return data, nil
}

func (k *K8sMetaData) GetDeployMent() (*k8smodel.K8s, error) {

	var kapi k8smodel.K8s

	path := ""
	switch k.Version {
	case "1.7":
		path = fmt.Sprintf("%s/apis/apps/v1beta1/namespaces/%s/deployments", k.Endpoint, k.Namespace)
	}

	if _const.DEBUG {
		log.Printf("[GetDeployMent] path:[%s]\n", path)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("admin", k.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(data, &kapi)
	if err != nil {
		return nil, err
	}

	return &kapi, nil
}

func (k *K8sMetaData) GetServiceV1BetaPods() (pods k8smodel.K8sPodsInfo, err error) {
	data, err := k.invokeK8sAPI(GetSpecPods)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get Specify Pods Error [%s]", err.Error()))
		return
	}

	if err = json.Unmarshal(data, &pods); err != nil {
		logrus.Println(string(data))
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s]", err.Error()))
	}

	return
}

// GetServiceV1Beta 获取V1Beta版本的Service数据
func (k *K8sMetaData) GetServiceV1Beta() (service k8smodel.K8sServiceInfo, err error) {
	data, err := k.invokeK8sAPI(GetSpecService)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get ALl Service Error [%s]", err.Error()))
		return
	}

	if err = json.Unmarshal(data, &service); err != nil {
		logrus.Println(string(data))
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s]", err.Error()))
	}

	return
}

// GetServiceV1Beta 获取V1版本的Service数据
func (k *K8sMetaData) GetServiceV1() (service k8smodel.K8sService, err error) {
	data, err := k.invokeK8sAPI(GetSpecService)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get ALl Service Error [%s]", err.Error()))
		return
	}

	if err = json.Unmarshal(data, &service); err != nil {
		logrus.Println(string(data))
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s]", err.Error()))
	}

	return
}

// GetDeployMentV1Beta 获取V1Beta版本的Deploymen数据
func (k *K8sMetaData) GetDeployMentV1Beta() (deploy k8smodel.K8sDeploymentInfo, err error) {
	data, err := k.invokeK8sAPI(GetAllDeployMent)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get ALl Deploy Error [%s]", err.Error()))
		return
	}

	if err = json.Unmarshal(data, &deploy); err != nil {
		logrus.Println(string(data))
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s]", err.Error()))
	}

	return
}

// GetDeployMentStatus 获取指定Deployment的运行状态
// DeployMent 对应的是容器云的Service
func (k *K8sMetaData) GetDeployMentStatus() (*k8smodel.K8sDeployment, error) {
	data, err := k.invokeK8sAPI(GetSpecDeployMent)
	if err != nil {
		return nil, err
	}

	k8d, err := k8smodel.K8dUnmarshal(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%v Raw Data[%s]", err, string(data)))
	}

	if k8d == nil {
		return nil, errors.New("Get K8sDeployment Error")
	}

	return k8d, nil
}

// GetPodsInNamespace 获取指定命名空间中的Pods信息
// Pods对应的是容器云中的实例信息
func (k *K8sMetaData) GetPodsInNamespace() (*k8smodel.K8sPods, error) {
	data, err := k.invokeK8sAPI(GetAllPods)
	if err != nil {
		return nil, err
	}

	k8p, err := k8smodel.K8pUnmarshal(data)
	if err != nil {
		return nil, err
	}

	if k8p == nil {
		return nil, errors.New("Get K8sDeployment Error")
	}

	return k8p, nil
}
