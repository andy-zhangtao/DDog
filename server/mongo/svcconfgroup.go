package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.

func MongoSvcConfGroup() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMonogSvcConfGroup)
}

func GetSvcConfGroupByName(name, ns string) (scg interface{}, err error) {
	err = MongoSvcConfGroup().Find(bson.M{"name": name, "namespace": ns}).One(&scg)
	return
}

// GetAllSvcConfGroupByNs 获取指定命名空间下的所有服务编排数据
func GetAllSvcConfGroupByNs(ns string) (scg []interface{}, err error) {
	err = MongoSvcConfGroup().Find(bson.M{"namespace": ns}).All(&scg)
	return
}

func SaveSvcConfGroup(scg interface{}) error {
	return MongoSvcConfGroup().Insert(&scg)
}

func DeleteSvcConfGroup(id string) (err error) {
	err = MongoSvcConfGroup().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return
}
