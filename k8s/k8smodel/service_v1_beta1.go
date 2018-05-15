package k8smodel

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.

type K8sServiceInfo struct {
	Kind       string                  `json:"kind"`
	ApiVersion string                  `json:"apiVersion"`
	Metadata   K8sServiceInfo_metadata `json:"metadata"`
	Items      []K8sServiceInfo_items  `json:"items"`
}

type K8sServiceInfo_metadata struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

type K8sServiceInfo_items struct {
	Metadata K8sServiceInfo_items_metadata `json:"metadata"`
	Spec     K8sServiceInfo_items_spec     `json:"spec"`
	Status   K8sServiceInfo_items_status   `json:"status"`
}

type K8sServiceInfo_items_metadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	SelfLink          string            `json:"selfLink"`
	Uid               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

type K8sServiceInfo_items_status struct {
	LoadBalancer map[string][]map[string]string `json:"loadBalancer" yaml:"loadBalancer"`
}

type K8sServiceInfo_items_spec_ports struct {
	Name       string `json:"name" yaml:"name"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Port       int    `json:"port" yaml:"port"`
	TargetPort int    `json:"targetPort" yaml:"targetPort"`
	NodePort   int    `json:"nodePort" yaml:"nodePort"`
}

type K8sServiceInfo_items_spec struct {
	Ports                 []K8sServiceInfo_items_spec_ports `json:"ports" yaml:"ports"`
	Selector              map[string]string                 `json:"selector" yaml:"selector"`
	ClusterIP             string                            `json:"clusterIP" yaml:"clusterIP"`
	Type                  string                            `json:"type" yaml:"type"`
	SessionAffinity       string                            `json:"sessionAffinity" yaml:"sessionAffinity"`
	ExternalTrafficPolicy string                            `json:"externalTrafficPolicy" yaml:"-"`
}
