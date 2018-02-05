package server

import (
	"github.com/nsqio/go-nsq"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/svcconf"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

const (
	ModuleName = "DDog-Agent-Nsq"
)

type AgentNsq struct {
	NsqEndpoint string
	StopChan    chan int
}

type DestroyAgent struct{}

var workerChan chan *nsq.Message

func (h *DestroyAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerChan <- m
	return nil
}

func (this *AgentNsq) RunDestoryAgent() {
	workerChan = make(chan *nsq.Message)

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(_const.SvcDestroyMsg, "Agent-"+_const.SvcDestroyMsg, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err}).Error(ModuleName)
		return
	}

	go func() {
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{"Destory Msg": string(m.Body)}).Info(ModuleName)
			msg := _const.DestoryMsg{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Info(ModuleName)
				continue
			}

			oper := svcconf.Operation{}
			err = oper.DeleteSvcConf(msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Delete Service Error":err,"Origin Svc":msg.Svcname, "Origin Namespace":msg.Namespace}).Error(ModuleName)
				continue
			}

			logrus.WithFields(logrus.Fields{"Delete Service":"Success","Origin Svc":msg.Svcname, "Origin Namespace":msg.Namespace}).Info(ModuleName)
			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&DestroyAgent{}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{"Destory Msg": "Listen..."}).Info(ModuleName)
	<-r.StopChan
}
