package service

import "encoding/json"

type ServiceInfo struct {
	ServiceName     string                     `json:"servicename"`
	ServiceDesc     string                     `json:"servicedesc"`
	Status          string                     `json:"status"`
	ReasonMap       map[string]int             `json:"reasonmap"`
	Reason          string                     `json:"reason"`
	RegionId        int                        `json:"regionid"`
	DesiredReplicas int                        `json:"desiredreplicas"`
	CurrentReplicas int                        `json:"currentreplicas"`
	LbId            string                     `json:"lbid"`
	LbStatus        string                     `json:"lbstatus"`
	CreatedAt       string                     `json:"createdat"`
	AccessType      string                     `json:"accesstype"`
	ServiceIp       string                     `json:"serviceip"`
	ExternalIp      string                     `json:"externalip"`
	Namespace       string                     `json:"namespace"`
	PortMappings    []ServiceInfo_portMappings `json:"portmappings"`
	Containers      []ServiceInfo_containers   `json:"containers"`
	Selector        ServiceInfo_selector       `json:"selector"`
	Labels          ServiceInfo_labels         `json:"labels"`
}

type ServiceInfo_portMappings struct {
	ContainerPort int    `json:"containerport"`
	LbPort        int    `json:"lbport"`
	NodePort      int    `json:"nodeport"`
	Protocol      string `json:"protocol"`
}

type ServiceInfo_containers struct {
	ContainerName string `json:"containername"`
	Image         string `json:"image"`
	Cpu           int    `json:"cpu"`
	Memory        int    `json:"memory"`
	Command       string `json:"command"`
}

type ServiceInfo_selector struct {
	Qcloud_app string `json:"qcloud_app"`
}

type ServiceInfo_labels struct {
	Qcloud_app string `json:"qcloud_app"`
}

// SvcInfoUnmarshal 解析服务详细数据
func SvcInfoUnmarshal(data []byte) (*ServiceInfo, error) {
	var svcInfo ServiceInfo

	err := json.Unmarshal(data, &svcInfo)
	if err != nil {
		return nil, err
	}

	return &svcInfo, nil
}
