package k8service

import (
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/k8s"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
//获取K8s Deployment数据

//GetK8sDeployMent 获取指定K8s集群的DeployMent
//region 集群区域
//name 集群控制命名空间名称
func GetK8sDeployMent(region, name string) (deploy k8smodel.K8sDeploymentInfo, err error) {
	logrus.WithFields(logrus.Fields{"Ready to extract K8s deployment data": true}).Info(ModuleName)

	k8sMeta, err := GetK8sCluster(region)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get K8s Cluster Error [%s]", err.Error()))
		return
	}

	k8m := k8s.K8sMetaData{
		Endpoint:  k8sMeta.Endpoint,
		Token:     k8sMeta.Token,
		Version:   "1.7",
		Namespace: name,
	}

	return k8m.GetDeployMentV1Beta()
}

//获取指定集群和命名空间的所有部署服务信息
func GetK8sDeployMents(k8c k8sconfig.K8sCluster, namespace string) (deploys k8smodel.K8sDeploymentInfo, err error) {
	k8m := k8s.K8sMetaData{
		Endpoint:  k8c.Endpoint,
		Token:     k8c.Token,
		Version:   "1.7",
		Namespace: namespace,
	}

	return k8m.GetDeployMentV1Beta()
}

//获取指定集群和命名空间的指定部署服务信息
func GetK8sSpecifyDeployMent(k8c k8sconfig.K8sCluster, namespace, svcname string) (deploy *k8smodel.K8sDeployment, err error) {
	k8m := k8s.K8sMetaData{
		Endpoint:  k8c.Endpoint,
		Token:     k8c.Token,
		Version:   "1.7",
		Namespace: namespace,
		Svcname:   svcname,
	}

	return k8m.GetDeployMentStatus()
}

//获取指定集群和命名空间的指定服务信息
func GetK8sSpecifyService(k8c k8sconfig.K8sCluster, namespace, svcname string) (service k8smodel.K8sService, err error) {
	k8m := k8s.K8sMetaData{
		Endpoint:  k8c.Endpoint,
		Token:     k8c.Token,
		Version:   "1.7",
		Namespace: namespace,
		Svcname:   svcname,
	}

	return k8m.GetServiceV1()
}
