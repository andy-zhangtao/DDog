package agents

import (
	"strings"
	"os"
	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/model/monitor"
	"fmt"
	"github.com/andy-zhangtao/DDog/bridge"
	"time"
)

/**
RetriAgent
在指定的命名空间中检索当前失败次数过多的服务实例
当发现有失败次数太多的服务实例之后，RetriDog将会
尝试获取失败原因，并且通知MonitorAgent.然后关闭
异常服务实例，同时将服务实例状态调整为失败
 */
type RetriAgent struct {
	Name        string
	NsqEndpoint string
	Namespace   []string
	StopChan    chan int
}

type RetriMsg struct {
	Name    string
	Svcname string
	Msg     string
	Count   int
}

func (this *RetriAgent) Run() {
	logrus.WithFields(logrus.Fields{RetriAgentName: "Start"}).Info(RetriAgentName)

	if os.Getenv("DDOG_AGENT_RETRI_NAMESPACE") != "" && strings.Contains(os.Getenv("DDOG_AGENT_RETRI_NAMESPACE"), ";") {
		this.Namespace = strings.Split(os.Getenv("DDOG_AGENT_RETRI_NAMESPACE"), ";")
	}
	if len(this.Namespace) == 0 {
		logrus.WithFields(logrus.Fields{RetriAgentName: "Namespace Emtpy!", "DDOG_AGENT_RETRI_NAMESPACE": os.Getenv("DDOG_AGENT_RETRI_NAMESPACE")}).Error(RetriAgentName)
		this.StopChan <- 1
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		logrus.WithFields(logrus.Fields{RetriAgentName: "Get MetaData Error!", "error": err}).Error(RetriAgentName)
		this.StopChan <- 1
	}

	for {
		now := time.Now()
		next := now.Add(time.Minute * 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		logrus.WithFields(logrus.Fields{"下次采集时间为": next.Format("200601021504")}).Info(RetriAgentName)

		select {
		case <-t.C:
			err = this.getSVC(md)
			if err != nil {
				logrus.WithFields(logrus.Fields{RetriAgentName: "Get SvcStatus Error!", "error": err}).Error(RetriAgentName)
				this.StopChan <- 1
				return
			}
		}
	}
}

func (this *RetriAgent) getSVC(md *metadata.MetaData) error {
	for _, ns := range this.Namespace {
		if ns != "" {
			logrus.WithFields(logrus.Fields{"RetriNamespace": ns}).Info(RetriAgentName)
			s := service.Svc{
				Pub: public.Public{
					SecretId: md.Sid,
					Region:   md.Region,
				},
				ClusterId: md.ClusterID,
				Namespace: ns,
				SecretKey: md.Skey,
			}

			s.SetDebug(true)
			sm, err := s.QuerySampleInfo()
			if err != nil {
				return err
			}

			q := service.Service{
				Pub: public.Public{
					SecretId: md.Sid,
					Region:   md.Region,
				},
				ClusterId: md.ClusterID,
				Namespace: ns,
				SecretKey: md.Skey,
			}

			for _, sv := range sm.Data.Services {
				logrus.WithFields(logrus.Fields{"name": sv}).Info(RetriAgentName)
				if strings.ToLower(sv.Status) != strings.ToLower("Normal") {
					q.ServiceName = sv.ServiceName
					event, err := q.DescribeServiceEvent()
					if err != nil {
						return err
					}

					rm := RetriMsg{Name: RetriAgentName, Svcname: sv.ServiceName}
					for _, e := range event.Data.EventList {
						if strings.ToLower(e.Level) == strings.ToLower("Warning") && e.Count > 10 {
							rm.Count = e.Count
							rm.Msg = e.Reason
						}
					}

					if rm.Count > 0 {
						data, _ := json.Marshal(monitor.MonitorModule{
							Kind:      rm.Name,
							Svcname:   rm.Svcname,
							Namespace: ns,
							Msg:       fmt.Sprintf("count:[%d]msg:[%s]", rm.Count, rm.Msg),
						})

						err = bridge.SendMonitorMsg(string(data))
						if err != nil {
							return err
						}
					}
				}
			}
		}

	}

	return nil
}
