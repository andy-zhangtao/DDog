package monitor

import (
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/Sirupsen/logrus"
	"reflect"
)

const (
	ModuleName = "MonitorModule"
)
//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.
type MonitorModule struct {
	Id        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Kind      string        `json:"kind"`
	Svcname   string        `json:"svcname"`
	Namespace string        `json:"namespace"`
	Status    string        `json:"status"`
	Msg       string        `json:"msg"`
	Num       int           `json:"num"`
}

// Save 保存监控信息
// 监控信息中Kind不得为空
// 初始监控状态为NotDeal。如果svcname或者namespace为空，则直接将此消息置为无效
func (mm *MonitorModule) Save() error {
	if mm.Kind == "" {
		return errors.New("Kind Empty!")
	}

	if mm.Status == "" {
		mm.Status = _const.NotDeal
	}

	if mm.Svcname == "" || mm.Namespace == "" {
		mm.Status = _const.DataError
	} else {
		oom, err := mongo.MongoGetMonitorByName(mm.Kind, mm.Svcname, mm.Namespace)
		if err != nil {
			if tool.IsNotFound(err) {
				//没有此类监控数据
				mm.Num++
				return mongo.SaveMonitor(mm)
			}
			return err
		}

		err = nil
		logrus.WithFields(logrus.Fields{"type":reflect.TypeOf(oom),"value":oom}).Info(ModuleName)
		om, err := conver(oom)
		if err != nil {
			return err
		}

		om.Msg = mm.Msg
		om.Num ++
		return mongo.ReplaceMonitor(om.Id.Hex(), om)
	}

	mm.Num++
	return mongo.SaveMonitor(mm)
}

func conver(m interface{}) (mm *MonitorModule, err error) {
	data, err := bson.Marshal(m)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &mm)
	if err != nil {
		return
	}

	return
}
