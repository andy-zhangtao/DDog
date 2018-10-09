package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/nsqio/go-nsq"
	"github.com/openzipkin/zipkin-go"
	zmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
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
			defer m.Finish()

			var errmessage = ""
			logrus.WithFields(logrus.Fields{_const.SvcMonitorMsg: string(m.Body)}).Info(ModuleName)
			msg := monitor.MonitorModule{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				errmessage = fmt.Sprintf("Unmarshal Msg [%s] Origin Byte [%s]", err.Error(), string(m.Body))
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
				//m.Finish()
				continue
			}

			span, reporter, err := tool.GetChildZipKinSpan(MonitorAgentName, tool.GetLocalIP(), true, msg.Span)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(MonitorAgentName)
			} else {
				logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(MonitorAgentName)
				span.Annotate(time.Now(), fmt.Sprintf("%s Receive Message", MonitorAgentName))
				span.Tag("service", msg.Svcname)
				span.Tag("namesapce", msg.Namespace)
				span.Tag("msg", msg.Msg)
				span.Tag("kind", msg.Kind)
				span.Tag("ip", fmt.Sprintf("%v", msg.Ip))
			}

			defer func() {
				span.Finish()
				reporter.Close()
			}()

			logrus.WithFields(logrus.Fields{"Kind": msg.Kind, "Origin Svc": msg.Svcname, "Origin Namespace": msg.Namespace, "ip": msg.Ip, "msg": msg.Msg}).Info(ModuleName)
			err = this.distMsg(&msg)
			if err != nil {
				errmessage = fmt.Sprintf("Save Msg Error [%s] Origin Byte [%s]", err.Error(), string(m.Body))
				logrus.WithFields(logrus.Fields{"Save Msg": err, "Origin Byte": string(m.Body)}).Error(ModuleName)
				//m.Finish()
				continue
			}

			logrus.WithFields(logrus.Fields{"HandlerMsg svc": msg.Svcname, "namespace": msg.Namespace, "msg": msg.Msg, "ip": msg.Ip}).Info(this.Name)
			//go func() {
			//	defer func() {
			//		if errmessage != "" {
			//			span.Annotate(time.Now(), fmt.Sprintf("%s Error [%s]", MonitorAgentName, errmessage))
			//		}
			//
			//		span.Finish()
			//		reporter.Close()
			//	}()
			//	err = this.handlerMsg(&msg, span)
			//	if err != nil {
			//		errmessage = err.Error()
			//		logrus.WithFields(logrus.Fields{"HandlerMsg Error": err}).Error(ModuleName)
			//	}
			//}()

			err = this.handlerMsg(&msg, span)
			if err != nil {
				errmessage = err.Error()
				logrus.WithFields(logrus.Fields{"HandlerMsg Error": err}).Error(ModuleName)
			}

			if errmessage != "" {
				span.Annotate(time.Now(), fmt.Sprintf("%s Error [%s]", MonitorAgentName, errmessage))
			}

			//span.Finish()
			//reporter.Close()
			//m.Finish()
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
func (this *MonitorAgent) handlerMsg(msg *monitor.MonitorModule, span zipkin.Span) error {
	switch msg.Kind {
	case RetriAgentName:
		return this.stopSVC(msg, span)
	case SpiderAgentName:
		return this.confirmSVC(msg, span)
	}

	return nil
}

// stopSVC 停掉服务并且将其置位失败
func (this *MonitorAgent) stopSVC(msg *monitor.MonitorModule, span zipkin.Span) error {
	var md *metadata.MetaData
	var err error
	switch msg.Namespace {
	case "proenv":
		fallthrough
	case "release":
		md, err = metadata.GetMetaDataByRegion("", msg.Namespace)
	default:
		md, err = metadata.GetMetaDataByRegion("")
	}
	//md, err := metadata.GetMetaDataByRegion("")
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
	sc.Span = span.Context()
	//sc.SvcName = ""
	this.NotifyDevEx(sc)
	return svcconf.UpdateSvcConf(sc)
}

// confirmSVC 确定服务状态
// 如果健康检测失败，则将服务置为失败。同时销毁服务
// 如果健康检测成功，则将服务置为成功
func (this *MonitorAgent) confirmSVC(msg *monitor.MonitorModule, span zipkin.Span) error {
	logrus.WithFields(logrus.Fields{"Confirm Svc": msg.Svcname, "namespace": msg.Namespace, "msg": msg.Msg, "ip": msg.Ip}).Info(this.Name)

	if strings.ToLower(msg.Msg) == "ok" {
		sc, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
		if err != nil {
			return err
		}

		if sc == nil {
			return errors.New(fmt.Sprintf("Can not find svcConf [%s][%s]", msg.Svcname, msg.Namespace))
		}

		if sc.Status == _const.ModifyReplica {
			logrus.WithFields(logrus.Fields{"ID": sc.SvcName, "Status": "Modify Replica"}).Info(ModuleName)
			return nil
		}

		mm, err := monitor.GetMonitroModule(msg.Kind, msg.Svcname, msg.Namespace)
		if err != nil {
			return err
		}

		logrus.WithFields(logrus.Fields{"service": msg.Svcname, "value": mm}).Info(ModuleName)
		if sc.Replicas == len(mm.Ip) {
			/*clear msg*/
			mm.Destory()
			/*获取负载信息*/
			var md *metadata.MetaData
			switch msg.Namespace {
			case "proenv":
				fallthrough
			case "release":
				md, err = metadata.GetMetaDataByRegion("", msg.Namespace)
			default:
				md, err = metadata.GetMetaDataByRegion("", )
			}

			//if msg.Namespace == "proenv" {
			//	md, err = metadata.GetMetaDataByRegion("", "proenv")
			//} else {
			//	md, err = metadata.GetMetaDataByRegion("")
			//}

			if err != nil {
				sc.Deploy = _const.DeployFailed
				sc.Msg = msg.Msg
				return svcconf.UpdateSvcConf(sc)
			}

			go func(span zipkin.Span) {
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
					var errmessage = ""
					span, reporter, err := tool.GetChildZipKinSpan(MonitorAgentName+"-Qcloud-Query-Service", tool.GetLocalIP(), true, span.Context())
					if err != nil {
						logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(MonitorAgentName)
					} else {
						//logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(MonitorAgentName)
						span.Annotate(time.Now(), fmt.Sprintf("%s Query Service Status", MonitorAgentName))
						span.Tag("service", sc.SvcName)
						span.Tag("namesapce", msg.Namespace)
						span.Tag("key", md.Skey)
						span.Tag("slusterID", md.ClusterID)
						span.Tag("secretid", md.Sid)
					}

					logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(ModuleName)
					defer func() {
						span.Finish()
						reporter.Close()
					}()

					resp, err := q.QuerySvcInfo()

					if err != nil {
						if resp.Code != 0 {
							logrus.WithFields(logrus.Fields{"RespCode": resp.Code, "Desc": resp.CodeDesc, "service": sc.SvcName}).Error(ModuleName)
							errmessage = fmt.Sprintf("RespCode [%v] Desc [%v]", resp.Code, resp.CodeDesc)
						} else {
							errmessage = err.Error()
							sc.Deploy = _const.DeployFailed
							sc.Msg = err.Error()
							if err := svcconf.UpdateSvcConf(sc); err != nil {
								errmessage = err.Error()
							}
						}

						//span.Finish()
						//reporter.Close()
						break
					}

					span.Annotate(time.Now(), strings.ToLower(resp.Data.ServiceInfo.Status))
					span.Annotate(time.Now(), resp.Data.ServiceInfo.ExternalIp)
					if strings.ToLower(resp.Data.ServiceInfo.Status) == "normal" && resp.Data.ServiceInfo.ExternalIp != "" {
						switch msg.Namespace {
						case "proenv":
							fallthrough
						case "release":
							//	预发布和正式环境，IP必须属于10.0.0.0/16网段
							if !strings.HasPrefix(resp.Data.ServiceInfo.ExternalIp, "10.0.") {
								time.Sleep(3 * time.Second)
								continue
							}
						case "devenv":
							fallthrough
						case "testenv":
							//	预发布和正式环境，IP必须属于192.168.0.0/16网段
							if !strings.HasPrefix(resp.Data.ServiceInfo.ExternalIp, "192.168.") {
								time.Sleep(3 * time.Second)
								continue
							}
						default:
							md, err = metadata.GetMetaDataByRegion("", )
						}
						lip := resp.Data.ServiceInfo.ExternalIp
						//if resp.Data.ServiceInfo.ExternalIp == "" {
						//	lip = resp.Data.ServiceInfo.ServiceIp
						//} else {
						//
						//}

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
						sc.Span = span.Context()
						this.NotifyDevEx(sc)
						if err := svcconf.UpdateSvcConf(sc); err != nil {
							errmessage = err.Error()
						}
						//span.Finish()
						//reporter.Close()
						break
					}
					if errmessage != "" {
						span.Annotate(time.Now(), errmessage)
					}

					//span.Finish()
					//reporter.Close()
					time.Sleep(5 * time.Second)
				}
			}(span)

		} else {
			return errors.New(fmt.Sprintf("sc.Replicas=[%v] len(mm.Ip)=[%v]", sc.Replicas, len(mm.Ip)))
		}
	} else if strings.ToLower(msg.Msg) == "deploy" {
		//	开始部署服务
		sc, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
		if err != nil {
			return err
		}
		if sc == nil {
			logrus.WithFields(logrus.Fields{"SvcConf": "nil"}).Error(ModuleName)
			//msg.Destory()
			return errors.New("SvcConf nil")
		}
		ns, reporter, err := tool.GetChildZipKinSpan(MonitorAgentName+"-Status-Service", tool.GetLocalIP(), true, span.Context())
		if err != nil {
			return errors.New(fmt.Sprintf("Generate Span Error [%v]", err))
		}

		logrus.WithFields(logrus.Fields{"Span ID": ns.Context().ID.String(), "Traceid": ns.Context().TraceID.String(), "Parentid": ns.Context().ParentID.String(), "P-Span ID": span.Context().ID, "P-Traceid": span.Context().TraceID, "P-Parentid": span.Context().ParentID}).Info(ModuleName)
		ns.Annotate(time.Now(), fmt.Sprintf("%s Query Service Status", MonitorAgentName))
		ns.Tag("service", sc.SvcName)
		ns.Tag("namesapce", msg.Namespace)
		ns.Tag("status", "Deploying")
		ns.Finish()
		reporter.Close()

		sc.Deploy = _const.DeployIng
		sc.Span = ns.Context()
		this.NotifyDevEx(sc)
	} else if strings.ToLower(msg.Msg) == "status" {
		ns, reporter, err := tool.GetChildZipKinSpan(MonitorAgentName+"-Check-Result", tool.GetLocalIP(), true, span.Context())
		if err != nil {
			return errors.New(fmt.Sprintf("Generate Span Error [%v]", err))
		}

		logrus.WithFields(logrus.Fields{"Span ID": ns.Context().ID.String(), "Traceid": ns.Context().TraceID.String(), "Parentid": ns.Context().ParentID.String(), "P-Span ID": span.Context().ID, "P-Traceid": span.Context().TraceID, "P-Parentid": span.Context().ParentID}).Info(ModuleName)
		ns.Annotate(time.Now(), fmt.Sprintf("%s Query Service Status", MonitorAgentName))
		ns.Tag("check result", msg.Status)
		ns.Finish()
		reporter.Close()
	} else {
		/*clear msg*/
		//msg.Destory()
		return this.stopSVC(msg, span)
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
		ProjectID   string             `json:"project_id"`
		Stage       int                `json:"stage"`
		DeployEnv   string             `json:"deploy_env"`
		LoadBalance []string           `json:"load_balance"`
		Span        zmodel.SpanContext `json:"span"`
	}{
		scf.Name,
		scf.Deploy,
		scf.Namespace,
		lb,
		scf.Span,
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
