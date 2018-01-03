package mongo

import (
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"errors"
	"log"
	"gopkg.in/mgo.v2"
)

var endpoint = os.Getenv(_const.EnvMongo)
var session *mgo.Session

func check() error {
	if endpoint == "" {
		return errors.New(_const.EnvMongoNotFound)
	}
	return nil
}

func init() {
	err := check()
	if err != nil {
		log.Panic(err)
	}

	session, err = mgo.Dial(endpoint)
	if err != nil {
		panic(err)
	}

}

func GetSession() *mgo.Session {
	return session
}

func getCloudMongo() *mgo.Database {
	return session.DB(_const.CloudMongoDBName)
}

func MongoCloudMetadata() *mgo.Collection {
	c := getCloudMongo()
	return c.C(_const.CloudMongoCollection)
}

func SaveMetaData(metadata interface{}) error {
	c := MongoCloudMetadata()
	return c.Insert(&metadata)
}
