package k8smodel

import (
	"testing"

	"gopkg.in/yaml.v2"
	"github.com/stretchr/testify/assert"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
func TestK8sDeployBackup(t *testing.T) {

	depBackup := K8sDeployBackup{
		Apiversion: "extensions/v1beta1",
		Kind:       "Deployment",
		Metadata: K8sDeployBackup_Metadata{
			Annotations: map[string]string{
				"deployment.changecourse":           "Updating",
				"deployment.kubernetes.io/revision": "13",
				"description":                       "Collect Container Resource",
			},
			CreationTimestamp: "2018-05-05T10:00:20Z",
			Generation:        13,
			Labels: map[string]string{
				"qcloud-app": "devex-resource-agent",
			},
			Name:            "devex-resource-agent",
			Namespace:       "scheduler",
			ResourceVersion: "10913887877",
			SelfLink:        "/apis/extensions/v1beta1/namespaces/scheduler/deployments/devex-resource-agent",
			Uid:             "1ebef947-504b-11e8-b75e-52540018543c",
		},
		Spec: K8sDeployBackup_spec{
			MinReadySeconds:      10,
			Replicas:             1,
			RevisionHistoryLimit: 5,
			Selector: K8sDeploymentV1_spec_selector{MatchLabels: map[string]string{
				"qcloud-app": "devex-resource-agent",
			}},
			Strategy: K8sDeploymentV1_spec_strategy{
				Type:          "RollingUpdate",
				RollingUpdate: K8sDeploymentV1_spec_strategy_rollingUpdate{MaxUnavailable: 0, MaxSurge: 1},
			},
			Template: K8sDeploymentV1_spec_template{
				Metadata: K8sDeploymentV1_metadata{CreationTimestamp: "", Labels: map[string]string{
					"qcloud-app":                "devex-resource-agent",
					"qcloud-redeploy-timestamp": "1525771236",
				},
				}, Spec: K8sDeploymentV1_Spec{
					Containers: []K8sDeploymentV1_Spec_containers{
						{
							Name:                     "devex-resource-agent",
							Image:                    "ccr.ccs.tencentyun.com/eqxiu/devex-resource-agent:latest",
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "",
							ImagePullPolicy:          "Always",
							Env: []map[string]string{

								{
									"name":  "PATH",
									"value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
								},
								{
									"name":  "Agent_Nsq_Endpoint",
									"value": "192.168.1.12:4150",
								},
								{
									"name":  "log_opt",
									"value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\"",
								},
								{
									"name":  "svcname",
									"value": "devex-resource-agent",
								},
								{
									"name":  "LOGCHAIN_DRIVER",
									"value": "influx",
								},
								{
									"name":  "Agent_Influx_Endpoint",
									"value": "192.168.1.14:8089",
								},
								{
									"name":  "Agent_Mongo_endpoint",
									"value": "mdb1.yqxiu.cn:27010,mdb.yqxiu.cn:27010",
								},
								{
									"name":  "Agent_Mongo_DB",
									"value": "data-mgr",
								},
								{
									"name":  "Agent_Influx_TCP_Endpoint",
									"value": "http://192.168.1.14:8086",
								},
							},
							Resources: struct {
								Limits   map[string]string `json:"limits" yaml:"limits"`;
								Requests map[string]string `json:"requests" yaml:"requests"`
							}{
								Limits: map[string]string{
									"cpu":    "500m",
									"memory": "50Mi",
								},
								Requests: map[string]string{
									"cpu":    "250m",
									"memory": "20Mi",
								},
							},
							SecurityContext: map[string]bool{
								"privileged": false,
							}},
					},
					RestartPolicy:                 "Always",
					TerminationGracePeriodSeconds: 30,
					DnsPolicy:                     "ClusterFirst",
					SecurityContext:               K8sDeploymentV1_Spec_securityContext{},
					ImagePullSecrets: []map[string]string{
						{"name": "qcloudregistrykey"},
					},
					SchedulerName:      "",
					ServiceAccountName: ""}}},
		Status: K8sDeploymentV1_status{ObservedGeneration: 13, Replicas: 1, UpdatedReplicas: 1, AvailableReplicas: 1,},
	}

	var expect = `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.changecourse: Updating
    deployment.kubernetes.io/revision: "13"
    description: Collect Container Resource
  creationTimestamp: 2018-05-05T10:00:20Z
  generation: 13
  labels:
    qcloud-app: devex-resource-agent
  name: devex-resource-agent
  namespace: scheduler
  resourceVersion: "10913887877"
  selfLink: /apis/extensions/v1beta1/namespaces/scheduler/deployments/devex-resource-agent
  uid: 1ebef947-504b-11e8-b75e-52540018543c
spec:
  minReadySeconds: 10
  replicas: 1
  revisionHistoryLimit: 5
  selector:
    matchLabels:
      qcloud-app: devex-resource-agent
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  template:
    metadata:
      creationTimestamp: ""
      labels:
        qcloud-app: devex-resource-agent
        qcloud-redeploy-timestamp: "1525771236"
    spec:
      containers:
      - name: devex-resource-agent
        image: ccr.ccs.tencentyun.com/eqxiu/devex-resource-agent:latest
        terminationMessagePath: /dev/termination-log
        imagePullPolicy: Always
        env:
        - name: PATH
          value: /usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
        - name: Agent_Nsq_Endpoint
          value: 192.168.1.12:4150
        - name: log_opt
          value: '"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10;
            --log-opt env=svcname;"'
        - name: svcname
          value: devex-resource-agent
        - name: LOGCHAIN_DRIVER
          value: influx
        - name: Agent_Influx_Endpoint
          value: 192.168.1.14:8089
        - name: Agent_Mongo_endpoint
          value: mdb1.yqxiu.cn:27010,mdb.yqxiu.cn:27010
        - name: Agent_Mongo_DB
          value: data-mgr
        - name: Agent_Influx_TCP_Endpoint
          value: http://192.168.1.14:8086
        resources:
          limits:
            cpu: 500m
            memory: 50Mi
          requests:
            cpu: 250m
            memory: 20Mi
        securityContext:
          privileged: false
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
      dnsPolicy: ClusterFirst
      securityContext: {}
      imagePullSecrets:
      - name: qcloudregistrykey
      serviceAccountName: ""
status:
  observedGeneration: 13
  replicas: 1
  updatedReplicas: 1
  availableReplicas: 1
`
	d, err := yaml.Marshal(&depBackup)

	assert.Nil(t, err)

	assert.Equal(t, expect, string(d))
}

func TestK8sServiceBackup(t *testing.T) {
	serBackup := K8sServiceBackup{
		Apiversion: "v1",
		Kind:       "Service",
		Metadata: K8sServiceBackup_Metadata{
			Annotations: map[string]string{
				"service.kubernetes.io/qcloud-loadbalancer-clusterid":         "cls-rfje0azd",
				"service.kubernetes.io/qcloud-loadbalancer-internal":          "96194",
				"service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov",
			},
			CreationTimestamp: "2018-05-02T06:33:10Z",
			Labels: map[string]string{
				"qcloud-app": "devex-agent",
			},
			Name:            "devex-agent",
			Namespace:       "scheduler",
			ResourceVersion: "10554019082",
			SelfLink:        "/api/v1/namespaces/scheduler/services/devex-agent",
			Uid:             "aee2f25e-4dd2-11e8-b75e-52540018543c",
		},
		Spec: K8sServiceInfo_items_spec{
			ClusterIP: "172.16.255.95",
			Ports: []K8sServiceInfo_items_spec_ports{
				{
					Name:       "tcp-8000-8000-3s2aq",
					NodePort:   30939,
					Port:       8000,
					Protocol:   "TCP",
					TargetPort: 8000,
				},
			},
			Selector: map[string]string{
				"qcloud-app": "devex-agent",
			},
			SessionAffinity: "None",
			Type:            "LoadBalancer",
		},
		Status: K8sServiceInfo_items_status{
			LoadBalancer: map[string][]map[string]string{
				"ingress": []map[string]string{
					{
						"ip": "192.168.1.40",
					},
				},
			},
		},
	}

	var expectd = `apiVersion: v1
kind: Service
metadata:
  annotations:
    service.kubernetes.io/qcloud-loadbalancer-clusterid: cls-rfje0azd
    service.kubernetes.io/qcloud-loadbalancer-internal: "96194"
    service.kubernetes.io/qcloud-loadbalancer-internal-subnetid: subnet-ba0hwkov
  creationTimestamp: 2018-05-02T06:33:10Z
  labels:
    qcloud-app: devex-agent
  name: devex-agent
  namespace: scheduler
  resourceVersion: "10554019082"
  selfLink: /api/v1/namespaces/scheduler/services/devex-agent
  uid: aee2f25e-4dd2-11e8-b75e-52540018543c
spec:
  ports:
  - name: tcp-8000-8000-3s2aq
    protocol: TCP
    port: 8000
    targetPort: 8000
    nodePort: 30939
  selector:
    qcloud-app: devex-agent
  clusterIP: 172.16.255.95
  type: LoadBalancer
  sessionAffinity: None
status:
  loadBalancer:
    ingress:
    - ip: 192.168.1.40
`
	d, err := yaml.Marshal(&serBackup)
	assert.Nil(t, err)
	assert.Equal(t, expectd, string(d))
}
