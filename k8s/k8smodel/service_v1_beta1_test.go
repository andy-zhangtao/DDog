package k8smodel

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.

func TestServiceV1Model(t *testing.T) {
	var data = `{
  "kind": "ServiceList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/namespaces/scheduler/services",
    "resourceVersion": "11346676879"
  },
  "items": [
    {
      "metadata": {
        "name": "caas-deploy-agent",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/caas-deploy-agent",
        "uid": "96a04c0f-4dd6-11e8-b75e-52540018543c",
        "resourceVersion": "10549846457",
        "creationTimestamp": "2018-05-02T07:01:08Z",
        "labels": {
          "qcloud-app": "caas-deploy-agent"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-9qd35",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 30133
          }
        ],
        "selector": {
          "qcloud-app": "caas-deploy-agent"
        },
        "clusterIP": "172.16.255.100",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.7"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "devex-agent",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/devex-agent",
        "uid": "aee2f25e-4dd2-11e8-b75e-52540018543c",
        "resourceVersion": "10554019082",
        "creationTimestamp": "2018-05-02T06:33:10Z",
        "labels": {
          "qcloud-app": "devex-agent"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-3s2aq",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 30939
          }
        ],
        "selector": {
          "qcloud-app": "devex-agent"
        },
        "clusterIP": "172.16.255.95",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.40"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "devex-gitlab-agent",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/devex-gitlab-agent",
        "uid": "ad72a729-4dd4-11e8-b75e-52540018543c",
        "resourceVersion": "10553197465",
        "creationTimestamp": "2018-05-02T06:47:27Z",
        "labels": {
          "qcloud-app": "devex-gitlab-agent"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-f1p1u",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 31690
          }
        ],
        "selector": {
          "qcloud-app": "devex-gitlab-agent"
        },
        "clusterIP": "172.16.255.236",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.37"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "devex-jenkins-agent",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/devex-jenkins-agent",
        "uid": "0eabbc42-4dd4-11e8-b75e-52540018543c",
        "resourceVersion": "10553163698",
        "creationTimestamp": "2018-05-02T06:43:00Z",
        "labels": {
          "qcloud-app": "devex-jenkins-agent"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-i0ult",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 31449
          }
        ],
        "selector": {
          "qcloud-app": "devex-jenkins-agent"
        },
        "clusterIP": "172.16.255.40",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.2"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "devex-resource-agent",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/devex-resource-agent",
        "uid": "d9ba35d5-5065-11e8-b75e-52540018543c",
        "resourceVersion": "10721383323",
        "creationTimestamp": "2018-05-05T13:11:40Z",
        "labels": {
          "qcloud-app": "devex-resource-agent"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-4kox4",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 31914
          }
        ],
        "selector": {
          "qcloud-app": "devex-resource-agent"
        },
        "clusterIP": "172.16.255.30",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.28"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "devex-ui",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/devex-ui",
        "uid": "6eae7cfb-4508-11e8-b75e-52540018543c",
        "resourceVersion": "10983396980",
        "creationTimestamp": "2018-04-21T02:05:15Z",
        "labels": {
          "qcloud-app": "devex-ui"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-80-80-cof0b",
            "protocol": "TCP",
            "port": 80,
            "targetPort": 80,
            "nodePort": 30009
          }
        ],
        "selector": {
          "qcloud-app": "devex-ui"
        },
        "clusterIP": "172.16.255.103",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.49"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "influxdb-chronograf",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/influxdb-chronograf",
        "uid": "d18a0219-4781-11e8-b75e-52540018543c",
        "resourceVersion": "9944316784",
        "creationTimestamp": "2018-04-24T05:39:12Z",
        "labels": {
          "qcloud-app": "influxdb-chronograf"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8888-8888-0h8jv",
            "protocol": "TCP",
            "port": 8888,
            "targetPort": 8888,
            "nodePort": 32617
          }
        ],
        "selector": {
          "qcloud-app": "influxdb-chronograf"
        },
        "clusterIP": "172.16.255.197",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.46"
            }
          ]
        }
      }
    },
    {
      "metadata": {
        "name": "scheduler",
        "namespace": "scheduler",
        "selfLink": "/api/v1/namespaces/scheduler/services/scheduler",
        "uid": "5747e6fd-0576-11e8-b75e-52540018543c",
        "resourceVersion": "9624045887",
        "creationTimestamp": "2018-01-30T04:30:46Z",
        "labels": {
          "qcloud-app": "scheduler"
        },
        "annotations": {
          "service.kubernetes.io/qcloud-loadbalancer-clusterid": "cls-rfje0azd",
          "service.kubernetes.io/qcloud-loadbalancer-internal": "96194",
          "service.kubernetes.io/qcloud-loadbalancer-internal-subnetid": "subnet-ba0hwkov"
        }
      },
      "spec": {
        "ports": [
          {
            "name": "tcp-8000-8000-qk5dj",
            "protocol": "TCP",
            "port": 8000,
            "targetPort": 8000,
            "nodePort": 31692
          }
        ],
        "selector": {
          "qcloud-app": "scheduler"
        },
        "clusterIP": "172.16.255.121",
        "type": "LoadBalancer",
        "sessionAffinity": "None",
        "externalTrafficPolicy": "Cluster"
      },
      "status": {
        "loadBalancer": {
          "ingress": [
            {
              "ip": "192.168.1.33"
            }
          ]
        }
      }
    }
  ]
}`

	var ks K8sServiceInfo

	err := json.Unmarshal([]byte(data), &ks)
	assert.Nil(t, err)
	assert.Equal(t, "caas-deploy-agent", ks.Items[0].Metadata.Labels["qcloud-app"])
	assert.Equal(t, "2018-05-02T07:01:08Z", ks.Items[0].Metadata.CreationTimestamp)
	assert.Equal(t, "subnet-ba0hwkov", ks.Items[0].Metadata.Annotations["service.kubernetes.io/qcloud-loadbalancer-internal-subnetid"])
	assert.Equal(t, "192.168.1.7", ks.Items[0].Status.LoadBalancer["ingress"][0]["ip"])
	assert.Equal(t, "192.168.1.33", ks.Items[len(ks.Items)-1].Status.LoadBalancer["ingress"][0]["ip"])
}
