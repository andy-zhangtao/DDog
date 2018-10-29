package cloudservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

// Restart 重启服务
// 逐个重启服务容器
func Restart(servicename, namespace string) (err error) {
	var md *metadata.MetaData
	switch namespace {
	case "proenv":
		//	预发布环境
		fallthrough
	case "release":
		//	预发布环境
		md, err = metadata.GetMetaDataByRegion("", namespace)
		if err != nil {
			return errors.New(_const.RegionNotFound)
		}
	default:
		md, err = metadata.GetMetaDataByRegion("")
		if err != nil {
			return errors.New(_const.RegionNotFound)
		}
	}

	scf, err := svcconf.GetSvcConfByName(servicename, namespace)
	if err != nil {
		return errors.New(fmt.Sprintf("Get Svc Error [%v] name [%v] env [%v]", err.Error(), servicename, namespace))
	}

	writer := MyWriter{
		header: make(map[string][]string),
	}

	url := fmt.Sprintf("?clusterid=%s&svcname=%s&namespace=%s", md.ClusterID, scf.SvcName, scf.Namespace)
	logrus.WithFields(logrus.Fields{"url": url}).Info(ModuleName)

	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	qcloud.ReinstallService(&writer, r)
	logrus.WithFields(logrus.Fields{"Status": writer.status, "Header": writer.header}).Info(ModuleName)

	if writer.status != 0 {
		return errors.New(fmt.Sprintf("QCloud Retrun Error Status [%v] Header [%v]", writer.status, writer.header))
	}

	producer, _ := nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())

	err = NotifyDevEx(servicename, namespace, producer, 27)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Notify DevEx Error": err}).Error(ModuleName)
		return
	}

	go waitRestartEnd(md, scf.SvcName, servicename, namespace, producer)
	return

}

//svcname 集群中deployment名称
//servicename 对外暴露的统一服务名
func waitRestartEnd(md *metadata.MetaData, svcname, servicename, namespace string, producer *nsq.Producer) {
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		Namespace:   namespace,
		SecretKey:   md.Skey,
		ServiceName: svcname,
	}
	q.SetDebug(true)

	for {
		resp, err := q.QuerySvcInfo()
		if err != nil || resp.Code != 0 {
			time.Sleep(5 * time.Second)
			continue
		}

		if strings.ToLower(resp.Data.ServiceInfo.Status) == "normal" {
			NotifyDevEx(servicename, namespace, producer, 28)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

// NotifyDevex 通知Devex当前服务状态
// 27正在重启
// 28重启结束
func NotifyDevEx(servicename, namespace string, producer *nsq.Producer, state int) error {

	req := struct {
		ProjectID string `json:"project_id"`
		Stage     int    `json:"stage"`
		DeployEnv string `json:"deploy_env"`
	}{
		servicename,
		state,
		namespace,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		//logrus.WithFields(logrus.Fields{"Marshal DevEx Request Error": err}).Error(ModuleName)
		return err
	}

	err = producer.Publish("DevEx-Request-Status", data)
	if err != nil {
		return err
	}

	return err
}
