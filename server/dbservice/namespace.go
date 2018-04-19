package dbservice

import (
	"github.com/andy-zhangtao/DDog/model/caasmodel"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
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

	if ns.CreateTime == "" {
		ns.CreateTime = time.Now().String()
	}
	return mongo.GetNameSpaceColleciton().Insert(&ns)
}

func DeleteNamespaceByID(id bson.ObjectId) (err error) {
	return mongo.GetNameSpaceColleciton().RemoveId(id)
}

func GetNamespaceByName(name string) (ns caasmodel.NameSpace, err error) {
	err = mongo.GetNameSpaceColleciton().Find(bson.M{"name": name}).One(&ns)
	return
}

func GetNamespaceByOwner(owner string) (ns []caasmodel.NameSpace, err error) {
	err = mongo.GetNameSpaceColleciton().Find(bson.M{"owner": owner}).All(&ns)
	return
}

func GetNamespaceByOwnerAndName(name, owner string) (ns caasmodel.NameSpace, err error) {
	err = mongo.GetNameSpaceColleciton().Find(bson.M{"name": name, "owner": owner}).One(&ns)
	return
}

func UpdateNamespace(namespace caasmodel.NameSpace) (err error) {
	ns, err := GetNamespaceByOwnerAndName(namespace.Name, namespace.Owner)
	if err != nil {
		return SaveNamespace(namespace)
	}

	if namespace.Desc != "" {
		err = DeleteNamespaceByID(ns.ID)
		if err != nil {
			return err
		}
		ns.Desc = namespace.Desc
		return SaveNamespace(namespace)
	}

	return nil
}
