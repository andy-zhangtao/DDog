package k8service

import (
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
//获取K8s DeployMent数据

func AddNewK8sClusterConfigeDate(kc k8sconfig.K8sCluster) (err error) {
	if kc.Region == "" || kc.Endpoint == "" || kc.Token == "" {
		return errors.New(fmt.Sprintf("Regin/Endpoint/Token Empty! [%v]", kc))
	}

	if err = mongo.SaveK8sClusterData(kc); err != nil {
		err = errors.New(fmt.Sprintf("Add New K8s Configure Error [%s] K8s [%v]", err.Error(), kc))
	}

	return
}

func GetALlK8sCluster() (kc []k8sconfig.K8sCluster, err error) {
	if kc, err = mongo.GetAllK8sClusterData(); err != nil {
		err = errors.New(fmt.Sprintf("Query All K8s Data Error [%s]", err.Error()))
	}

	return
}

func GetK8sCluster(region string) (kc k8sconfig.K8sCluster, err error) {
	if region == "" {
		region = os.Getenv(_const.EnvRegion)
	}

	if region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if kc, err = mongo.GetK8sClusterData(region); err != nil {
		err = errors.New(fmt.Sprintf("Query K8s Cluster Data Error [%s] Region [%s]", err.Error(), region))
	}

	return
}

func UpdateK8sCluster(kc k8sconfig.K8sCluster) (err error) {
	if kc.Region == "" {
		kc.Region = os.Getenv(_const.EnvRegion)
	}

	if kc.Region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if err = mongo.UpdateK8sClusterData(kc); err != nil {
		errors.New(fmt.Sprintf("Update K8s Cluster Data Error [%s] K8s [%s]", err.Error(), kc))
	}

	return
}

func DeleteK8sCluster(region string) (err error) {
	if region == "" {
		region = os.Getenv(_const.EnvRegion)
	}

	if region == "" {
		err = errors.New(fmt.Sprintf("Region Empty!"))
		return
	}

	if kc, err := mongo.GetK8sClusterData(region); err != nil {
		return errors.New(fmt.Sprintf("Query K8s Cluster Data Error [%s] Region [%s]", err.Error(), region))
	} else {
		if err = mongo.DeleteK8sClusterDataByID(kc.ID); err != nil {
			return errors.New(fmt.Sprintf("Delete K8s Cluster Data Error [%s] ID [%s]", err.Error(), kc.ID))
		}
	}

	return
}

func DeleteK8sClusterByID(id string) (err error) {
	if err = mongo.DeleteK8sClusterDataByID(bson.ObjectIdHex(id)); err != nil {
		return errors.New(fmt.Sprintf("Delete K8s Cluster Data Error [%s] ID [%s]", err.Error(), id))
	}

	return
}
