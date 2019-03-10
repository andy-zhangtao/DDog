package agents

import (
	"encoding/json"
	"os"
	"time"

	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

const (
	ModuleName          = "DDog-Agent-Nsq"
	DestoryAgent        = "DestoryAgent"
	MonitorAgentName    = "MonitorAgent"
	RetriAgentName      = "RetriAgent"
	SpiderAgentName     = "SpiderAgent"
	DeployAgentName     = "DeployAgent"
	ReplicaAgentName    = "ReplicaAgent"
	K8sMonitorAgentName = "K8sMonitorAgent"
)

type AgentNsq struct {
	Name        string
	NsqEndpoint string
	StopChan    chan int
}

type DestroyAgent struct{}

var workerHome map[string]chan *nsq.Message

func init() {
	workerHome = make(map[string]chan *nsq.Message)
}

func checkEnv() {
	for _, e := range []string{_const.EnvMongo, _const.EnvMongoDB} {
		if os.Getenv(e) == "" {
			logrus.WithFields(logrus.Fields{"Env Empty": e}).Panic(ModuleName)
		}
	}
}

func (h *DestroyAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerHome[DestoryAgent] <- m
	return nil
}

func (this *AgentNsq) RunDestoryAgent() {
	checkEnv()
	workerChan := make(chan *nsq.Message)

	workerHome[this.Name] = workerChan

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
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
				data, _ := json.Marshal(monitor.MonitorModule{
					Kind:      DestoryAgent,
					Svcname:   msg.Svcname,
					Namespace: msg.Namespace,
					Msg:       err.Error(),
				})

				m.Finish()
				bridge.SendMonitorMsg(string(data))
				continue
			}

			span, reporter, err := tool.GetChildZipKinSpan(DestoryAgent, tool.GetLocalIP(), true, msg.Span)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(DestoryAgent)
			} else {
				logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(DestoryAgent)
				span.Tag("service", msg.Svcname)
				span.Tag("namespace", msg.Namespace)
				span.Annotate(time.Now(), "Destory Agent Server Receive Message")
			}

			oper := svcconf.Operation{}
			err = oper.DeleteSvcConf(msg,span)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Delete Service Error": err, "Origin Svc": msg.Svcname, "Origin Namespace": msg.Namespace}).Error(ModuleName)
				data, _ := json.Marshal(monitor.MonitorModule{
					Kind:      DestoryAgent,
					Svcname:   msg.Svcname,
					Namespace: msg.Namespace,
					Msg:       err.Error(),
				})

				m.Finish()
				bridge.SendMonitorMsg(string(data))
				continue
			}

			logrus.WithFields(logrus.Fields{"Delete Service": "Success", "Origin Svc": msg.Svcname, "Origin Namespace": msg.Namespace}).Info(ModuleName)
			m.Finish()
			span.Finish()
			reporter.Close()
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
