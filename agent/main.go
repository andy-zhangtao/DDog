package main

import (
	"github.com/andy-zhangtao/_hulk_client"
	"github.com/andy-zhangtao/DDog/agent/agents"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/sirupsen/logrus"
	"os"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.
//DDog Agent负责执行DDog分发的任务

const (
	ModuleName = "DDog Agent"
)

var _VERSION_ = "v1.1.0"

func init() {
	_hulk_client.Run()
	if os.Getenv(_const.EnvNsqdEndpoint) == "" {
		logrus.WithFields(logrus.Fields{"Env Empty": _const.EnvNsqdEndpoint,}).Panic(_const.EnvNsqdEndpoint)
	}

}

func main() {

	logrus.WithFields(logrus.Fields{"Version": _VERSION_, "Agent": os.Getenv("DDOG_AGENT_NAME")}).Info(ModuleName)

	switch (os.Getenv("DDOG_AGENT_NAME")) {
	case agents.MonitorAgentName:
		mm := &agents.MonitorAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.MonitorAgentName}
		go mm.Run()
		<-mm.StopChan

	case agents.RetriAgentName:
		ret := &agents.RetriAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.RetriAgentName, Namespace: []string{}}
		go ret.Run()
		<-ret.StopChan

	case agents.SpiderAgentName:
		ret := &agents.SpiderAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), AlivaChan: make(chan int), Name: agents.SpiderAgentName, Namespace: []string{}}
		go ret.Run()
		<-ret.StopChan
	case agents.DeployAgentName:
		ret := &agents.DeployAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.DeployAgentName}
		go ret.Run()
		<-ret.StopChan
	case agents.ReplicaAgentName:
		ret := &agents.ReplicaAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.ReplicaAgentName}
		go ret.Run()
		<-ret.StopChan
	case agents.K8sMonitorAgentName:
		ret := &agents.K8sMonitorAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.K8sMonitorAgentName}
		go ret.Run()
		<-ret.StopChan
	default:
		agn := &agents.AgentNsq{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.DestoryAgent}

		go agn.RunDestoryAgent()
		<-agn.StopChan
	}
}
