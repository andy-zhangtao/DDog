package agents

import (
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
	"strconv"
	"github.com/drael/GOnetstat"
	"fmt"
	"github.com/andy-zhangtao/DDog/bridge"
	"encoding/json"
	"github.com/mitchellh/go-ps"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/27.

/**
SpiderAgent
服务探针，用来探查服务是否健康。
 */
type SpiderAgent struct {
	Name        string
	NsqEndpoint string
	Namespace   []string
	Port        []int64
	Cmd         string
	StopChan    chan int
}

type SpiderMsg struct {
	Name    string
	Svcname string
	Msg     string
}

func (this *SpiderAgent) Run() {
	logrus.WithFields(logrus.Fields{SpiderAgentName: "Start"}).Info(SpiderAgentName)
	if os.Getenv("DDOG_AGENT_SPIDER_SVC") == "" {
		logrus.WithFields(logrus.Fields{"Env Empty": "DDOG_AGENT_SPIDER_SVC"}).Error(SpiderAgentName)
		this.StopChan <- 0
		return
	}

	if os.Getenv("DDOG_AGENT_SPIDER_PORT") != "" {
		for _, p := range strings.Split(os.Getenv("DDOG_AGENT_SPIDER_PORT"), ";") {
			if p != "" {
				pp, err := strconv.ParseInt(p, 10, 64)
				if err != nil {
					logrus.WithFields(logrus.Fields{"Parse Port Error": err}).Error(SpiderAgentName)
					this.StopChan <- 0
					return
				}
				this.Port = append(this.Port, pp)
			}
		}
		this.checkPort()
	} else if os.Getenv("DDOG_AGENT_SPIDER_CMD") != "" {
		this.Cmd = os.Getenv("DDOG_AGENT_SPIDER_CMD")
		this.checkCmd()
	} else {
		logrus.WithFields(logrus.Fields{"Return": "Nothing To Do"}).Error(SpiderAgentName)
		this.StopChan <- 0
	}
}

func (this *SpiderAgent) checkPort() {

	tcp_data := GOnetstat.Tcp()
	portMap := make(map[int64]int)

	for _, p := range this.Port {
		portMap[p] = 1
	}

	for _, td := range tcp_data {
		if _, ok := portMap[td.Port]; ok {
			delete(portMap, td.Port)
		}

		if len(portMap) == 0 {
			break
		} else {
			msg := ""
			for k, _ := range portMap {
				msg += fmt.Sprintf("Port[%d] Not Listen. ", k)
			}
			sm := SpiderMsg{Name: this.Name, Svcname: os.Getenv("DDOG_AGENT_SPIDER_SVC"), Msg: msg}
			data, err := json.Marshal(&sm)
			if err != nil {
				logrus.WithFields(logrus.Fields{"Marshal Error": err}).Error(SpiderAgentName)
				continue
			}

			bridge.SendMonitorMsg(string(data))
		}
	}
}

func (this *SpiderAgent) checkCmd() {
	porcess, err := ps.Processes()
	if err != nil {
		logrus.WithFields(logrus.Fields{"Get Process Error": err}).Error(SpiderAgentName)
	}

	for _, p := range porcess {
		logrus.WithFields(logrus.Fields{"name": p.Executable()}).Info(SpiderAgentName)
	}
}
