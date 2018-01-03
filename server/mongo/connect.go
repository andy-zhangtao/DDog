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
	b, err := session.BuildInfo()
	if err != nil {
		panic(err)
	}

	log.Printf("Mongo Server[%s]\n", b.Version)
}

func GetSession() *mgo.Session {
	return session
}

func getCloudMongo() *mgo.Database {
	return session.Clone().DB(_const.CloudMongoDBName)
}

func MongoMetadataCol() *mgo.Collection {
	c := getCloudMongo()
	return c.C(_const.CloudMongoMeataDataCol)
}

func SaveMetaData(metadata interface{}) error {
	c := MongoMetadataCol()
	return c.Insert(&metadata)
}

func FindMetaDataByRegion(region string) (int, error) {
	c := MongoMetadataCol()
	return c.Find(bson.M{"region": region}).Count()
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

func MongoClusterCol() *mgo.Collection {
	c := getCloudMongo()
	return c.C(_const.CloudMongoClusterCol)
}

func SaveCluster(cluster interface{}) error {
	c := MongoClusterCol()
	return c.Insert(&cluster)
}

func DeleteCluster(id string) error {
	c := MongoClusterCol()
	change, err := c.RemoveAll(bson.M{"clusterid": id})
	if err != nil {
		return err
	}

	if change.Removed == 0 {
		return errors.New("There is no match record!")
	}
	return nil
}

func GetClusterByRegion(region string) (m []interface{}, err error) {
	c := MongoClusterCol()
	err = c.Find(bson.M{"region": region}).All(&m)
	return
}

func MongoNamespaceCol() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoNamespaceCol)
}

func SaveNamespace(namespace interface{}) error {
	return MongoNamespaceCol().Insert(&namespace)
}

//func GetNamespaceByName(name string)(namespace interface{}, err error){
//	c := MongoNamespaceCol()
//	c.Find(bson.M{""})
//}

func DeleteNamespaceByName(clusterID, name string) error {
	change, err := MongoNamespaceCol().RemoveAll(bson.M{"name": name, "clusterid": clusterID})
	if err != nil {
		return err
	}

	log.Println(change.Removed, change.Matched)
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

	log.Println(change.Removed, change.Matched)
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
