package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.

// MongoMonitor 获取监控数据库实例
func MongoMonitor() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoMonitor)
}

// MongoGetALlMonitorByKind 获取指定类型的监控数据
// kind 事件类型
func MongoGetALlMonitorByKind(kind string) (mb []interface{}, err error) {
	err = MongoMonitor().Find(bson.M{"kind": kind}).All(&mb)
	return
}

// MongoGetAllMonitorByKindAndStatus 获取指定类型指定状态的监控数据
// kind 事件类型
// status 事件状态
func MongoGetAllMonitorByKindAndStatus(kind, status string) (mb []interface{}, err error) {
	err = MongoMonitor().Find(bson.M{"kind": kind, "status": status}).All(&mb)
	return
}

// SaveMonitor 插入一条监控数据
func SaveMonitor(scg interface{}) error {
	return MongoMonitor().Insert(&scg)
}

// DeleteMonitorByKind 删除指定类型的监控数据
// kind 事件类型
func DeleteMonitorByKind(kind string) (err error) {
	_, err = MongoMonitor().RemoveAll(bson.M{"kind": kind})
	return
}

// DeleteMonitorByKindAndStatus 删除指定类型指定状态的监控数据
// kind 事件类型
// status 事件状态
func DeleteMonitorByKindAndStatus(kind, status string) (err error) {
	_, err = MongoMonitor().RemoveAll(bson.M{"kind": kind, "status": status})
	return
}
