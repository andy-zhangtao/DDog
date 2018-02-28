package monitor

import (
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/tool"
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
	Ip        []string      `json:"ip"`
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
		//logrus.WithFields(logrus.Fields{"type": reflect.TypeOf(oom), "value": oom}).Info(ModuleName)
		om, err := conver(oom)
		if err != nil {
			return err
		}

		om.Msg = mm.Msg
		om.Num ++

		if len(om.Ip) > 0 {
			isCheck := false
			for _, p := range om.Ip {
				if p == mm.Ip[0] {
					isCheck = true
					break
				}
			}

			if !isCheck {
				om.Ip = append(om.Ip, mm.Ip...)
			}
		} else {
			om.Ip = mm.Ip
		}
		/*Merge MonitorMsg*/
		mm = om
		return mongo.ReplaceMonitor(om.Id.Hex(), om)
	}

	mm.Num++
	return mongo.SaveMonitor(mm)
}

// GetMonitroModule 获取指定类型的监控信息
func GetMonitroModule(kind, svcname, namespace string) (*MonitorModule, error){
	oom, err := mongo.MongoGetMonitorByName(kind, svcname, namespace)
	if err != nil {
		return nil, err
	}

	err = nil
	//logrus.WithFields(logrus.Fields{"type": reflect.TypeOf(oom), "value": oom}).Info(ModuleName)
	return conver(oom)
}

func (mm *MonitorModule) Replace() error{
	oom, err := mongo.MongoGetMonitorByName(mm.Kind, mm.Svcname, mm.Namespace)
	if err != nil {
		if tool.IsNotFound(err) {
			return mongo.SaveMonitor(mm)
		}
		return err
	}

	err = nil
	om, err := conver(oom)
	if err != nil {
		return err
	}

	return mongo.ReplaceMonitor(om.Id.Hex(), mm)
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
