package _const

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

const (
	SvcDestroyMsg    = "DestroySvc"
	SvcMonitorMsg    = "MonitorSvc"
	SvcDeployMsg     = "DeploySvc"
	SvcReplicaMsg    = "ReplicaSvc"
	SvcK8sMonitorMsg = "K8sMonitorSvc"
)

type DestoryMsg struct {
	Svcname   string `json:"svcname"`
	Namespace string `json:"namespace"`
}
