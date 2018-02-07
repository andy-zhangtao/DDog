package main

import (
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/agent/agents"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.
//DDog Agent负责执行DDog分发的任务

const (
	ModuleName = "DDog Agent"
)

func init() {

	if os.Getenv(_const.EnvNsqdEndpoint) == "" {
		logrus.WithFields(logrus.Fields{"Env Empty": _const.EnvNsqdEndpoint,}).Panic(_const.EnvNsqdEndpoint)
	}

}

func main() {

	logrus.WithFields(logrus.Fields{"Version": "v0.6.4"}).Info(ModuleName)

	switch(os.Getenv("DDOG_AGENT_NAME")) {
	case agents.MonitorAgentName:
		mm := &agents.MonitorAgent{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.MonitorMsgName}
		go mm.Run()
		<-mm.StopChan

	default:
		agn := &agents.AgentNsq{NsqEndpoint: os.Getenv(_const.EnvNsqdEndpoint), StopChan: make(chan int), Name: agents.DestoryAgent}

		go agn.RunDestoryAgent()
		<-agn.StopChan
	}
}
