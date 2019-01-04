package k8smodel

import "time"

type K8sPodsInfo struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink        string `json:"selfLink"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []struct {
		Metadata struct {
			Name              string    `json:"name"`
			GenerateName      string    `json:"generateName"`
			Namespace         string    `json:"namespace"`
			SelfLink          string    `json:"selfLink"`
			UID               string    `json:"uid"`
			ResourceVersion   string    `json:"resourceVersion"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Labels            struct {
				PodTemplateHash string `json:"pod-template-hash"`
				QcloudApp       string `json:"qcloud-app"`
			} `json:"labels"`
			Annotations struct {
				KubernetesIoCreatedBy string `json:"kubernetes.io/created-by"`
			} `json:"annotations"`
			OwnerReferences []struct {
				APIVersion         string `json:"apiVersion"`
				Kind               string `json:"kind"`
				Name               string `json:"name"`
				UID                string `json:"uid"`
				Controller         bool   `json:"controller"`
				BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
			} `json:"ownerReferences"`
		} `json:"metadata"`
		Spec struct {
			Volumes []struct {
				Name   string `json:"name"`
				Secret struct {
					SecretName  string `json:"secretName"`
					DefaultMode int    `json:"defaultMode"`
				} `json:"secret"`
			} `json:"volumes"`
			Containers []struct {
				Name  string `json:"name"`
				Image string `json:"image"`
				Env   []struct {
					Name  string `json:"name"`
					Value string `json:"value,omitempty"`
				} `json:"env"`
				Resources struct {
					Limits struct {
						Memory string `json:"memory"`
					} `json:"limits"`
					Requests struct {
						Memory string `json:"memory"`
					} `json:"requests"`
				} `json:"resources"`
				VolumeMounts []struct {
					Name      string `json:"name"`
					ReadOnly  bool   `json:"readOnly"`
					MountPath string `json:"mountPath"`
				} `json:"volumeMounts"`
				LivenessProbe struct {
					TCPSocket struct {
						Port int `json:"port"`
					} `json:"tcpSocket"`
					InitialDelaySeconds int `json:"initialDelaySeconds"`
					TimeoutSeconds      int `json:"timeoutSeconds"`
					PeriodSeconds       int `json:"periodSeconds"`
					SuccessThreshold    int `json:"successThreshold"`
					FailureThreshold    int `json:"failureThreshold"`
				} `json:"livenessProbe,omitempty"`
				ReadinessProbe struct {
					TCPSocket struct {
						Port int `json:"port"`
					} `json:"tcpSocket"`
					InitialDelaySeconds int `json:"initialDelaySeconds"`
					TimeoutSeconds      int `json:"timeoutSeconds"`
					PeriodSeconds       int `json:"periodSeconds"`
					SuccessThreshold    int `json:"successThreshold"`
					FailureThreshold    int `json:"failureThreshold"`
				} `json:"readinessProbe,omitempty"`
				TerminationMessagePath   string `json:"terminationMessagePath"`
				TerminationMessagePolicy string `json:"terminationMessagePolicy"`
				ImagePullPolicy          string `json:"imagePullPolicy"`
				SecurityContext          struct {
					Privileged bool `json:"privileged"`
				} `json:"securityContext"`
			} `json:"containers"`
			RestartPolicy                 string `json:"restartPolicy"`
			TerminationGracePeriodSeconds int    `json:"terminationGracePeriodSeconds"`
			DNSPolicy                     string `json:"dnsPolicy"`
			ServiceAccountName            string `json:"serviceAccountName"`
			ServiceAccount                string `json:"serviceAccount"`
			NodeName                      string `json:"nodeName"`
			SecurityContext               struct {
			} `json:"securityContext"`
			ImagePullSecrets []struct {
				Name string `json:"name"`
			} `json:"imagePullSecrets"`
			SchedulerName string `json:"schedulerName"`
		} `json:"spec"`
		Status struct {
			Phase      string `json:"phase"`
			Conditions []struct {
				Type               string      `json:"type"`
				Status             string      `json:"status"`
				LastProbeTime      interface{} `json:"lastProbeTime"`
				LastTransitionTime time.Time   `json:"lastTransitionTime"`
			} `json:"conditions"`
			HostIP            string    `json:"hostIP"`
			PodIP             string    `json:"podIP"`
			StartTime         time.Time `json:"startTime"`
			ContainerStatuses []struct {
				Name  string `json:"name"`
				State struct {
					Running struct {
						StartedAt time.Time `json:"startedAt"`
					} `json:"running"`
				} `json:"state"`
				LastState struct {
				} `json:"lastState"`
				Ready        bool   `json:"ready"`
				RestartCount int    `json:"restartCount"`
				Image        string `json:"image"`
				ImageID      string `json:"imageID"`
				ContainerID  string `json:"containerID"`
			} `json:"containerStatuses"`
			QosClass string `json:"qosClass"`
		} `json:"status"`
	} `json:"items"`
}
