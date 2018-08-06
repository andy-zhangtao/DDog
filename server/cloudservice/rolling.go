package cloudservice

import (
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"net/http"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/andy-zhangtao/DDog/server/qcloud"
)

const ModuleName = "RollingUpService"

type MyWriter struct {
	header http.Header
	data   string
	status int
}

func (this *MyWriter) Header() http.Header {
	return this.header
}

func (this *MyWriter) Write(d []byte) (int, error) {
	this.data = string(d)
	return len(this.data), nil
}

func (this *MyWriter) WriteHeader(h int) {
	this.status = h
	logrus.WithFields(logrus.Fields{"header": h, "status": this.status}).Info("DeployAgent-Writer")
}

// RollingUpService 灰度升级
func RollingUpService(servicename, namespace string, percent float64) (err error) {

	scf, err := svcconf.GetSvcConfByName(servicename, namespace)
	if err != nil {
		return
	}

	rollCons, _, left := scf.CountInstances(percent)
	if left <= 0 && len(rollCons) == 0 {
		scf.Deploy = 3
		svcconf.UpdateSvcConf(scf)
		return
	}

	if len(rollCons) == 0 {
		return
	}

	writer := MyWriter{
		header: make(map[string][]string),
	}

	logrus.WithFields(logrus.Fields{"url": fmt.Sprintf("?svcname=%s&namespace=%s&replicas=%d", servicename, scf.Namespace, len(rollCons))}).Info(ModuleName)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("?svcname=%s&namespace=%s&replicas=%d", servicename, scf.Namespace, len(rollCons)), nil)
	if err != nil {
		return err
	}

	qcloud.RunService(&writer, r)
	logrus.WithFields(logrus.Fields{"Status": writer.status, "Header": writer.header}).Info(ModuleName)

	return
}
