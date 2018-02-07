package monitor

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/7.
type MonitorModule struct {
	Kind      string `json:"kind"`
	Svcname   string `json:"svcname"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
}
