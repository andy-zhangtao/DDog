package qcloud

import (
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/server/svcconf"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/4.

// GetInstanceInfo 获取指定服务的实例信息
func GetInstanceInfo(name, namespace string) (instance []service.Instance, err error) {

	svc, err := svcconf.GetSvcConfByName(name, namespace)
	if err != nil {
		return
	}

	md, err := metadata.GetMetaDataByRegion("")
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

	return smData.Data.Instance, nil
}
