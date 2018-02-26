package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/26.
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
