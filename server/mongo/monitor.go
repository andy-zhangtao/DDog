package mongo

import (
	"gopkg.in/mgo.v2"
	"github.com/andy-zhangtao/DDog/const"
	"gopkg.in/mgo.v2/bson"
	"github.com/Sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.

// MongoMonitor 获取监控数据库实例
func MongoMonitor() *mgo.Collection {
	return getCloudMongo().C(_const.CloudMongoMonitor)
}

// MongoGetMonitorByName 通过服务名和命名空间查询监控信息
// Agent agent名称
// Svcname 服务名称
// Namespace 命名空间
func MongoGetMonitorByName(agent, svcname, namespace string) (mb interface{}, err error) {
	logrus.WithFields(logrus.Fields{"kind": agent, "svcname": svcname, "namespace": namespace}).Info("MongoGetMonitorByName")
	err = MongoMonitor().Find(bson.M{"kind": agent, "svcname": svcname, "namespace": namespace}).One(&mb)
	return
}

// MongoGetALlMonitorByKind 获取指定Agent的监控数据
// kind Agent名称
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

// ReplaceMonitor 替换监控数据
// id 旧数据ID
// scg 监控数据
func ReplaceMonitor(id string, scg interface{}) error {
	logrus.WithFields(logrus.Fields{"id": id, "scg": scg}).Info(ModuleName)
	err := DeleteMonitorByID(id)
	if err != nil {
		return err
	}

	return SaveMonitor(scg)
}

// SaveMonitor 插入一条监控数据
func SaveMonitor(scg interface{}) error {
	return MongoMonitor().Insert(&scg)
}

func DeleteMonitorByID(id string) (error) {
	return MongoMonitor().Remove(bson.M{"_id": bson.ObjectIdHex(id)})
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
