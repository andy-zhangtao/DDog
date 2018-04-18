package dbservice

import (
	"github.com/andy-zhangtao/DDog/model/caasmodel"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/18.
//对Namespace的Mongo操作

// GetAllNamesapce 获取所有命名空间
func GetAllNamesapce() (ns []caasmodel.NameSpace, err error) {
	err = mongo.GetNameSpaceColleciton().Find(nil).All(&ns)
	return
}

func SaveNamespace(ns caasmodel.NameSpace) (err error) {
	if ns.ID == "" {
		ns.ID = bson.NewObjectId()
	}
	return mongo.GetNameSpaceColleciton().Insert(&ns)
}

func DeleteNamespaceByID(id string) (err error) {
	return mongo.GetNameSpaceColleciton().RemoveId(id)
}

func GetNamespaceByName(name string) (ns caasmodel.NameSpace, err error) {
	err = mongo.GetNameSpaceColleciton().Find(bson.M{"name": name}).One(&ns)
	return
}
