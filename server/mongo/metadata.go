package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.

func MongoMetadataCol() *mgo.Collection {
	c := getCloudMongo()
	return c.C(_const.CloudMongoMeataDataCol)
}

func SaveMetaData(metadata interface{}) error {
	c := MongoMetadataCol()
	return c.Insert(&metadata)
}

func DeleteMetaData(id string) (err error) {
	return MongoMetadataCol().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

func FindMetaDataByRegion(region string) (md interface{}, err error) {
	err = MongoMetadataCol().Find(bson.M{"region": region}).One(&md)
	return
}

// FindAllMetaData 检索所有的MetaData数据
func FindAllMetaData() (md []interface{}, err error) {
	err = MongoMetadataCol().Find(nil).All(&md)
	return
}

func GetMetaDataByRegion(region string, metadata interface{}) (err error) {
	c := MongoMetadataCol()
	err = c.Find(bson.M{"region": region}).One(metadata)
	return
}

func GetALlMetaData() (m []interface{}, err error) {
	c := MongoMetadataCol()
	err = c.Find(nil).All(&m)
	return
}