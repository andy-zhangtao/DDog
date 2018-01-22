package service

import "encoding/json"

type Instance struct {
	Name         string              `json:"name"`
	Status       string              `json:"status"`
	Reason       string              `json:"reason"`
	NodeIp       string              `json:"nodeIp"`
	NodeName     string              `json:"nodeName"`
	Ip           string              `json:"ip"`
	RestartCount int                 `json:"restartCount"`
	ReadyCount   int                 `json:"readyCount"`
	CreatedAt    string              `json:"createdAt"`
	Containers   []InstanceContainer `json:"containers"`
}

type InstanceContainer struct {
	Name        string `json:"name"`
	ContainerId string `json:"container_id"`
	Status      string `json:"status"`
	Reason      string `json:"reason"`
	Image       string `json:"image"`
}

func InstanceUnmarshal(data []byte) (*Instance, error) {
	var inst Instance

	err := json.Unmarshal(data, &inst)
	if err != nil {
		return nil, err
	}

	return &inst, nil
}
