package k8smodel

// Generate By json.golang.chinazt.cc
// Please Don't Edit it !

type K8s struct {
	Kind       string       `json:"kind"`
	ApiVersion string       `json:"apiVersion"`
	Metadata   K8s_metadata `json:"metadata"`
	Items      []K8s_items  `json:"items"`
}

type K8s_metadata struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}
type K8s_items struct {
	Metadata K8s_items_metadata `json:"metadata"`
	Spec     K8s_items_spec     `json:"spec"`
	Status   K8s_items_status   `json:"status"`
}

type K8s_items_metadata_ownerReferences struct {
	ApiVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	Uid                string `json:"uid"`
	Controller         bool   `json:"controller"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
}
type K8s_items_metadata struct {
	Name              string                    `json:"name"`
	GenerateName      string                    `json:"generateName"`
	Namespace         string                    `json:"namespace"`
	SelfLink          string                    `json:"selfLink"`
	Uid               string                    `json:"uid"`
	ResourceVersion   string                    `json:"resourceVersion"`
	CreationTimestamp string                    `json:"creationTimestamp"`
	Labels            K8s_items_metadata_labels `json:"labels"`
	//Annotations       K8s_items_metadata_annotations       `json:"annotations"`
	OwnerReferences []K8s_items_metadata_ownerReferences `json:"ownerReferences"`
}
type K8s_items_status_conditions struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastTransitionTime string `json:"lasttransitiontime"`
}
type K8s_items_status_containerStatuses struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartcount"`
	Image        string `json:"image"`
	ImageID      string `json:"imageid"`
	ContainerID  string `json:"containerid"`
}

type K8s_items_status struct {
	Phase             string                               `json:"phase"`
	Conditions        []K8s_items_status_conditions        `json:"conditions"`
	HostIP            string                               `json:"hostIP"`
	PodIP             string                               `json:"podIP"`
	StartTime         string                               `json:"startTime"`
	ContainerStatuses []K8s_items_status_containerStatuses `json:"containerStatuses"`
	QosClass          string                               `json:"qosClass"`
}
type K8s_items_metadata_labels struct {
	Pod_template_hash         string `json:"pod-template-hash"`
	Qcloud_app                string `json:"qcloud-app"`
	Qcloud_redeploy_timestamp string `json:"qcloud-redeploy-timestamp"`
}

type K8s_items_metadata_annotations struct {
	Kubernetes_Io_Created_by string `json:"kubernetes.io/created_by"`
}

type K8s_items_spec_volumes struct {
	Name string `json:"name"`
}
type K8s_items_spec_containers struct {
	Name                     string `json:"name"`
	Image                    string `json:"image"`
	WorkingDir               string `json:"workingDir"`
	TerminationMessagePath   string `json:"terminationmessagepath"`
	TerminationMessagePolicy string `json:"terminationmessagepolicy"`
	ImagePullPolicy          string `json:"imagepullpolicy"`
}

type K8s_items_spec_securityContext struct {
}

type K8s_items_spec_imagePullSecrets struct {
	Name string `json:"name"`
}

type K8s_items_spec struct {
	Volumes                       []K8s_items_spec_volumes          `json:"volumes"`
	Containers                    []K8s_items_spec_containers       `json:"containers"`
	RestartPolicy                 string                            `json:"restartpolicy"`
	TerminationGracePeriodSeconds int                               `json:"terminationgraceperiodseconds"`
	DnsPolicy                     string                            `json:"dnspolicy"`
	ServiceAccountName            string                            `json:"serviceaccountname"`
	ServiceAccount                string                            `json:"serviceaccount"`
	NodeName                      string                            `json:"nodename"`
	SecurityContext               K8s_items_spec_securityContext    `json:"securitycontext"`
	ImagePullSecrets              []K8s_items_spec_imagePullSecrets `json:"imagepullsecrets"`
	SchedulerName                 string                            `json:"schedulername"`
}
