package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.

func MongoNamespaceCol() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoNamespaceCol)
}

func SaveNamespace(namespace interface{}) error {
	return MongoNamespaceCol().Insert(&namespace)
}

func DeleteNamespaceByName(clusterID, name string) error {
	change, err := MongoNamespaceCol().RemoveAll(bson.M{"name": name, "clusterid": clusterID})
	if err != nil {
		return err
	}

	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil

}

func DeleteAllNamespaceByCID(clusterID string) error {
	change, err := MongoNamespaceCol().RemoveAll(bson.M{"cluster_id": clusterID})
	if err != nil {
		return err
	}

	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}

func GetAllNamespaceByCID(clusterID string) (ns []interface{}, err error) {
	err = MongoNamespaceCol().Find(bson.M{"clusterid": clusterID}).All(&ns)
	return
}

func GetNamespaceByName(clusterID, name string) (ns interface{}, err error) {
	err = MongoNamespaceCol().Find(bson.M{"clusterid": clusterID, "name": name}).One(&ns)
	return
}