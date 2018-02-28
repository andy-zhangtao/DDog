package agents

import (
	"github.com/nsqio/go-nsq"
	"github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"strings"
	"strconv"
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

			go this.handlerMsg(&msg)

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

// handlerMsg 处理消息
func (this *MonitorAgent) handlerMsg(msg *monitor.MonitorModule) error {
	switch msg.Kind {
	case RetriAgentName:
		return this.stopSVC(msg)
	case SpiderAgentName:
		return this.confirmSVC(msg)
	}

	return nil
}

// stopSVC 停掉服务并且将其置位失败
func (this *MonitorAgent) stopSVC(msg *monitor.MonitorModule) error {
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		logrus.WithFields(logrus.Fields{MonitorAgentName: "Get MetaData Error!", "error": err}).Error(MonitorAgentName)
		this.StopChan <- 1
	}

	name := ""
	ss := strings.Split(msg.Svcname, "-")
	_, err = strconv.ParseInt(ss[len(ss)-1],10,64)
	if err != nil{
		name = msg.Svcname
	}else{
		/*最后一个字段如果不包含字母，那么就是自动生成的时间戳*/
		name = strings.Join(strings.Split(msg.Svcname, "-")[:len(ss)-2], "-")
	}

	logrus.WithFields(logrus.Fields{"Query Svc": name, "namespace": msg.Namespace}).Info(MonitorAgentName)
	sc, err := svcconf.GetSvcConfByName(name, msg.Namespace)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{"Stop Svc": sc.SvcName, "namespace": msg.Namespace}).Info(MonitorAgentName)
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		Namespace:   msg.Namespace,
		SecretKey:   md.Skey,
		ServiceName: sc.SvcName,
	}

	_, err = q.DeleteService()
	if err != nil {
		return err
	}

	sc.Deploy = _const.DeployFailed
	sc.Msg = msg.Msg
	return svcconf.UpdateSvcConf(sc)
}

// confirmSVC 确定服务状态
// 如果健康检测失败，则将服务置为失败。同时销毁服务
// 如果健康检测成功，则将服务置为成功
func (this *MonitorAgent) confirmSVC(msg *monitor.MonitorModule) error {
	logrus.WithFields(logrus.Fields{"Confirm Svc": msg.Svcname, "namespace": msg.Namespace, "msg": msg.Msg, "ip": msg.Ip}).Info(this.Name)

	if strings.ToLower(msg.Msg) == "ok" {
		sc, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
		if err != nil {
			return err
		}

		if sc.Replicas == len(msg.Ip) {
			/*clear msg ip*/
			msg.Ip = []string{}
			msg.Num = 0
			msg.Replace()
			sc.Deploy = _const.DeploySuc
			sc.Msg = msg.Msg
			return svcconf.UpdateSvcConf(sc)
		}
	} else {
		/*clear msg ip*/
		msg.Ip = []string{}
		msg.Num = 0
		msg.Replace()
		return this.stopSVC(msg)
	}

	return nil
}
