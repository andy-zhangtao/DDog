package k8smodel

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
type K8sDeploymentInfo struct {
	Kind       string                 `json:"kind"`
	ApiVersion string                 `json:"apiversion"`
	Metadata   K8sDeployment_metadata `json:"metadata"`
	Items      []K8sDeploymentV1      `json:"items"`
}
type K8sDeploymentV1 struct {
	Metadata K8sDeploymentV1_metadata `json:"metadata"`
	Spec     K8sDeploymentV1_spec     `json:"spec"`
	Status   K8sDeploymentV1_status   `json:"status"`
}
type K8sDeploymentV1_metadata struct {
	Name              string            `json:"name" yaml:"-"`
	Namespace         string            `json:"namespace" yaml:"-"`
	SelfLink          string            `json:"selfLink" yaml:"-"`
	Uid               string            `json:"uid" yaml:"-"`
	ResourceVersion   string            `json:"resourceVersion" yaml:"-"`
	Generation        int               `json:"generation" yaml:"-"`
	CreationTimestamp string            `json:"creationTimestamp" yaml:"creationTimestamp"`
	Labels            map[string]string `json:"labels" yaml:"labels"`
	Annotations       map[string]string `json:"annotations" yaml:"-"`
}
type K8sDeploymentV1_spec_selector_matchLabels struct {
	Qcloud_app string `json:"qcloud_app"`
}
type K8sDeploymentV1_spec_strategy_rollingUpdate struct {
	MaxUnavailable int         `json:"maxUnavailable" yaml:"maxUnavailable"`
	MaxSurge       interface{} `json:"maxSurge" yaml:"maxSurge"`
}
type K8sDeploymentV1_spec struct {
	Replicas             int                           `json:"replicas"`
	Selector             K8sDeploymentV1_spec_selector `json:"selector"`
	Template             K8sDeploymentV1_spec_template `json:"template"`
	Strategy             K8sDeploymentV1_spec_strategy `json:"strategy"`
	MinReadySeconds      int                           `json:"minreadyseconds"`
	RevisionHistoryLimit int                           `json:"revisionhistorylimit"`
}
type K8sDeploymentV1_status struct {
	ObservedGeneration int                                 `json:"observedGeneration" yaml:"observedGeneration"`
	Replicas           int                                 `json:"replicas" yaml:"replicas"`
	UpdatedReplicas    int                                 `json:"updatedReplicas" yaml:"updatedReplicas"`
	ReadyReplicas      int                                 `json:"readyreplicas" yaml:"-"`
	AvailableReplicas  int                                 `json:"availableReplicas" yaml:"availableReplicas"`
	Conditions         []K8sDeploymentV1_status_conditions `json:"conditions" yaml:"-"`
}
type K8sDeploymentV1_status_conditions struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastUpdateTime     string `json:"lastupdatetime"`
	LastTransitionTime string `json:"lasttransitiontime"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
}
type K8sDeploymentV1_metadata_labels struct {
	Qcloud_app string `json:"qcloud_app"`
}

type K8sDeploymentV1_spec_selector struct {
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}
type K8sDeploymentV1_spec_template struct {
	Metadata K8sDeploymentV1_metadata `json:"metadata" yaml:"metadata"`
	Spec     K8sDeploymentV1_Spec     `json:"spec" yaml:"spec"`
}
type K8sDeploymentV1_spec_strategy struct {
	Type          string                                      `json:"type"`
	RollingUpdate K8sDeploymentV1_spec_strategy_rollingUpdate `json:"rollingUpdate" yaml:"rollingUpdate"`
}

type K8sDeploymentV1_Spec_containers struct {
	Name                     string              `json:"name" yaml:"name"`
	Image                    string              `json:"image" yaml:"image"`
	TerminationMessagePath   string              `json:"terminationMessagePath" yaml:"terminationMessagePath"`
	TerminationMessagePolicy string              `json:"terminationmessagepolicy" yaml:"-"`
	ImagePullPolicy          string              `json:"imagePullPolicy" yaml:"imagePullPolicy"`
	Env                      []map[string]string `json:"env" yaml:"env"`
	Resources struct {
		Limits   map[string]string `json:"limits" yaml:"limits"`
		Requests map[string]string `json:"requests" yaml:"requests"`
	} `json:"resources" yaml:"resources"`
	SecurityContext map[string]bool `json:"securityContext" yaml:"securityContext"`
}
type K8sDeploymentV1_Spec_securityContext struct {
}
type K8sDeploymentV1_Spec_imagePullSecrets struct {
	Name string `json:"name"`
}
type K8sDeploymentV1_Spec struct {
	Containers                    []K8sDeploymentV1_Spec_containers    `json:"containers" yaml:"containers"`
	RestartPolicy                 string                               `json:"restartPolicy" yaml:"restartPolicy"`
	TerminationGracePeriodSeconds int                                  `json:"terminationGracePeriodSeconds" yaml:"terminationGracePeriodSeconds"`
	DnsPolicy                     string                               `json:"dnsPolicy" yaml:"dnsPolicy"`
	SecurityContext               K8sDeploymentV1_Spec_securityContext `json:"securityContext" yaml:"securityContext"`
	ImagePullSecrets              []map[string]string                  `json:"imagePullSecrets" yaml:"imagePullSecrets"`
	SchedulerName                 string                               `json:"schedulerName" yaml:"-"`
	ServiceAccountName            string                               `json:"serviceAccountName" yaml:"serviceAccountName"`
}
