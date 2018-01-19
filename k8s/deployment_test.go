package k8s

import (
	"testing"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestK8sMetaData_GetDeployMent(t *testing.T) {
	jsons := `{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/namespaces/devenv/pods",
    "resourceVersion": "5342202077"
  },
  "items": [
    {
      "metadata": {
        "name": "base-server-491048210-3fv8v",
        "generateName": "base-server-491048210-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/base-server-491048210-3fv8v",
        "uid": "78d63601-fa87-11e7-a645-52540018543c",
        "resourceVersion": "5265429285",
        "creationTimestamp": "2018-01-16T06:35:41Z",
        "labels": {
          "pod-template-hash": "491048210",
          "qcloud-app": "base-server",
          "qcloud-redeploy-timestamp": "1515988415"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"base-server-491048210\",\"uid\":\"a9261988-f9a7-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5265314905\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "base-server-491048210",
            "uid": "a9261988-f9a7-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "base-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "workingDir": "/app",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/jvm/java-1.8-openjdk/jre/bin:/usr/lib/jvm/java-1.8-openjdk/bin"
              },
              {
                "name": "LANG",
                "value": "C.UTF-8"
              },
              {
                "name": "JAVA_HOME",
                "value": "/usr/lib/jvm/default-jvm/jre"
              },
              {
                "name": "JAVA_VERSION",
                "value": "8u111"
              },
              {
                "name": "JAVA_ALPINE_VERSION",
                "value": "8.111.14-r0"
              }
            ],
            "resources": {
              "limits": {
                "cpu": "200m"
              },
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.15",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:43Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          }
        ],
        "hostIP": "192.168.1.15",
        "podIP": "172.16.0.103",
        "startTime": "2018-01-16T06:35:41Z",
        "containerStatuses": [
          {
            "name": "base-server",
            "state": {
              "running": {
                "startedAt": "2018-01-16T06:35:42Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server@sha256:2f2e0b691ada35ab7fc225305c11a893f311c9723e1741aeb8cb65ab4471aa62",
            "containerID": "docker://6c1747aca768fc124d6a8f71f1f76173b361a1fd40709f49093d76a43d12ac77"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "scene-service-provider-eqshow-cn-4247796272-pg1tc",
        "generateName": "scene-service-provider-eqshow-cn-4247796272-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/scene-service-provider-eqshow-cn-4247796272-pg1tc",
        "uid": "79f37d64-f9c0-11e7-a645-52540018543c",
        "resourceVersion": "5229280702",
        "creationTimestamp": "2018-01-15T06:51:13Z",
        "labels": {
          "pod-template-hash": "4247796272",
          "qcloud-app": "scene-service-provider-eqshow-cn",
          "qcloud-redeploy-timestamp": "1515999073"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"scene-service-provider-eqshow-cn-4247796272\",\"uid\":\"79f17a29-f9c0-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5229279803\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "scene-service-provider-eqshow-cn-4247796272",
            "uid": "79f17a29-f9c0-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "scene-service-provider-eqshow-cn",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-scene-service-provider:1.01",
            "workingDir": "/app",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/jvm/java-1.8-openjdk/jre/bin:/usr/lib/jvm/java-1.8-openjdk/bin"
              },
              {
                "name": "LANG",
                "value": "C.UTF-8"
              },
              {
                "name": "JAVA_HOME",
                "value": "/usr/lib/jvm/default-jvm/jre"
              },
              {
                "name": "JAVA_VERSION",
                "value": "8u111"
              },
              {
                "name": "JAVA_ALPINE_VERSION",
                "value": "8.111.14-r0"
              },
              {
                "name": "DOMAIN",
                "value": "scene.service.provider.eqshow.cn"
              }
            ],
            "resources": {
              "limits": {
                "cpu": "500m",
                "memory": "1Gi"
              },
              "requests": {
                "cpu": "250m",
                "memory": "256Mi"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.15",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-15T06:51:13Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-15T06:51:15Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-15T06:51:13Z"
          }
        ],
        "hostIP": "192.168.1.15",
        "podIP": "172.16.0.100",
        "startTime": "2018-01-15T06:51:13Z",
        "containerStatuses": [
          {
            "name": "scene-service-provider-eqshow-cn",
            "state": {
              "running": {
                "startedAt": "2018-01-15T06:51:14Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-scene-service-provider:1.01",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-scene-service-provider@sha256:8615589079bb75fd19ebeee9a079439773252438eb113b886a67d1e25b590af5",
            "containerID": "docker://6b9cefec60e3dd64f1b9a12391034017f0c349ca2512b2519b76b5c3638b8703"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "test1-4227749671-b518f",
        "generateName": "test1-4227749671-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/test1-4227749671-b518f",
        "uid": "560c5139-fc1f-11e7-a645-52540018543c",
        "resourceVersion": "5341785219",
        "creationTimestamp": "2018-01-18T07:15:17Z",
        "labels": {
          "pod-template-hash": "4227749671",
          "qcloud-app": "test1",
          "qcloud-redeploy-timestamp": "1516259717"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"test1-4227749671\",\"uid\":\"560732c2-fc1f-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5341767734\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "test1-4227749671",
            "uid": "560732c2-fc1f-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "eqxiu-base-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "resources": {},
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "livenessProbe": {
              "tcpSocket": {
                "port": 40000
              },
              "initialDelaySeconds": 30,
              "timeoutSeconds": 5,
              "periodSeconds": 10,
              "successThreshold": 1,
              "failureThreshold": 5
            },
            "readinessProbe": {
              "tcpSocket": {
                "port": 40000
              },
              "initialDelaySeconds": 30,
              "timeoutSeconds": 5,
              "periodSeconds": 10,
              "successThreshold": 1,
              "failureThreshold": 5
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.16",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:15:17Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:15:57Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:15:17Z"
          }
        ],
        "hostIP": "192.168.1.16",
        "podIP": "172.16.1.115",
        "startTime": "2018-01-18T07:15:17Z",
        "containerStatuses": [
          {
            "name": "eqxiu-base-server",
            "state": {
              "running": {
                "startedAt": "2018-01-18T07:15:24Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server@sha256:687f74f6de750e1ead38f24c88fb337b3b6923b965bea77c68617943776d27fc",
            "containerID": "docker://a7d74ab22da0536ffed5b285f9fb863342948f2d77730453a18f4cff3d990a15"
          }
        ],
        "qosClass": "BestEffort"
      }
    },
    {
      "metadata": {
        "name": "test1-4227749671-vq243",
        "generateName": "test1-4227749671-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/test1-4227749671-vq243",
        "uid": "6df5a5cf-fc1f-11e7-a645-52540018543c",
        "resourceVersion": "5341802937",
        "creationTimestamp": "2018-01-18T07:15:57Z",
        "labels": {
          "pod-template-hash": "4227749671",
          "qcloud-app": "test1",
          "qcloud-redeploy-timestamp": "1516259717"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"test1-4227749671\",\"uid\":\"560732c2-fc1f-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5341785252\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "test1-4227749671",
            "uid": "560732c2-fc1f-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "eqxiu-base-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "resources": {},
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "livenessProbe": {
              "tcpSocket": {
                "port": 40000
              },
              "initialDelaySeconds": 30,
              "timeoutSeconds": 5,
              "periodSeconds": 10,
              "successThreshold": 1,
              "failureThreshold": 5
            },
            "readinessProbe": {
              "tcpSocket": {
                "port": 40000
              },
              "initialDelaySeconds": 30,
              "timeoutSeconds": 5,
              "periodSeconds": 10,
              "successThreshold": 1,
              "failureThreshold": 5
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.5",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:15:57Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:16:37Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-18T07:15:57Z"
          }
        ],
        "hostIP": "192.168.1.5",
        "podIP": "172.16.2.153",
        "startTime": "2018-01-18T07:15:57Z",
        "containerStatuses": [
          {
            "name": "eqxiu-base-server",
            "state": {
              "running": {
                "startedAt": "2018-01-18T07:15:59Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server@sha256:687f74f6de750e1ead38f24c88fb337b3b6923b965bea77c68617943776d27fc",
            "containerID": "docker://c087ce4adcd2080f4daec6c2e695802ab9887f6f402c110cb130884af1092338"
          }
        ],
        "qosClass": "BestEffort"
      }
    },
    {
      "metadata": {
        "name": "user-server-1369396299-7fwx7",
        "generateName": "user-server-1369396299-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/user-server-1369396299-7fwx7",
        "uid": "78d99439-fa87-11e7-a645-52540018543c",
        "resourceVersion": "5265429340",
        "creationTimestamp": "2018-01-16T06:35:41Z",
        "labels": {
          "pod-template-hash": "1369396299",
          "qcloud-app": "user-server",
          "qcloud-redeploy-timestamp": "1515998228"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"user-server-1369396299\",\"uid\":\"82ade022-f9be-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5265314918\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "user-server-1369396299",
            "uid": "82ade022-f9be-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "user-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server:1.2.0",
            "workingDir": "/app",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/jvm/java-1.8-openjdk/jre/bin:/usr/lib/jvm/java-1.8-openjdk/bin"
              },
              {
                "name": "LANG",
                "value": "C.UTF-8"
              },
              {
                "name": "JAVA_HOME",
                "value": "/usr/lib/jvm/default-jvm/jre"
              },
              {
                "name": "JAVA_VERSION",
                "value": "8u111"
              },
              {
                "name": "JAVA_ALPINE_VERSION",
                "value": "8.111.14-r0"
              }
            ],
            "resources": {
              "limits": {
                "cpu": "2",
                "memory": "1Gi"
              },
              "requests": {
                "cpu": "200m",
                "memory": "512Mi"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.16",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:43Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          }
        ],
        "hostIP": "192.168.1.16",
        "podIP": "172.16.1.112",
        "startTime": "2018-01-16T06:35:41Z",
        "containerStatuses": [
          {
            "name": "user-server",
            "state": {
              "running": {
                "startedAt": "2018-01-16T06:35:43Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server:1.2.0",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server@sha256:fd34a6f916260cbd91874331d67499d940260f73e071311dce9b3643bdcdf617",
            "containerID": "docker://8f6ceeabed5d2d6e667f0ea6bfb92b2de176fccfd79bf30c207df22f7d527c39"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "user-service-2005460954-df4k3",
        "generateName": "user-service-2005460954-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/user-service-2005460954-df4k3",
        "uid": "bc3bb618-fb38-11e7-a645-52540018543c",
        "resourceVersion": "5298445648",
        "creationTimestamp": "2018-01-17T03:44:35Z",
        "labels": {
          "pod-template-hash": "2005460954",
          "qcloud-app": "user-service",
          "qcloud-redeploy-timestamp": "1515997697"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"user-service-2005460954\",\"uid\":\"45d7889c-f9bd-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5228704827\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "user-service-2005460954",
            "uid": "45d7889c-f9bd-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "user-service",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-service:1.2.0",
            "workingDir": "/app",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/jvm/java-1.8-openjdk/jre/bin:/usr/lib/jvm/java-1.8-openjdk/bin"
              },
              {
                "name": "LANG",
                "value": "C.UTF-8"
              },
              {
                "name": "JAVA_HOME",
                "value": "/usr/lib/jvm/default-jvm/jre"
              },
              {
                "name": "JAVA_VERSION",
                "value": "8u111"
              },
              {
                "name": "JAVA_ALPINE_VERSION",
                "value": "8.111.14-r0"
              }
            ],
            "resources": {
              "limits": {
                "cpu": "200m"
              },
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.5",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-17T03:44:35Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-17T03:44:41Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-17T03:44:35Z"
          }
        ],
        "hostIP": "192.168.1.5",
        "podIP": "172.16.2.127",
        "startTime": "2018-01-17T03:44:35Z",
        "containerStatuses": [
          {
            "name": "user-service",
            "state": {
              "running": {
                "startedAt": "2018-01-17T03:44:41Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-service:1.2.0",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-service@sha256:b0a5e9c52004d66bd3d57d268f0bedc6ae4623344ebe0845ee92f0ab93461141",
            "containerID": "docker://a294aa819b5b65554a3640c8e018431b0bd7927b02fcef50df2d5ef244aa2b51"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "user-service-eqshow-cn-2826347142-zvcd3",
        "generateName": "user-service-eqshow-cn-2826347142-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/user-service-eqshow-cn-2826347142-zvcd3",
        "uid": "78dee429-fa87-11e7-a645-52540018543c",
        "resourceVersion": "5284915033",
        "creationTimestamp": "2018-01-16T06:35:41Z",
        "labels": {
          "pod-template-hash": "2826347142",
          "qcloud-app": "user-service-eqshow-cn"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"user-service-eqshow-cn-2826347142\",\"uid\":\"b3605516-f670-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5265314926\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "user-service-eqshow-cn-2826347142",
            "uid": "b3605516-f670-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "user-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server:1.2.0",
            "resources": {},
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.12",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T19:04:18Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          }
        ],
        "hostIP": "192.168.1.12",
        "podIP": "172.16.3.91",
        "startTime": "2018-01-16T06:35:41Z",
        "containerStatuses": [
          {
            "name": "user-server",
            "state": {
              "running": {
                "startedAt": "2018-01-16T19:04:18Z"
              }
            },
            "lastState": {
              "terminated": {
                "exitCode": 0,
                "reason": "Completed",
                "startedAt": "2018-01-16T18:58:53Z",
                "finishedAt": "2018-01-16T18:59:07Z",
                "containerID": "docker://fd936bc7d01fbd3e4e539bd73e38b03cd4f6be354fbc1e79d521db96ab3a4e32"
              }
            },
            "ready": true,
            "restartCount": 144,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server:1.2.0",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-user-server@sha256:d3d85eed6918669b7bebeef8a38a74ae4885b78aacaca16f04c73b3e92f86159",
            "containerID": "docker://5a6c594283339379b430700313986238f32e665bd27f1115c8c40e0512d1bc53"
          }
        ],
        "qosClass": "BestEffort"
      }
    },
    {
      "metadata": {
        "name": "webserver-671749410-0133h",
        "generateName": "webserver-671749410-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/webserver-671749410-0133h",
        "uid": "4d62d94e-e62a-11e7-a645-52540018543c",
        "resourceVersion": "4401700168",
        "creationTimestamp": "2017-12-21T08:38:21Z",
        "labels": {
          "pod-template-hash": "671749410",
          "qcloud-app": "webserver"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"webserver-671749410\",\"uid\":\"4d61a5eb-e62a-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"4401699347\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "webserver-671749410",
            "uid": "4d61a5eb-e62a-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "webserver",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "TZ",
                "value": "Asia/Shanghai"
              }
            ],
            "resources": {
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.15",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:38:21Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:38:23Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:38:21Z"
          }
        ],
        "hostIP": "192.168.1.15",
        "podIP": "172.16.0.13",
        "startTime": "2017-12-21T08:38:21Z",
        "containerStatuses": [
          {
            "name": "webserver",
            "state": {
              "running": {
                "startedAt": "2017-12-21T08:38:23Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver@sha256:c873355b990107545aad170426f985f8c3cf0aebdef19f58d7992b553067bec0",
            "containerID": "docker://7999efe3b1492d5a29252688f576ab4060b787cd5f4288486ffbb24c7995b009"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "webserver-671749410-3dnpg",
        "generateName": "webserver-671749410-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/webserver-671749410-3dnpg",
        "uid": "7752700d-e62a-11e7-a645-52540018543c",
        "resourceVersion": "4401725484",
        "creationTimestamp": "2017-12-21T08:39:32Z",
        "labels": {
          "pod-template-hash": "671749410",
          "qcloud-app": "webserver"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"webserver-671749410\",\"uid\":\"4d61a5eb-e62a-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"4401724267\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "webserver-671749410",
            "uid": "4d61a5eb-e62a-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "webserver",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "TZ",
                "value": "Asia/Shanghai"
              }
            ],
            "resources": {
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.16",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:35Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          }
        ],
        "hostIP": "192.168.1.16",
        "podIP": "172.16.1.24",
        "startTime": "2017-12-21T08:39:32Z",
        "containerStatuses": [
          {
            "name": "webserver",
            "state": {
              "running": {
                "startedAt": "2017-12-21T08:39:34Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver@sha256:c873355b990107545aad170426f985f8c3cf0aebdef19f58d7992b553067bec0",
            "containerID": "docker://fe1a7cc3093a6b5cd12791e0228cb26fad614cd39c882505e5ca77023ad27d44"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "webserver-671749410-4t7j4",
        "generateName": "webserver-671749410-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/webserver-671749410-4t7j4",
        "uid": "78e3047d-fa87-11e7-a645-52540018543c",
        "resourceVersion": "5265431917",
        "creationTimestamp": "2018-01-16T06:35:41Z",
        "labels": {
          "pod-template-hash": "671749410",
          "qcloud-app": "webserver"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"webserver-671749410\",\"uid\":\"4d61a5eb-e62a-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5265314933\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "webserver-671749410",
            "uid": "4d61a5eb-e62a-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "webserver",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "TZ",
                "value": "Asia/Shanghai"
              }
            ],
            "resources": {
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.12",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:49Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-16T06:35:41Z"
          }
        ],
        "hostIP": "192.168.1.12",
        "podIP": "172.16.3.92",
        "startTime": "2018-01-16T06:35:41Z",
        "containerStatuses": [
          {
            "name": "webserver",
            "state": {
              "running": {
                "startedAt": "2018-01-16T06:35:48Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver@sha256:c873355b990107545aad170426f985f8c3cf0aebdef19f58d7992b553067bec0",
            "containerID": "docker://3d89224de58a992b7f623d58c96c50187d92e4c7c2e19462b9456b667097183d"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "webserver-671749410-x35p3",
        "generateName": "webserver-671749410-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/webserver-671749410-x35p3",
        "uid": "77527a75-e62a-11e7-a645-52540018543c",
        "resourceVersion": "4401725157",
        "creationTimestamp": "2017-12-21T08:39:32Z",
        "labels": {
          "pod-template-hash": "671749410",
          "qcloud-app": "webserver"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"webserver-671749410\",\"uid\":\"4d61a5eb-e62a-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"4401724267\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "webserver-671749410",
            "uid": "4d61a5eb-e62a-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "webserver",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "TZ",
                "value": "Asia/Shanghai"
              }
            ],
            "resources": {
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.16",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:34Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          }
        ],
        "hostIP": "192.168.1.16",
        "podIP": "172.16.1.23",
        "startTime": "2017-12-21T08:39:32Z",
        "containerStatuses": [
          {
            "name": "webserver",
            "state": {
              "running": {
                "startedAt": "2017-12-21T08:39:34Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver@sha256:c873355b990107545aad170426f985f8c3cf0aebdef19f58d7992b553067bec0",
            "containerID": "docker://ba8886cdf67c045361889ace03ff8f7765323abb5bcb42098c4fe8e8398b8731"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "webserver-671749410-zmp5c",
        "generateName": "webserver-671749410-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/webserver-671749410-zmp5c",
        "uid": "775288ad-e62a-11e7-a645-52540018543c",
        "resourceVersion": "4401725047",
        "creationTimestamp": "2017-12-21T08:39:32Z",
        "labels": {
          "pod-template-hash": "671749410",
          "qcloud-app": "webserver"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"webserver-671749410\",\"uid\":\"4d61a5eb-e62a-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"4401724267\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "webserver-671749410",
            "uid": "4d61a5eb-e62a-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "webserver",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "TZ",
                "value": "Asia/Shanghai"
              }
            ],
            "resources": {
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.12",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:34Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2017-12-21T08:39:32Z"
          }
        ],
        "hostIP": "192.168.1.12",
        "podIP": "172.16.3.24",
        "startTime": "2017-12-21T08:39:32Z",
        "containerStatuses": [
          {
            "name": "webserver",
            "state": {
              "running": {
                "startedAt": "2017-12-21T08:39:33Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver:0.4",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-webserver@sha256:c873355b990107545aad170426f985f8c3cf0aebdef19f58d7992b553067bec0",
            "containerID": "docker://276ab234d785d096782ae42e30b0160dd2ba384481400873c6ca499438037ac9"
          }
        ],
        "qosClass": "Burstable"
      }
    },
    {
      "metadata": {
        "name": "www-helloworld-eqshow-cn-592900821-5b44p",
        "generateName": "www-helloworld-eqshow-cn-592900821-",
        "namespace": "devenv",
        "selfLink": "/api/v1/namespaces/devenv/pods/www-helloworld-eqshow-cn-592900821-5b44p",
        "uid": "c2bba0fa-f674-11e7-a645-52540018543c",
        "resourceVersion": "5080903686",
        "creationTimestamp": "2018-01-11T02:11:40Z",
        "labels": {
          "pod-template-hash": "592900821",
          "qcloud-app": "www-helloworld-eqshow-cn",
          "qcloud-redeploy-timestamp": "1515636700"
        },
        "annotations": {
          "kubernetes.io/created-by": "{\"kind\":\"SerializedReference\",\"apiVersion\":\"v1\",\"reference\":{\"kind\":\"ReplicaSet\",\"namespace\":\"devenv\",\"name\":\"www-helloworld-eqshow-cn-592900821\",\"uid\":\"c2b9f5c4-f674-11e7-a645-52540018543c\",\"apiVersion\":\"extensions\",\"resourceVersion\":\"5080902492\"}}\n"
        },
        "ownerReferences": [
          {
            "apiVersion": "extensions/v1beta1",
            "kind": "ReplicaSet",
            "name": "www-helloworld-eqshow-cn-592900821",
            "uid": "c2b9f5c4-f674-11e7-a645-52540018543c",
            "controller": true,
            "blockOwnerDeletion": true
          }
        ]
      },
      "spec": {
        "volumes": [
          {
            "name": "default-token-mmlt1",
            "secret": {
              "secretName": "default-token-mmlt1",
              "defaultMode": 420
            }
          }
        ],
        "containers": [
          {
            "name": "www-helloworld-eqshow-cn",
            "image": "ccr.ccs.tencentyun.com/eqxiu/hello-world:latest",
            "env": [
              {
                "name": "PATH",
                "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
              },
              {
                "name": "NPM_CONFIG_LOGLEVEL",
                "value": "info"
              },
              {
                "name": "NODE_VERSION",
                "value": "4.2.6"
              }
            ],
            "resources": {
              "limits": {
                "cpu": "200m"
              },
              "requests": {
                "cpu": "200m"
              }
            },
            "volumeMounts": [
              {
                "name": "default-token-mmlt1",
                "readOnly": true,
                "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
              }
            ],
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File",
            "imagePullPolicy": "Always",
            "securityContext": {
              "privileged": false
            }
          }
        ],
        "restartPolicy": "Always",
        "terminationGracePeriodSeconds": 30,
        "dnsPolicy": "ClusterFirst",
        "serviceAccountName": "default",
        "serviceAccount": "default",
        "nodeName": "192.168.1.15",
        "securityContext": {},
        "imagePullSecrets": [
          {
            "name": "qcloudregistrykey"
          }
        ],
        "schedulerName": "default-scheduler"
      },
      "status": {
        "phase": "Running",
        "conditions": [
          {
            "type": "Initialized",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-11T02:11:40Z"
          },
          {
            "type": "Ready",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-11T02:11:42Z"
          },
          {
            "type": "PodScheduled",
            "status": "True",
            "lastProbeTime": null,
            "lastTransitionTime": "2018-01-11T02:11:40Z"
          }
        ],
        "hostIP": "192.168.1.15",
        "podIP": "172.16.0.79",
        "startTime": "2018-01-11T02:11:40Z",
        "containerStatuses": [
          {
            "name": "www-helloworld-eqshow-cn",
            "state": {
              "running": {
                "startedAt": "2018-01-11T02:11:41Z"
              }
            },
            "lastState": {},
            "ready": true,
            "restartCount": 0,
            "image": "ccr.ccs.tencentyun.com/eqxiu/hello-world:latest",
            "imageID": "docker-pullable://ccr.ccs.tencentyun.com/eqxiu/hello-world@sha256:b9bac7a91a8a2a55ef3902b02a23f5221fcecf334f90156e557a45b1668087f5",
            "containerID": "docker://00561f2f86ff990bc3d7ad964ed4d87fb195325f5be9f62a5836a50034fae42d"
          }
        ],
        "qosClass": "Burstable"
      }
    }
  ]
}`
	var kapi k8smodel.K8s
	err := json.Unmarshal([]byte(jsons), &kapi)
	assert.Nil(t, err)
	assert.Equal(t, "test1", kapi.Items[2].Metadata.Labels.Qcloud_app, "Qcloud_app Error")
	assert.Equal(t, "Running", kapi.Items[2].Status.Phase, "Qcloud_app Error")
	assert.Equal(t, "test1", kapi.Items[3].Metadata.Labels.Qcloud_app, "Qcloud_app Error")
	assert.Equal(t, "Running", kapi.Items[3].Status.Phase, "Qcloud_app Error")
}
