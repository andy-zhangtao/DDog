package agents

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"strings"
	"strconv"
	"time"
	"os"
	"fmt"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.
// 采集失败任务名称，并择机重新执行

type MonitorAgent struct {
	Name        string
	NsqEndpoint string
	StopChan    chan int
}

var producer *nsq.Producer

func (h *MonitorAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerHome[MonitorAgentName] <- m
	return nil
}

func (this *MonitorAgent) Run() {
	workerChan := make(chan *nsq.Message)

	workerHome[this.Name] = workerChan
	producer, _ = nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())
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

			logrus.WithFields(logrus.Fields{"HandlerMsg svc": msg.Svcname, "namespace": msg.Namespace, "msg": msg.Msg, "ip": msg.Ip}).Info(this.Name)
			go func() {
				err = this.handlerMsg(&msg)
				if err != nil {
					logrus.WithFields(logrus.Fields{"HandlerMsg Error": err}).Error(ModuleName)
				}
			}()

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
	_, err = strconv.ParseInt(ss[len(ss)-1], 10, 64)
	if err != nil {
		name = msg.Svcname
	} else {
		/*最后一个字段如果不包含字母，那么就是自动生成的时间戳*/
		name = strings.Join(strings.Split(msg.Svcname, "-")[:len(ss)-2], "-")
	}

	logrus.WithFields(logrus.Fields{"Query Svc": name, "namespace": msg.Namespace}).Info(MonitorAgentName)
	sc, err := svcconf.GetSvcConfByName(name, msg.Namespace)
	if err != nil {
		return err
	}

	if sc == nil {
		logrus.WithFields(logrus.Fields{"Not Found SVC": name, "namespace": msg.Namespace}).Info(MonitorAgentName)
		return nil
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
	//sc.SvcName = ""
	this.NotifyDevEx(sc)
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

		mm, err := monitor.GetMonitroModule(msg.Kind, msg.Svcname, msg.Namespace)
		if err != nil {
			return err
		}

		if sc.Replicas == len(mm.Ip) {
			/*clear msg*/
			mm.Destory()
			/*获取负载信息*/
			md, err := metadata.GetMetaDataByRegion("")
			if err != nil {
				sc.Deploy = _const.DeploySuc
				sc.Msg = msg.Msg
				return svcconf.UpdateSvcConf(sc)
			}

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
			q.SetDebug(true)

			for {

				resp, err := q.QuerySvcInfo()

				if err != nil || resp.Code != 0 {
					sc.Deploy = _const.DeploySuc
					sc.Msg = msg.Msg
					return svcconf.UpdateSvcConf(sc)
				}

				if strings.ToLower(resp.Data.ServiceInfo.Status) == "normal" {
					lip := ""
					if resp.Data.ServiceInfo.ExternalIp == "" {
						lip = resp.Data.ServiceInfo.ServiceIp
					} else {
						lip = resp.Data.ServiceInfo.ExternalIp
					}
					lb := svcconf.LoadBlance{
						IP: lip,
					}

					var port []int
					for _, c := range resp.Data.ServiceInfo.PortMappings {
						port = append(port, c.LbPort)
					}
					lb.Port = port
					sc.LbConfig = lb
					sc.Deploy = _const.DeploySuc
					sc.Msg = msg.Msg

					this.NotifyDevEx(sc)
					return svcconf.UpdateSvcConf(sc)
				}
				time.Sleep(3 * time.Second)
			}
		}
	} else {
		/*clear msg*/
		msg.Destory()
		return this.stopSVC(msg)
	}

	return nil
}

// NotifyDevEx 通知Devex更新状态
func (this *MonitorAgent) NotifyDevEx(scf *svcconf.SvcConf) {
	var lb []string
	for _, port := range scf.LbConfig.Port {
		lb = append(lb, fmt.Sprintf("%s:%d", scf.LbConfig.IP, port))
	}
	req := struct {
		ProjectID   string   `json:"project_id"`
		Stage       int      `json:"stage"`
		DeployEnv   string   `json:"deploy_env"`
		LoadBalance []string `json:"load_balance"`
	}{
		scf.Name,
		scf.Deploy,
		scf.Namespace,
		lb,
	}

	data, err := json.Marshal(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Marshal DevEx Request Error": err}).Error(this.Name)
		return
	}

	err = producer.Publish("DevEx-Request-Status", data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Notify DevEx Error": err}).Error(this.Name)
		return
	}

	return
}
