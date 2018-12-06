package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/k8service"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/nsqio/go-nsq"
	"github.com/openzipkin/zipkin-go"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

//调用K8s APIServer API获取服务状态数据

type K8sMonitorAgent struct {
	Name        string
	NsqEndpoint string
	StopChan    chan int
}

func (k *K8sMonitorAgent) HandleMessage(m *nsq.Message) error {
	m.DisableAutoResponse()
	workerHome[K8sMonitorAgentName] <- m
	return nil
}

var currentDeploySvc map[string]int

const (
	INSTANCE_INIT = iota
	INSTANCE_LOOKUP
)

func (this *K8sMonitorAgent) Run() {
	workerChan := make(chan *nsq.Message)

	workerHome[K8sMonitorAgentName] = workerChan
	producer, _ = nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())

	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 1000

	topic := fmt.Sprintf("%s_%s", _const.SvcK8sMonitorMsg, os.Getenv(_const.ENV_WATCH_MONITOR_NAMESPACE))
	r, err := nsq.NewConsumer(topic, K8sMonitorAgentName, cfg)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Create Consumer Error": err, "Agent": K8sMonitorAgentName, "topic": topic}).Error(this.Name)
		return
	}

	k8sMasters, err := k8service.GetALlK8sCluster()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Query K8s API Server Error": err}).Error(K8sMonitorAgentName)
		return
	}

	currentDeploySvc = make(map[string]int)

	logrus.WithFields(logrus.Fields{"K8s API Service": k8sMasters, "Watch-Namespace": os.Getenv(_const.ENV_WATCH_MONITOR_NAMESPACE)}).Info(K8sMonitorAgentName)

	go func() {
		for m := range workerChan {
			logrus.WithFields(logrus.Fields{topic: string(m.Body)}).Info(K8sMonitorAgentName)
			msg := monitor.MonitorModule{}

			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Unmarshal Msg": err, "Origin Byte": string(m.Body)}).Error(this.Name)
				m.Finish()
				continue
			}

			span, reporter, err := tool.GetChildZipKinSpan(K8sMonitorAgentName, tool.GetLocalIP(), true, msg.Span)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Get ZipKin Span Error": err}).Error(K8sMonitorAgentName)
			} else {
				logrus.WithFields(logrus.Fields{"span": span.Context()}).Info(K8sMonitorAgentName)
				span.Annotate(time.Now(), fmt.Sprintf("%s Receive Message", K8sMonitorAgentName))
				span.Tag("service", msg.Svcname)
				span.Tag("namesapce", msg.Namespace)
				span.Tag("msg", msg.Msg)
				span.Tag("kind", msg.Kind)
				span.Tag("ip", fmt.Sprintf("%v", msg.Ip))
			}

			logrus.WithFields(logrus.Fields{"Kind": msg.Kind, "Svc": msg.Svcname, "Namespace": msg.Namespace, "ip": msg.Ip, "msg": msg.Msg}).Info(K8sMonitorAgentName)
			isStand := false
			var apiServer k8sconfig.K8sCluster
			for _, a := range k8sMasters {
				if a.Namespace == msg.Namespace {
					isStand = true
					apiServer = a
					break
				}
			}

			if !isStand {
				//	使用默认的K8s集群数据
				for _, a := range k8sMasters {
					if a.Namespace == "devenv" {
						apiServer = a
						apiServer.Namespace = msg.Namespace
						break
					}
				}
			}

			logrus.WithFields(logrus.Fields{"isStand": isStand, "api": apiServer.Name}).Info(K8sMonitorAgentName)
			go this.handlerMsg(apiServer, &msg, &span)
			span.Finish()
			reporter.Close()
			m.Finish()
		}
	}()

	r.AddConcurrentHandlers(&K8sMonitorAgent{}, 20)

	err = r.ConnectToNSQD(this.NsqEndpoint)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	logrus.WithFields(logrus.Fields{K8sMonitorAgentName: "Listen...", "Topic": _const.SvcK8sMonitorMsg}).Info(this.Name)
	<-r.StopChan
}

func (this *K8sMonitorAgent) handlerMsg(apiServer k8sconfig.K8sCluster, msg *monitor.MonitorModule, span *zipkin.Span) {
	logrus.WithFields(logrus.Fields{"msg": msg.Msg, "name": msg.Svcname}).Info(K8sMonitorAgentName)
	sc, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "name": msg.Svcname}).Error(K8sMonitorAgentName)
		return
	}

	if sc == nil {
		err = errors.New(fmt.Sprintf("Can not find svcConf [%s][%s]", msg.Svcname, msg.Namespace))
		logrus.WithFields(logrus.Fields{"err": err, "name": msg.Svcname}).Error(K8sMonitorAgentName)
		(*span).Annotate(time.Now(), err.Error())
		return
	}

	switch strings.ToLower(msg.Msg) {
	case "deploy":
		//	接受到部署消息
		if _, ok := currentDeploySvc[msg.Svcname]; !ok {
			//只添加新增的服务
			currentDeploySvc[msg.Svcname] = INSTANCE_INIT
			//	发送服务创建事件，同时通知Devex修改状态
			sc.Deploy = _const.HealthCheck
			//sc.Msg = sc.SvcName
			//sc.Span = (*span).Context()
			//NotifyEvent(sc, _const.CREATESERVICE)
			NotifyDevEx(sc)
		}
	case "ok":
		//	实例部署成功
		//sc, err := svcconf.GetSvcConfByName(msg.Svcname, msg.Namespace)
		//if err != nil {
		//	logrus.WithFields(logrus.Fields{"err": err, "name": msg.Svcname}).Error(K8sMonitorAgentName)
		//	return
		//}
		//
		//if sc == nil {
		//	err = errors.New(fmt.Sprintf("Can not find svcConf [%s][%s]", msg.Svcname, msg.Namespace))
		//	logrus.WithFields(logrus.Fields{"err": err, "name": msg.Svcname}).Error(K8sMonitorAgentName)
		//	(*span).Annotate(time.Now(), err.Error())
		//	return
		//}

		logrus.WithFields(logrus.Fields{"name": sc.SvcName, "namespace": sc.Namespace, "stat": currentDeploySvc[msg.Svcname]}).Info(ModuleName)
		if currentDeploySvc[msg.Svcname] == INSTANCE_INIT {
			currentDeploySvc[msg.Svcname] = INSTANCE_LOOKUP

			//	开始轮询当前实例状态
			//	Spider上送的是统一服务名称，这里需要的实际服务名称
			msg.Svcname = sc.SvcName
			for {
				isReady, err := checkServiceStata(apiServer, msg, span)
				if err != nil {
					(*span).Annotate(time.Now(), fmt.Sprintf("Service Status Check Error %s ", err.Error()))
					return
				}

				logrus.WithFields(logrus.Fields{"name": sc.SvcName, "namespace": sc.Namespace, "isReady": isReady}).Info(ModuleName)
				if isReady {
					ip, port, err := getServiceLB(apiServer, msg, span)
					if err != nil {
						(*span).Annotate(time.Now(), fmt.Sprintf("Service LB Check Error %s ", err.Error()))
						return
					}

					lb := svcconf.LoadBlance{
						IP: ip,
					}

					lb.Port = port
					sc.LbConfig = lb
					sc.Deploy = _const.DeploySuc
					sc.Status = _const.NeedDeploy
					sc.Msg = msg.Msg
					sc.Span = (*span).Context()
					logrus.WithFields(logrus.Fields{"name": sc.SvcName, "namespace": sc.Namespace, "lb": lb}).Info(ModuleName)
					NotifyDevEx(sc)

					break
				}

				time.Sleep(3 * time.Second)
			}
			delete(currentDeploySvc, sc.Name)
		}

	}
}

func checkServiceStata(apiServer k8sconfig.K8sCluster, msg *monitor.MonitorModule, span *zipkin.Span) (isReady bool, err error) {
	isReady = false
	deploy, err := k8service.GetK8sSpecifyDeployMent(apiServer, msg.Namespace, msg.Svcname)
	if err != nil {
		return
	}

	(*span).Annotate(time.Now(), fmt.Sprintf("Repllicas [%v] ReadyReplicas [%v]", deploy.Status.Replicas, deploy.Status.ReadyReplicas))
	logrus.WithFields(logrus.Fields{"Repllicas": deploy.Status.Replicas, "ReadyReplicas": deploy.Status.ReadyReplicas, "UpdatedReplicas": deploy.Status.UpdatedReplicas, "name": deploy.Metadata.Name}).Info(ModuleName)
	if (deploy.Status.Replicas == deploy.Status.ReadyReplicas) && (deploy.Status.ReadyReplicas == deploy.Status.UpdatedReplicas) {
		return true, nil
	}

	return
}

func getServiceLB(apiServer k8sconfig.K8sCluster, msg *monitor.MonitorModule, span *zipkin.Span) (ip string, port []int, err error) {

	for {
		service, err := k8service.GetK8sSpecifyService(apiServer, msg.Namespace, msg.Svcname)
		if err != nil {
			return ip, port, err
		}

		(*span).Annotate(time.Now(), fmt.Sprintf("Service LB [%v]", service.Status.LoadBalancer.Ingress))

		if len(service.Status.LoadBalancer.Ingress) > 0 {
			switch msg.Namespace {
			case _const.PROENV:
				fallthrough
			case _const.RELEASEENV:
				if !strings.HasPrefix(service.Status.LoadBalancer.Ingress[0].IP, "10.0.") {
					time.Sleep(3 * time.Second)
					continue
				}

				ip = service.Status.LoadBalancer.Ingress[0].IP
				for _, p := range service.Spec.Ports {
					port = append(port, p.Port)
				}

				return ip, port, err
			case _const.DEVENV:
				fallthrough
			case _const.TESTENV:
				fallthrough
			default:
				//	开发和测试环境，IP属于192.168.0.0/16网段
				if !strings.HasPrefix(service.Status.LoadBalancer.Ingress[0].IP, "192.168.") {
					time.Sleep(3 * time.Second)
					continue
				}
				ip = service.Status.LoadBalancer.Ingress[0].IP
				for _, p := range service.Spec.Ports {
					port = append(port, p.Port)
				}
				return ip, port, err
			}
		} else {
			time.Sleep(3 * time.Second)
			continue
		}
	}

}
