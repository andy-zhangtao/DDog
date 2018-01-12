package cluster

import (
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
)

// GetClusterByID 根据ID获取集群数据
func GetClusterByID(clusterid string) (cluster cvm.ClusterInfo_data_clusters, err error) {

	cs, err := mongo.GetClusterById(clusterid)
	if err != nil {
		return
	}

	data, err := bson.Marshal(cs)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &cluster)
	if err != nil {
		return
	}

	return
}
