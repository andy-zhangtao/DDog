package caasmodel

import "gopkg.in/mgo.v2/bson"

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/18. 

// NameSpace 保存Caas环境中所有的Namespace数据
type NameSpace struct {
	ID         bson.ObjectId `json:"_id" bson:"_id"`
	Name       string        `json:"name"`
	Owner      string        `json:"owner"`
	CreateTime string        `json:"create_time"`
	Desc       string        `json:"desc"`
}
