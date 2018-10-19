package k8smodel

import "time"

type K8sService struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		Name              string    `json:"name"`
		Namespace         string    `json:"namespace"`
		SelfLink          string    `json:"selfLink"`
		UID               string    `json:"uid"`
		ResourceVersion   string    `json:"resourceVersion"`
		CreationTimestamp time.Time `json:"creationTimestamp"`
		Labels            struct {
			QcloudApp string `json:"qcloud-app"`
		} `json:"labels"`
		Annotations struct {
			ServiceKubernetesIoQcloudLoadbalancerClusterid        string `json:"service.kubernetes.io/qcloud-loadbalancer-clusterid"`
			ServiceKubernetesIoQcloudLoadbalancerInternal         string `json:"service.kubernetes.io/qcloud-loadbalancer-internal"`
			ServiceKubernetesIoQcloudLoadbalancerInternalSubnetid string `json:"service.kubernetes.io/qcloud-loadbalancer-internal-subnetid"`
		} `json:"annotations"`
	} `json:"metadata"`
	Spec struct {
		Ports []struct {
			Name       string `json:"name"`
			Protocol   string `json:"protocol"`
			Port       int    `json:"port"`
			TargetPort int    `json:"targetPort"`
			NodePort   int    `json:"nodePort"`
		} `json:"ports"`
		Selector struct {
			QcloudApp string `json:"qcloud-app"`
		} `json:"selector"`
		ClusterIP             string `json:"clusterIP"`
		Type                  string `json:"type"`
		SessionAffinity       string `json:"sessionAffinity"`
		ExternalTrafficPolicy string `json:"externalTrafficPolicy"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
			Ingress []struct {
				IP string `json:"ip"`
			} `json:"ingress"`
		} `json:"loadBalancer"`
	} `json:"status"`
}
