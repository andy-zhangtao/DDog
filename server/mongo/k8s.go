package mongo

import (
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strings"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
func MongoK8sConfigCol() *mgo.Collection {
	return getCloudMongo().C(_const.CloudK8sClusterCol)
}

// GetK8sClusterData 按照地区获取K8s数据
func GetK8sClusterData(region string) (kc k8sconfig.K8sCluster, err error) {
	err = MongoK8sConfigCol().Find(bson.M{"region": region}).One(&kc)
	return
}

func GetAllK8sClusterData() (kc []k8sconfig.K8sCluster, err error) {
	err = MongoK8sConfigCol().Find(nil).All(&kc)
	return
}

func DeleteK8sClusterDataByID(id bson.ObjectId) (err error) {
	return MongoK8sConfigCol().RemoveId(id)
}

func SaveK8sClusterData(kc k8sconfig.K8sCluster) (err error) {
	if kc.ID == "" {
		kc.ID = bson.NewObjectId()
	}

	kc.UpdateTime = time.Now().Format("2006-01-02T15:04")
	return MongoK8sConfigCol().Insert(&kc)
}

func UpdateK8sClusterData(kc k8sconfig.K8sCluster) (err error) {
	if kc.ID != "" {
		if err = DeleteK8sClusterDataByID(kc.ID); err != nil {
			return
		}

		return SaveK8sClusterData(kc)
	}

	if okc, err := GetK8sClusterData(kc.Region); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return SaveK8sClusterData(kc)
		}
		return err
	} else {
		okc.Endpoint = kc.Endpoint
		okc.Token = kc.Token
		okc.Name = kc.Name
		okc.Namespace = kc.Namespace
		okc.UpdateTime = time.Now().Format("2006-01-02T15:04")
		return MongoK8sConfigCol().Update(bson.M{"_id": okc.ID}, okc)
	}
}
