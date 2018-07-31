package agents

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/DDog/const"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/agent"
	"net/http"
	"fmt"
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"errors"
	"os"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/server/tool"
	"time"
)

/*DeployAgent 部署Agent.
从NSQ接受DDOG发来的部署通知然后执行部署*/
type DeployAgent struct {
	Name        string
	NsqEndpoint string
	StopChan    chan int
}

func (this *DeployAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerHome[DeployAgentName] <- m
	return nil
}

func (this *DeployAgent) Run() {
	workerChan := make(chan *nsq.Message)

	workerHome[this.Name] = workerChan

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000
	r, err := nsq.NewConsumer(_const.SvcDeployMsg, DeployAgentName, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": _const.SvcDeployMsg}).Error(this.Name)
		return
	}

	go func() {
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{_const.SvcDeployMsg: string(m.Body)}).Info(this.Name)
			msg := agent.DeployMsg{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				continue
			}

			logrus.WithFields(logrus.Fields{"SvcName": msg.SvcName, "NameSpace": msg.NameSpace, "Upgrade": msg.Upgrade, "Replicas": msg.Replicas}).Info(this.Name)

			span, reporter, err := tool.GetChildZipKinSpan(DeployAgentName, tool.GetLocalIP(), true, msg.Span)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(DeployAgentName)
			} else {
				logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(DeployAgentName)
				span.Annotate(time.Now(), "Deploy Agent Server Receive Message")
			}

			go func() {
				var errmessage = ""
				defer func() {
					if errmessage != "" {
						span.Annotate(time.Now(), fmt.Sprintf("%s-Deploy Error [%s]", DeployAgentName, errmessage))
					}
					span.Finish()
					reporter.Close()
				}()
				err = this.handlerMsg(&msg)
				if err != nil {
					errmessage = err.Error()
					logrus.WithFields(logrus.Fields{"HandlerMsg Error": err}).Error(this.Name)
				}
			}()

			m.Finish()

		}
	}()

	r.AddConcurrentHandlers(&DeployAgent{}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{DeployAgentName: "Listen...", "Topic": _const.SvcDeployMsg}).Info(this.Name)
	<-r.StopChan
}

// handlerMsg 调用部署API开始部署服务
func (this *DeployAgent) handlerMsg(msg *agent.DeployMsg) error {
	var traceid = ""
	var parentid = ""
	var id = ""
	if msg.Span.TraceID.String() != "" {
		traceid = fmt.Sprintf("&traceid=%s", msg.Span.TraceID.String())
	}
	if msg.Span.ID.String() != "" {
		id = fmt.Sprintf("&id=%s", msg.Span.ID.String())
	}
	if msg.Span.ParentID.String() != "" {
		parentid = fmt.Sprintf("&parentid=%s", msg.Span.ParentID.String())
	}

	logrus.WithFields(logrus.Fields{"url": fmt.Sprintf("/v1/cloud/svc/deploy?svcname=%s&namespace=%s&upgrade=%v%s%s%s", msg.SvcName, msg.NameSpace, msg.Upgrade, traceid, id, parentid)}).Info(this.Name)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/cloud/svc/deploy?svcname=%s&namespace=%s&upgrade=%v%s%s%s", msg.SvcName, msg.NameSpace, msg.Upgrade, traceid, id, parentid), nil)
	if err != nil {
		return err
	}

	writer := MyWriter{
		header: make(map[string][]string),
	}

	qcloud.RunService(&writer, r)

	logrus.WithFields(logrus.Fields{"Status": writer.status, "Header": writer.header}).Info(this.Name)
	if writer.status != http.StatusOK && writer.status != 0 {
		logrus.WithFields(logrus.Fields{"Deploy Service Error": writer.data}).Error(this.Name)
		m := monitor.MonitorModule{
			Kind:      RetriAgentName,
			Svcname:   msg.SvcName,
			Namespace: msg.NameSpace,
		}

		data, err := json.Marshal(&m)
		if err != nil {
			return err
		}
		producer, _ := nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())
		producer.Publish(_const.SvcMonitorMsg, data)
		return errors.New(writer.data)
	}
	return nil
}

type MyWriter struct {
	header http.Header
	data   string
	status int
}

func (this *MyWriter) Header() http.Header {
	return this.header
}

func (this *MyWriter) Write(d []byte) (int, error) {
	this.data = string(d)
	return len(this.data), nil
}

func (this *MyWriter) WriteHeader(h int) {
	this.status = h
	logrus.WithFields(logrus.Fields{"header": h, "status": this.status}).Info("DeployAgent-Writer")
}
