package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/agent"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
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
	producer, _ = nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())

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

			go func() {
				span, reporter, err := tool.GetChildZipKinSpan(DeployAgentName, tool.GetLocalIP(), true, msg.Span)
				if err != nil {
					logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(DeployAgentName)
				} else {
					logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(DeployAgentName)
					span.Tag("service", msg.SvcName)
					span.Tag("namespace", msg.NameSpace)
					span.Tag("upgrade", fmt.Sprintf("%v", msg.Upgrade))
					span.Tag("replicas", fmt.Sprintf("%v", msg.Replicas))
					span.Annotate(time.Now(), "Deploy Agent Server Receive Message")
				}

				var errmessage = ""

				err = this.handlerMsg(&msg)
				if err != nil {
					logrus.WithFields(logrus.Fields{"HandlerMsg Error": err}).Error(this.Name)
					errmessage = err.Error()
				}

				if errmessage != "" {
					span.Annotate(time.Now(), fmt.Sprintf("%s-Deploy Error [%s]", DeployAgentName, errmessage))
				}
				span.Finish()
				reporter.Close()
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

	//删除可能会存在的Monitor Message, 如果不删除, 则在根据IP数量判断服务状态时会永远失败
	mm, err := monitor.GetMonitroModule(SpiderAgentName, msg.SvcName, msg.NameSpace)
	if err == nil && mm.Id != "" {
		mm.Destory()
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

	sc, err := svcconf.GetSvcConfByName(msg.SvcName, msg.NameSpace)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "name": msg.SvcName}).Error(DeployAgentName)
		return err
	}

	if sc == nil {
		err = errors.New(fmt.Sprintf("Can not find svcConf [%s][%s]", msg.SvcName, msg.NameSpace))
		logrus.WithFields(logrus.Fields{"err": err, "name": msg.SvcName}).Error(DeployAgentName)
		return err
	}

	logrus.WithFields(logrus.Fields{"Status": writer.status, "Header": writer.header}).Info(this.Name)
	//发送服务创建信息
	sc.Deploy = _const.DeployIng
	//sc.Msg = msg.SvcName
	NotifyEvent(sc, _const.CREATESERVICE)
	NotifyDevEx(sc)

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
