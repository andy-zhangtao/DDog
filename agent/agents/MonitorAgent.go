package agents

import (
	"github.com/nsqio/go-nsq"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/const"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.
//采集失败任务名称，并择机重新执行

type MonitorAgent struct {
	Name        string
	NsqEndpoint string
	StopChan    chan int
}

func (h *MonitorAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerHome[MonitorAgentName] <- m
	return nil
}

func (this *MonitorAgent) Run() {
	workerChan := make(chan *nsq.Message)

	workerHome[this.Name] = workerChan

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(_const.SvcMonitorMsg, MonitorAgentName, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": _const.SvcMonitorMsg}).Error(ModuleName)
		return
	}

	go func() {
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{_const.SvcMonitorMsg: string(m.Body)}).Info(ModuleName)
			msg := monitor.MonitorModule{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
				continue
			}

			logrus.WithFields(logrus.Fields{"Kind": msg.Kind, "Origin Svc": msg.Svcname, "Origin Namespace": msg.Namespace}).Info(ModuleName)
			err = this.distMsg(&msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Save Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
				continue
			}

			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&MonitorAgent{}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{MonitorAgentName: "Listen..."}).Info(ModuleName)
	<-r.StopChan
}

func SendMonitor(msg []byte) {
	logrus.WithFields(logrus.Fields{"Monitor Msg": string(msg)}).Info(ModuleName)

}

// distMsg 消息分发
func (this *MonitorAgent) distMsg(msg *monitor.MonitorModule) error {
	return msg.Save()
}
