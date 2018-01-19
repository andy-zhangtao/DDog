package k8smodel

import "encoding/json"

// K8s Deployment Struct & Function

// Generate By Json.Golang.Chinazt.CC
// Please Don't Edit Manual!

type K8sDeployment struct {
	Kind       string                 `json:"kind"`
	ApiVersion string                 `json:"apiversion"`
	Metadata   K8sDeployment_metadata `json:"metadata"`
	Spec       K8sDeployment_spec     `json:"spec"`
	Status     K8sDeployment_status   `json:"status"`
}

type K8sDeployment_metadata struct {
	Name              string                             `json:"name"`
	Namespace         string                             `json:"namespace"`
	SelfLink          string                             `json:"selflink"`
	Uid               string                             `json:"uid"`
	ResourceVersion   string                             `json:"resourceversion"`
	Generation        int                                `json:"generation"`
	CreationTimestamp string                             `json:"creationtimestamp"`
	Labels            K8sDeployment_metadata_labels      `json:"labels"`
	Annotations       K8sDeployment_metadata_annotations `json:"annotations"`
}

type K8sDeployment_metadata_labels struct {
	Qcloud_app string `json:"qcloud_app"`
}
type K8sDeployment_spec_selector_matchLabels struct {
	Qcloud_app string `json:"qcloud_app"`
}
type K8sDeployment_spec_selector struct {
	MatchLabels K8sDeployment_spec_selector_matchLabels `json:"matchlabels"`
}
type K8sDeployment_spec_strategy struct {
	Type          string                                    `json:"type"`
	RollingUpdate K8sDeployment_spec_strategy_rollingUpdate `json:"rollingupdate"`
}
type K8sDeployment_status_conditions struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastUpdateTime     string `json:"lastupdatetime"`
	LastTransitionTime string `json:"lasttransitiontime"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
}
type K8sDeployment_status struct {
	ObservedGeneration int                               `json:"observedgeneration"`
	Replicas           int                               `json:"replicas"`
	UpdatedReplicas    int                               `json:"updatedreplicas"`
	ReadyReplicas      int                               `json:"readyreplicas"`
	AvailableReplicas  int                               `json:"availablereplicas"`
	Conditions         []K8sDeployment_status_conditions `json:"conditions"`
}

type K8sDeployment_spec struct {
	Replicas             int                         `json:"replicas"`
	Selector             K8sDeployment_spec_selector `json:"selector"`
	Template             K8sDeployment_spec_template `json:"template"`
	Strategy             K8sDeployment_spec_strategy `json:"strategy"`
	RevisionHistoryLimit int                         `json:"revisionhistorylimit"`
}

type K8sDeployment_metadata_annotations struct {
	DeploymentChangecourse         string `json:"deployment.changecourse"`
	DeploymentKubernetesIoRevision string `json:"deployment.kubernetes.io/revision"`
}

type K8sDeployment_spec_template struct {
	Metadata K8sDeployment_metadata `json:"metadata"`
	//Spec     K8sDeployment_spec     `json:"spec"`
}
type K8sDeployment_spec_strategy_rollingUpdate struct {
	MaxUnavailable int `json:"maxunavailable"`
	MaxSurge       int `json:"maxsurge"`
}


// K8dUnmarshal 将字符串解析成K8sDeployment对象指针
func K8dUnmarshal(data []byte) (*K8sDeployment, error) {
	var k8d K8sDeployment

	err := json.Unmarshal(data, &k8d)
	if err != nil {
		return nil, err
	}

	return &k8d, nil
}
