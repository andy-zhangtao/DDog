package mongo

import (
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"errors"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var endpoint = os.Getenv(_const.EnvMongo)
var username = os.Getenv(_const.EnvMongoName)
var password = os.Getenv(_const.EnvMongoPasswd)
var dbname = os.Getenv(_const.EnvMongoDB)
var session *mgo.Session

func check() error {
	if endpoint == "" {
		return errors.New(_const.EnvMongoNotFound)
	}

	if dbname == "" {
		return errors.New(_const.EnvMongoDBNotFound)
	}
	return nil
}

func init() {
	err := check()
	if err != nil {
		log.Panic(err)
	}

	if username != "" || password != "" {
		dialInfo := &mgo.DialInfo{
			Addrs:    []string{endpoint},
			Database: dbname,
			Username: username,
			Password: password,
		}

		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
	} else {
		session, err = mgo.Dial(endpoint)
	}

}

func GetSession() *mgo.Session {
	return session
}

func getCloudMongo() *mgo.Database {
	return session.Clone().DB(_const.CloudMongoDBName)
}

func MongoCloudMetadata() *mgo.Collection {
	c := getCloudMongo()
	return c.C(_const.CloudMongoCollection)
}

func SaveMetaData(metadata interface{}) error {
	c := MongoCloudMetadata()
	return c.Insert(&metadata)
}

func FindMetaDataByRegion(region string) (int, error) {
	c := MongoCloudMetadata()
	return c.Find(bson.M{"region": region}).Count()
}

func GetMetaDataByRegion(region string, metadata interface{}) (err error) {
	c := MongoCloudMetadata()
	err = c.Find(bson.M{"region": region}).One(metadata)
	return
}

func GetALlMetaData() (m []interface{}, err error) {
	c := MongoCloudMetadata()
	//var m []interface{}
	err = c.Find(nil).All(&m)
	return
}
