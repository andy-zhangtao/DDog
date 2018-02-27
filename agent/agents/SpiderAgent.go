package agents

import (
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
	"strconv"
	"github.com/drael/GOnetstat"
	"fmt"
	"github.com/mitchellh/go-ps"
	"time"
	"github.com/andy-zhangtao/DDog/bridge"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/monitor"
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
	AlivaChan   chan int
}

type SpiderMsg struct {
	Name    string
	Svcname string
	Msg     string
}

// 只有当目标服务处于活动状态的时候才需要开始检测,
// 如果目标服务被Kill掉，则可以认为K8s健康检测失败
// 此时Spider就不需要再检查，开始上报检查结果
var needCheck bool
var msg = ""

func (this *SpiderAgent) Run() {
	logrus.WithFields(logrus.Fields{SpiderAgentName: "Start"}).Info(SpiderAgentName)
	needCheck = false
	go func() {
		for {
			now := time.Now()
			next := now.Add(time.Second * 1)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			logrus.WithFields(logrus.Fields{"下次采集进程时间为": next.Format("200601021504")}).Info(SpiderAgentName)

			select {
			case <-t.C:
				porcess, err := ps.Processes()
				if err != nil {
					logrus.WithFields(logrus.Fields{"Get Process Error": err}).Error(SpiderAgentName)
				}
				logrus.WithFields(logrus.Fields{"Porcess": porcess}).Info(SpiderAgentName)
				this.AlivaChan <- len(porcess)
			}
		}
	}()

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
		logrus.WithFields(logrus.Fields{"Check Port": this.Port}).Info(SpiderAgentName)
		for {
			//now := time.Now()
			//next := now.Add(time.Second * 3)
			//next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
			//t := time.NewTimer(next.Sub(now))
			//logrus.WithFields(logrus.Fields{"needCheck": needCheck, "port": this.Port, "下次端口采集时间为": next.Format("200601021504")}).Info(SpiderAgentName)

			select {
			//case <-t.C:
			//if needCheck {
			//
			//}
			case p := <-this.AlivaChan:
				if p <= 1 {
					needCheck = false
				} else {
					needCheck = true
					this.checkPort()
				}

				if !needCheck && msg != "" {
					data, _ := json.Marshal(monitor.MonitorModule{
						Kind:      this.Name,
						Svcname:   os.Getenv("DDOG_AGENT_SPIDER_SVC"),
						Namespace: os.Getenv("DDOG_AGENT_SPIDER_NS"),
						Msg:       msg,
					})
					//sm := SpiderMsg{Name: this.Name, Svcname: , Msg: msg}
					//data, err := json.Marshal(&sm)
					//if err != nil {
					//	logrus.WithFields(logrus.Fields{"Marshal Error": err}).Error(SpiderAgentName)
					//}

					bridge.SendMonitorMsg(string(data))
					msg = ""
				}
			}
		}
	} else if os.Getenv("DDOG_AGENT_SPIDER_CMD") != "" {
		this.Cmd = os.Getenv("DDOG_AGENT_SPIDER_CMD")
		for {
			now := time.Now()
			next := now.Add(time.Second * 5)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			logrus.WithFields(logrus.Fields{"下次命令采集时间为": next.Format("200601021504")}).Info(SpiderAgentName)

			select {
			case <-t.C:
				this.checkCmd()
			}
		}

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
		logrus.WithFields(logrus.Fields{"tcp": td}).Info(SpiderAgentName)
		if _, ok := portMap[td.Port]; ok {
			delete(portMap, td.Port)
		}
	}

	if len(portMap) != 0 {
		for k, _ := range portMap {
			msg = fmt.Sprintf("Port[%d] Check Failed. ", k)
		}
		logrus.WithFields(logrus.Fields{"msg": msg}).Info(SpiderAgentName)
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

//func (this *SpiderAgent) isAlive() bool {
//	for {
//		now := time.Now()
//		next := now.Add(time.Second * 1)
//		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
//		t := time.NewTimer(next.Sub(now))
//		logrus.WithFields(logrus.Fields{"下次采集进程时间为": next.Format("200601021504")}).Info(SpiderAgentName)
//
//		select {
//		case <-t.C:
//			porcess, err := ps.Processes()
//			if err != nil {
//				logrus.WithFields(logrus.Fields{"Get Process Error": err}).Error(SpiderAgentName)
//			}
//
//			logrus.WithFields(logrus.Fields{"Porcess": porcess}).Info(SpiderAgentName)
//			if len(porcess) > 1 {
//				return true
//			}
//		}
//	}
//}
