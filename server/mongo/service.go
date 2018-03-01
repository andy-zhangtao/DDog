package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.
func MongoSVCCol() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoSVCCol)
}

func SaveService(svc interface{}) error {
	return MongoSVCCol().Insert(&svc)
}

func DeleteSvcByName(ns, name string) error {
	change, err := MongoSVCCol().RemoveAll(bson.M{"servicename": name, "namespace": ns})
	if err != nil {
		return err
	}

	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}

func DeleteAllSvcByNs(ns string) error {
	change, err := MongoSVCCol().RemoveAll(bson.M{"namespace": ns})
	if err != nil {
		return err
	}

	//log.Println(change.Removed, change.Matched)
	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}

func GetAllSvcByNs(ns string) (svc []interface{}, err error) {
	err = MongoSVCCol().Find(bson.M{"namespace": ns}).All(&svc)
	return
}

func GetSvcByName(ns, name string) (svc interface{}, err error) {
	err = MongoSVCCol().Find(bson.M{"servicename": name, "namespace": ns}).One(&svc)
	return
}