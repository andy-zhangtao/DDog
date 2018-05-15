package k8smodel

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15. 
type K8sDeployBackup struct {
	Apiversion string                   `yaml:"apiVersion"`
	Kind       string                   `yaml:"kind"`
	Metadata   K8sDeployBackup_Metadata `yaml:"metadata"`
	Spec       K8sDeployBackup_spec     `yaml:"spec"`
	Status     K8sDeploymentV1_status   `yaml:"status"`
}

type K8sDeployBackup_Metadata struct {
	Annotations       map[string]string `yaml:"annotations"`
	CreationTimestamp string            `yaml:"creationTimestamp"`
	Generation        int               `yaml:"generation"`
	Labels            map[string]string `yaml:"labels"`
	Name              string            `yaml:"name"`
	Namespace         string            `yaml:"namespace"`
	ResourceVersion   string            `yaml:"resourceVersion"`
	SelfLink          string            `yaml:"selfLink"`
	Uid               string            `yaml:"uid"`
}

type K8sDeployBackup_spec struct {
	MinReadySeconds      int                           `yaml:"minReadySeconds"`
	Replicas             int                           `yaml:"replicas"`
	RevisionHistoryLimit int                           `yaml:"revisionHistoryLimit"`
	Selector             K8sDeploymentV1_spec_selector `yaml:"selector"`
	Strategy             K8sDeploymentV1_spec_strategy `yaml:"strategy"`
	Template             K8sDeploymentV1_spec_template `yaml:"template"`
}

type K8sServiceBackup struct {
	Apiversion string                      `yaml:"apiVersion"`
	Kind       string                      `yaml:"kind"`
	Metadata   K8sServiceBackup_Metadata   `yaml:"metadata"`
	Spec       K8sServiceInfo_items_spec   `yaml:"spec"`
	Status     K8sServiceInfo_items_status `yaml:"status"`
}

type K8sServiceBackup_Metadata struct {
	Annotations       map[string]string `yaml:"annotations"`
	CreationTimestamp string            `yaml:"creationTimestamp"`
	Labels            map[string]string `yaml:"labels"`
	Name              string            `yaml:"name"`
	Namespace         string            `yaml:"namespace"`
	ResourceVersion   string            `yaml:"resourceVersion"`
	SelfLink          string            `yaml:"selfLink"`
	Uid               string            `yaml:"uid"`
}
