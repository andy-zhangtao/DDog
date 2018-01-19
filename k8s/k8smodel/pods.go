package k8smodel

import "encoding/json"

type K8sPods struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiversion"`
	Metadata   K8sPods_metadata `json:"metadata"`
	Items      []K8sPods_items  `json:"items"`
}

type K8sPods_metadata struct {
	SelfLink        string `json:"selflink"`
	ResourceVersion string `json:"resourceversion"`
}
type K8sPods_items_spec_volumes struct {
	Name string `json:"name"`
}
type K8sPods_items struct {
	Metadata K8sPods_items_metadata `json:"metadata"`
	Spec     K8sPods_items_spec     `json:"spec"`
	Status   K8sPods_items_status   `json:"status"`
}

type K8sPods_items_status_conditions struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastTransitionTime string `json:"lasttransitiontime"`
}
type K8sPods_items_status_containerStatuses struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartcount"`
	Image        string `json:"image"`
	ImageID      string `json:"imageid"`
	ContainerID  string `json:"containerid"`
}
type K8sPods_items_spec_securityContext struct {
}
type K8sPods_items_spec_imagePullSecrets struct {
	Name string `json:"name"`
}
type K8sPods_items_spec struct {
	Volumes                       []K8sPods_items_spec_volumes          `json:"volumes"`
	Containers                    []K8sPods_items_spec_containers       `json:"containers"`
	RestartPolicy                 string                                `json:"restartpolicy"`
	TerminationGracePeriodSeconds int                                   `json:"terminationgraceperiodseconds"`
	DnsPolicy                     string                                `json:"dnspolicy"`
	ServiceAccountName            string                                `json:"serviceaccountname"`
	ServiceAccount                string                                `json:"serviceaccount"`
	NodeName                      string                                `json:"nodename"`
	SecurityContext               K8sPods_items_spec_securityContext    `json:"securitycontext"`
	ImagePullSecrets              []K8sPods_items_spec_imagePullSecrets `json:"imagepullsecrets"`
	SchedulerName                 string                                `json:"schedulername"`
}
type K8sPods_items_metadata_labels struct {
	Pod_template_hash         string `json:"pod_template_hash"`
	Qcloud_app                string `json:"qcloud-app"`
	Qcloud_redeploy_timestamp string `json:"qcloud-redeploy-timestamp"`
}
type K8sPods_items_metadata_annotations struct {
	KubernetesIoCreatedBy string `json:"kubernetes.io/created_by"`
}
type K8sPods_items_metadata_ownerReferences struct {
	ApiVersion         string `json:"apiversion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	Uid                string `json:"uid"`
	Controller         bool   `json:"controller"`
	BlockOwnerDeletion bool   `json:"blockownerdeletion"`
}
type K8sPods_items_metadata struct {
	Name              string                                   `json:"name"`
	GenerateName      string                                   `json:"generatename"`
	Namespace         string                                   `json:"namespace"`
	SelfLink          string                                   `json:"selflink"`
	Uid               string                                   `json:"uid"`
	ResourceVersion   string                                   `json:"resourceversion"`
	CreationTimestamp string                                   `json:"creationtimestamp"`
	Labels            K8sPods_items_metadata_labels            `json:"labels"`
	Annotations       K8sPods_items_metadata_annotations       `json:"annotations"`
	OwnerReferences   []K8sPods_items_metadata_ownerReferences `json:"ownerreferences"`
}
type K8sPods_items_spec_containers struct {
	Name                     string `json:"name"`
	Image                    string `json:"image"`
	WorkingDir               string `json:"workingdir"`
	TerminationMessagePath   string `json:"terminationmessagepath"`
	TerminationMessagePolicy string `json:"terminationmessagepolicy"`
	ImagePullPolicy          string `json:"imagepullpolicy"`
}
type K8sPods_items_status struct {
	Phase             string                                   `json:"phase"`
	Conditions        []K8sPods_items_status_conditions        `json:"conditions"`
	HostIP            string                                   `json:"hostip"`
	PodIP             string                                   `json:"podip"`
	StartTime         string                                   `json:"starttime"`
	ContainerStatuses []K8sPods_items_status_containerStatuses `json:"containerstatuses"`
	QosClass          string                                   `json:"qosclass"`
}

// K8pUnmarshal 将字符串解析成K8sDeployment对象指针
func K8pUnmarshal(data []byte) (*K8sPods, error) {
	var k8p K8sPods

	err := json.Unmarshal(data, &k8p)
	if err != nil {
		return nil, err
	}

	return &k8p, nil
}
