package k8s

type K8sMetaData struct {
	Endpoint  string `json:"endpoint"`
	Namespace string `json:"namespace"`
	Svcname   string `json:"svcname"`
	Version   string `json:"version"`
	Token     string `json:"token"`
}
