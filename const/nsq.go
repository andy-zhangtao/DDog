package _const

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

import zmodel "github.com/openzipkin/zipkin-go/model"

const (
	SvcDestroyMsg    = "DestroySvc"
	SvcMonitorMsg    = "MonitorSvc"
	SvcDeployMsg     = "DeploySvc"
	SvcReplicaMsg    = "ReplicaSvc"
	SvcK8sMonitorMsg = "K8sMonitorSvc"
	SvcK8sAPIDelete  = "k8s-%s-delete"
)

type DestoryMsg struct {
	Svcname   string             `json:"svcname"`
	Namespace string             `json:"namespace"`
	Span      zmodel.SpanContext `json:"span"`
}

// ActionMsgForK8sAPI 调用K8sAPI
type ActionMsgForK8sAPI struct {
	Action         string `json:"action"`
	NameSpace      string `json:"nameSpace"`
	DeploymentName string `json:"deploymentName"`
}
