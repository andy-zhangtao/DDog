package mongo

import (
	"errors"
	"os"
	"time"

	"github.com/andy-zhangtao/DDog/const"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var endpoint = os.Getenv(_const.EnvMongo)
var username = os.Getenv(_const.EnvMongoName)
var password = os.Getenv(_const.EnvMongoPasswd)
var dbname = os.Getenv(_const.EnvMongoDB)
var session *mgo.Session

const (
	ModuleName = "Mongo InIt"
)

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
	//_hulk_client.Run()
	logrus.Println("=====Connect Mongo=====")
	err := check()
	if err != nil {
		logrus.Panic(err)
	}

	if username != "" || password != "" {
		dialInfo := &mgo.DialInfo{
			Addrs:    []string{endpoint},
			Database: dbname,
			Username: username,
			Password: password,
			Timeout:  30 * time.Second,
		}

		session, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			panic(err)
		}
	} else {
		session, err = mgo.Dial(endpoint)
	}
	b, err := session.BuildInfo()
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{"Mongo Server": b.Version}).Info(ModuleName)
}

func GetSession() *mgo.Session {
	return session
}

func getCloudMongo() *mgo.Database {
	return session.Clone().DB(_const.CloudMongoDBName)
}

func GetNameSpaceColleciton() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoNamespaceCol)
}
