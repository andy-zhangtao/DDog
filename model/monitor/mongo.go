package monitor

import (
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.
type MonitorModule struct {
	Kind      string `json:"kind"`
	Svcname   string `json:"svcname"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Msg       string `json:"msg"`
	Num       int    `json:"num"`
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
	}

	mm.Num++
	return mongo.SaveMonitor(mm)
}
