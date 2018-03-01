package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.
func MongoSvcConfig() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoSvcConfig)
}

func SaveSvcConfig(conf interface{}) error {
	return MongoSvcConfig().Insert(&conf)
}

func GetSvcConfByID(id string) (conf interface{}, err error) {
	err = MongoSvcConfig().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&conf)
	return
}

func GetSvcConfNs(ns string) (conf []interface{}, err error) {
	err = MongoSvcConfig().Find(bson.M{"namespace": ns}).All(&conf)
	return
}

func GetSvcConfByName(name, ns string) (conf interface{}, err error) {
	err = MongoSvcConfig().Find(bson.M{"name": name, "namespace": ns}).One(&conf)
	return
}

func DeleteSvcConfById(id string) (err error) {
	err = MongoSvcConfig().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return
}

func DeleteSvcConfByNs(ns string) error {
	change, err := MongoSvcConfig().RemoveAll(bson.M{"namespace": ns})
	if err != nil {
		return err
	}

	//log.Println(change.Removed, change.Matched)
	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}