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
			//err = this.distMsg(&msg)
			//if err != nil {
			//	logrus.WithFields(logrus.Fields{"Save Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
			//	continue
			//}

			go func() {
				err = this.handlerMsg(&msg)
				if err != nil {
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
	logrus.WithFields(logrus.Fields{"url": fmt.Sprintf("/v1/cloud/svc/deploy?svcname=%s&namespace=%s&upgrade=%v", msg.SvcName, msg.NameSpace, msg.Upgrade)}).Info(this.Name)
	r, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/cloud/svc/deploy?svcname=%s&namespace=%s&upgrade=%v", msg.SvcName, msg.NameSpace, msg.Upgrade), nil)
	if err != nil {
		return err
	}

	writer := MyWriter{
		header: make(map[string][]string),
	}

	qcloud.RunService(&writer, r)

	logrus.WithFields(logrus.Fields{"Status": writer.status, "Header": writer.header}).Info(this.Name)
	if writer.status != http.StatusOK || writer.status != 0 {
		logrus.WithFields(logrus.Fields{"Deploy Service Error": writer.data}).Error(this.Name)
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
