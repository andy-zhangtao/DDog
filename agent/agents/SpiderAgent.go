package agents

import (
	"encoding/json"
	"fmt"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"github.com/drael/GOnetstat"
	"github.com/mitchellh/go-ps"
	zmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

// Write by zhangtao<ztao8607@gmail.com> . In 2018/2/27.

/**
SpiderAgent
服务检查探针，用来探查服务是否健康。
如果检查失败，则通知MonitorAgent销毁.
如果检查成功，则通过StatusAgent更新状态.
*/
type SpiderAgent struct {
	Name        string
	NsqEndpoint string
	Namespace   []string
	Port        []int64
	Cmd         string
	Ip          string
	StopChan    chan int
	AlivaChan   chan int
}

// type SpiderMsg struct {
// 	Name    string
// 	Svcname string
// 	Msg     string
// }

// 只有当目标服务处于活动状态的时候才需要开始检测,
// 如果目标服务被Kill掉，则可以认为K8s健康检测失败
// 此时Spider就不需要再检查，开始上报检查结果
var needCheck bool
var msg = ""
var errNum = 0

func (this *SpiderAgent) Run() {
	ctx := getSpan()
	logrus.WithFields(logrus.Fields{SpiderAgentName: "Start"}).Info(SpiderAgentName)
	logrus.WithFields(logrus.Fields{"Send Start Signal": os.Getenv("DDOG_AGENT_SPIDER_SVC")}).Info(SpiderAgentName)
	logrus.WithFields(logrus.Fields{"span": ctx}).Info(ModuleName)
	data, _ := json.Marshal(monitor.MonitorModule{
		Kind:      this.Name,
		Svcname:   os.Getenv("DDOG_AGENT_SPIDER_SVC"),
		Namespace: os.Getenv("DDOG_AGENT_SPIDER_NS"),
		Msg:       "deploy",
		Span:      ctx,
		//Ip:        []string{this.Ip},
	})

	bridge.SendMonitorMsg(string(data))

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
				if len(porcess) == 2 && ((porcess[0].Executable() == "pause" && porcess[1].Executable() == "ddog-agent") || (porcess[1].Executable() == "pause" && porcess[0].Executable() == "ddog-agent")) {
					//	另外一个容器已经退出
					logrus.WithFields(logrus.Fields{"Send Quit Msg": msg}).Info(ModuleName)
					if msg != "" {
						data, _ := json.Marshal(monitor.MonitorModule{
							Kind:      this.Name,
							Svcname:   os.Getenv("DDOG_AGENT_SPIDER_SVC"),
							Namespace: os.Getenv("DDOG_AGENT_SPIDER_NS"),
							Msg:       "status",
							Status:    msg,
							Span:      ctx,
						})

						bridge.SendMonitorMsg(string(data))
					}

				}
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
			select {
			case p := <-this.AlivaChan:
				if p <= 1 {
					needCheck = false
				} else {
					needCheck = true
					this.checkPort(ctx)
				}

				/*errNum == 60*10几乎等同为10分钟*/
				if (!needCheck && msg != "") || (errNum == 60*10) {
					data, _ := json.Marshal(monitor.MonitorModule{
						Kind:      this.Name,
						Svcname:   os.Getenv("DDOG_AGENT_SPIDER_SVC"),
						Namespace: os.Getenv("DDOG_AGENT_SPIDER_NS"),
						Msg:       msg,
						Ip:        []string{this.Ip},
						Span:      ctx,
					})

					bridge.SendMonitorMsg(string(data))
					msg = ""
					errNum = 0
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
		// 不能退出,否则k8s会认为此服务处于异常状态
		// this.StopChan <- 0
	}
}

func (this *SpiderAgent) checkPort(ctx zmodel.SpanContext) {

	portMap := make(map[int64]int)

	for _, p := range this.Port {
		portMap[p] = 1
	}

	for _, td := range GOnetstat.Tcp() {
		logrus.WithFields(logrus.Fields{"tcp Port": td.Port, "state": td.State}).Info(SpiderAgentName)
		if strings.ToUpper(td.State) == "LISTEN" {
			if _, ok := portMap[td.Port]; ok {
				delete(portMap, td.Port)
			}
		} else {
			this.Ip = td.Ip
		}
	}

	//再过滤一遍TCP6
	for _, td := range GOnetstat.Tcp6() {
		logrus.WithFields(logrus.Fields{"tcp Port": td.Port, "state": td.State}).Info(SpiderAgentName)
		if strings.ToUpper(td.State) == "LISTEN" {
			if _, ok := portMap[td.Port]; ok {
				delete(portMap, td.Port)
			}
		}
	}

	if len(portMap) != 0 {
		msg = ""
		errNum ++
		for k, _ := range portMap {
			msg += fmt.Sprintf("Port[%d] Check Failed. ", k)
		}
		logrus.WithFields(logrus.Fields{"msg": msg}).Info(SpiderAgentName)
	} else {
		/* Task End */
		logrus.WithFields(logrus.Fields{"check end": true}).Info(SpiderAgentName)
		//ctx := zmodel.SpanContext{}
		//var traceid = ""
		//var id = ""
		//var parentid = ""
		//if os.Getenv("ZIPKIN_TRACID") != "" {
		//	traceid = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_TRACID"))
		//}
		//if os.Getenv("ZIPKIN_ID") != "" {
		//	id = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_ID"))
		//}
		//if os.Getenv("ZIPKIN_PARENTID") != "" {
		//	id = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_PARENTID"))
		//}
		//
		//_tracid := new(zmodel.TraceID)
		//_tracid.UnmarshalJSON([]byte(traceid))
		//ctx.TraceID = *_tracid
		//
		//_id := new(zmodel.ID)
		//_id.UnmarshalJSON([]byte(id))
		//ctx.ID = *_id
		//
		//_parentid := new(zmodel.ID)
		//_parentid.UnmarshalJSON([]byte(parentid))
		//ctx.ParentID = _parentid

		//ctx := getSpan()
		data, _ := json.Marshal(monitor.MonitorModule{
			Span:      ctx,
			Kind:      this.Name,
			Svcname:   os.Getenv("DDOG_AGENT_SPIDER_SVC"),
			Namespace: os.Getenv("DDOG_AGENT_SPIDER_NS"),
			Msg:       "ok",
			Ip:        []string{this.Ip},
		})

		bridge.SendMonitorMsg(string(data))
		<-make(chan int)
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

// func (this *SpiderAgent) isAlive() bool {
// 	for {
// 		now := time.Now()
// 		next := now.Add(time.Second * 1)
// 		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
// 		t := time.NewTimer(next.Sub(now))
// 		logrus.WithFields(logrus.Fields{"下次采集进程时间为": next.Format("200601021504")}).Info(SpiderAgentName)
//
// 		select {
// 		case <-t.C:
// 			porcess, err := ps.Processes()
// 			if err != nil {
// 				logrus.WithFields(logrus.Fields{"Get Process Error": err}).Error(SpiderAgentName)
// 			}
//
// 			logrus.WithFields(logrus.Fields{"Porcess": porcess}).Info(SpiderAgentName)
// 			if len(porcess) > 1 {
// 				return true
// 			}
// 		}
// 	}
// }

func getSpan() (ctx zmodel.SpanContext) {
	//ctx := zmodel.SpanContext{}
	var traceid = ""
	var id = ""
	var parentid = ""
	if os.Getenv("ZIPKIN_TRACID") != "" {
		traceid = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_TRACID"))
	}
	if os.Getenv("ZIPKIN_ID") != "" {
		id = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_ID"))
	}
	if os.Getenv("ZIPKIN_PARENTID") != "" {
		parentid = fmt.Sprintf("0%s0", os.Getenv("ZIPKIN_PARENTID"))
	}

	_tracid := new(zmodel.TraceID)
	_tracid.UnmarshalJSON([]byte(traceid))
	ctx.TraceID = *_tracid

	_id := new(zmodel.ID)
	_id.UnmarshalJSON([]byte(id))
	ctx.ID = *_id

	_parentid := new(zmodel.ID)
	_parentid.UnmarshalJSON([]byte(parentid))
	ctx.ParentID = _parentid

	return
}
