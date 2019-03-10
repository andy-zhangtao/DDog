package _const

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

import zmodel "github.com/openzipkin/zipkin-go/model"

const (
	SvcDestroyMsg    = "DestroySvc"
	SvcMonitorMsg    = "MonitorSvc"
	SvcDeployMsg     = "DeploySvc"
	SvcReplicaMsg    = "ReplicaSvc"
	SvcK8sMonitorMsg = "K8sMonitorSvc"
)

type DestoryMsg struct {
	Svcname   string             `json:"svcname"`
	Namespace string             `json:"namespace"`
	Span      zmodel.SpanContext `json:"span"`
}
