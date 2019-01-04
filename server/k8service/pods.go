package k8service

import (
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/k8s"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/sirupsen/logrus"
)

func GetDeploymentsPods(region, namespace, svcname string) (pods k8smodel.K8sPodsInfo, err error) {
	logrus.WithFields(logrus.Fields{"Ready to extract K8s pods data": namespace}).Info(ModuleName)

	svc, err := svcconf.GetSvcConfByName(svcname, namespace)
	if err != nil {
		return
	}

	k8sMeta, err := GetK8sClusterWithNamespace(region, namespace)
	if err != nil {
		err = errors.New(fmt.Sprintf("Get K8s Cluster Error [%s]", err.Error()))
		return
	}

	k8m := k8s.K8sMetaData{
		Endpoint:  k8sMeta.Endpoint,
		Token:     k8sMeta.Token,
		Version:   "1.7",
		Namespace: namespace,
		Svcname:   svc.SvcName,
	}

	return k8m.GetServiceV1BetaPods()
}
