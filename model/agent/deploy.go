package agent

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

type DeployMsg struct {
	SvcName   string `json:"svc_name"`
	NameSpace string `json:"name_space"`
	Upgrade   bool   `json:"upgrade"`
	Replicas  int    `json:"replicas"`
}
