package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.
func MongoContainerCol() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoContainerCol)
}

func SaveContainer(con interface{}) error {
	return MongoContainerCol().Insert(&con)
}

func GetContaienrByName(conName, svcName, nsme string) (con interface{}, err error) {
	err = MongoContainerCol().Find(bson.M{"name": conName, "svc": svcName, "nsme": nsme}).One(&con)
	return
}

func GetContainerByID(id string) (con interface{}, err error) {
	err = MongoContainerCol().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&con)
	return
}

func GetContaienrBySvc(svcname, ns string) (con []interface{}, err error) {
	err = MongoContainerCol().Find(bson.M{"svc": svcname, "nsme": ns}).Sort("-idx").All(&con)
	return
}

func DeleteContainerById(id string) (err error) {
	err = MongoContainerCol().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return
}

func DeleteAllContainer(svcname, ns string) error {
	change, err := MongoContainerCol().RemoveAll(bson.M{"svc": svcname, "nsme": ns})
	if err != nil {
		return err
	}

	//log.Println(change.Removed, change.Matched)
	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}