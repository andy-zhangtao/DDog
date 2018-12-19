package qcloud

import (
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	sc "github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/sirupsen/logrus"
	"strings"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/4.

// GetInstanceInfo 获取指定服务的实例信息
func GetInstanceInfo(name, namespace string) (instance []service.Instance, err error) {

	svc, err := svcconf.GetSvcConfByName(name, namespace)
	if err != nil {
		return
	}

	var md *metadata.MetaData
	switch namespace {
	case "proenv":
		fallthrough
	case "release":
		md, err = metadata.GetMetaDataByRegion("", namespace)
	case "testenv":
		md, err = metadata.GetMetaDataByRegion("", "testenv")
	case "autoenv":
		md, err = metadata.GetMetaDataByRegion("", "autoenv")
	default:
		md, err = metadata.GetMetaDataByRegion("")
	}
	//if namespace == "proenv" {
	//	md, err = metadata.GetMetaDataByRegion("", namespace)
	//} else {
	//	md, err = metadata.GetMetaDataByRegion("")
	//}

	if err != nil {
		return
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: svc.SvcName,
		Namespace:   namespace,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)

	smData, err := q.QueryInstance()
	if err != nil {
		return
	}

	if smData.Code != 0 {
		err = errors.New(fmt.Sprintf("Query Instance Error [%d] Msg[%s][%s]", smData.Code, smData.Message, smData.CodeDesc))
		return
	}

	for _, i := range smData.Data.Instance {
		logrus.WithFields(logrus.Fields{"status": i.Status, "name": i.Name}).Info(ModuleName)
		if strings.ToLower(i.Status) == "running" {
			instance = append(instance, i)
		}
	}
	return instance, nil
}

//ModifyInstancesReplica  修改实例副本集数量
func ModifyInstancesReplica(name, namespace string, replica int) (err error) {
	scf, err := svcconf.GetSvcConfByName(name, namespace)
	if err != nil {
		return err
	}

	//scf.Replicas = replica
	var md *metadata.MetaData

	switch namespace {
	case "proenv":
		fallthrough
	case "release":
		md, err = metadata.GetMetaDataByRegion("", namespace)
	case "testenv":
		md, err = metadata.GetMetaDataByRegion("", "testenv")
	case "autoenv":
		md, err = metadata.GetMetaDataByRegion("", "autoenv")
	default:
		md, err = metadata.GetMetaDataByRegion("")
	}
	if err != nil {
		return err
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: scf.SvcName,
		Namespace:   scf.Namespace,
		ScaleTo:     replica,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)

	_, err = q.ModeifyInstance()
	if err != nil {
		return
	}

	scf.Deploy = _const.ModifyReplica
	scf.Status = _const.ModifyReplica
	return sc.UpdateSvcConf(&scf)
}
