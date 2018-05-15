package k8sconfig

import "gopkg.in/mgo.v2/bson"

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
//K8sCluster 集群数据
//必须和集群数据相一致
type K8sCluster struct {
	ID         bson.ObjectId `json:"_id" bson:"_id"`
	Name       string        `json:"name" bson:"name"`
	Region     string        `json:"region" bson:"region"`
	Token      string        `json:"token" bson:"token"`
	Endpoint   string        `json:"endpoint" bson:"endpoint"`
	Namespace  string        `json:"namespace" bson:"namespace"`
	UpdateTime string        `json:"update_time" bson:"update"`
}
