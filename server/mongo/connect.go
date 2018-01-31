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
	log.Println("=====Connect Mongo=====")
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

func GetClusterById(id string) (ns interface{}, err error) {
	err = MongoClusterCol().Find(bson.M{"clusterid": id}).One(&ns)
	return
}
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

func MongoSvcConfGroup() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMonogSvcConfGroup)
}

func GetSvcConfGroupByName(name, ns string) (scg interface{}, err error) {
	err = MongoSvcConfGroup().Find(bson.M{"name": name, "namespace": ns}).One(&scg)
	return
}

// GetAllSvcConfGroupByNs 获取指定命名空间下的所有服务编排数据
func GetAllSvcConfGroupByNs(ns string) (scg []interface{}, err error) {
	err = MongoSvcConfGroup().Find(bson.M{"namespace": ns}).All(&scg)
	return
}

func SaveSvcConfGroup(scg interface{}) error {
	return MongoSvcConfGroup().Insert(&scg)
}

func DeleteSvcConfGroup(id string) (err error) {
	err = MongoSvcConfGroup().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return
}
