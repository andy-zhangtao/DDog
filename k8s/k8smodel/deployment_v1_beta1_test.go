package k8smodel

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/14.
func TestDeployMentV1Model(t *testing.T) {
	var data = `
{
  "kind": "DeploymentList",
  "apiVersion": "apps/v1beta1",
  "metadata": {
    "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments",
    "resourceVersion": "11335571561"
  },
  "items": [
    {
      "metadata": {
        "name": "caas-deploy-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/caas-deploy-agent",
        "uid": "969d2e57-4dd6-11e8-b75e-52540018543c",
        "resourceVersion": "11330488116",
        "generation": 10,
        "creationTimestamp": "2018-05-02T07:01:08Z",
        "labels": {
          "qcloud-app": "caas-deploy-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "10",
          "description": "The GraphQL Of Caas"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "caas-deploy-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "caas-deploy-agent",
              "qcloud-redeploy-timestamp": "1526275108"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "caas-deploy-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_CLUSTER_ID",
                    "value": "cls-rfje0azd"
                  },
                  {
                    "name": "DDOG_REGION",
                    "value": "sh"
                  },
                  {
                    "name": "DDOG_LOG_OPT",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=100;             --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "svcname",
                    "value": "caas-deploy-agent"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 10,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-14T06:50:46Z",
            "lastTransitionTime": "2018-05-14T06:50:46Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-agent",
        "uid": "aedfc96e-4dd2-11e8-b75e-52540018543c",
        "resourceVersion": "11129795624",
        "generation": 53,
        "creationTimestamp": "2018-05-02T06:33:10Z",
        "labels": {
          "qcloud-app": "devex-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "53"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-agent",
              "qcloud-redeploy-timestamp": "1526004612"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "192.168.0.17:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "svcname",
                    "value": "devex-agent"
                  },
                  {
                    "name": "Agent_Jenkins_Endpoint",
                    "value": "http://build.yqxiu.cn:18080"
                  },
                  {
                    "name": "Agent_Jenkins_User",
                    "value": "admin"
                  },
                  {
                    "name": "Agent_Jenkins_Passwd",
                    "value": "admin"
                  },
                  {
                    "name": "Agent_Caas_Endpoint",
                    "value": "http://192.168.1.7:8000/api"
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "Agent_Influx_Endpoint",
                    "value": "http://192.168.1.14:8086"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 53,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-11T09:55:09Z",
            "lastTransitionTime": "2018-05-11T09:55:09Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-caas-deploy-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-caas-deploy-agent",
        "uid": "072edfdd-4dd2-11e8-b75e-52540018543c",
        "resourceVersion": "11125147697",
        "generation": 16,
        "creationTimestamp": "2018-05-02T06:28:29Z",
        "labels": {
          "qcloud-app": "devex-caas-deploy-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "15"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-caas-deploy-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-caas-deploy-agent",
              "qcloud-redeploy-timestamp": "1526020048"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-caas-deploy-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-caas-deploy-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "Agent_Caas_Endpoint",
                    "value": "http://192.168.1.7:8000/api"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "svcname",
                    "value": "devex-caas-deploy-agent"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "192.168.0.17:27010"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 16,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-11T08:17:49Z",
            "lastTransitionTime": "2018-05-11T08:17:49Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-gitlab-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-gitlab-agent",
        "uid": "ad6f4d19-4dd4-11e8-b75e-52540018543c",
        "resourceVersion": "10913895643",
        "generation": 5,
        "creationTimestamp": "2018-05-02T06:47:27Z",
        "labels": {
          "qcloud-app": "devex-gitlab-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "5",
          "description": "Get Kong Configure From Gitlab"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-gitlab-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-gitlab-agent",
              "qcloud-redeploy-timestamp": "1525771248"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-gitlab-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/gitlab-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "192.168.0.17:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "svcname",
                    "value": "devex-gitlab-agent"
                  },
                  {
                    "name": "Agent_Gitlab_Endpoint",
                    "value": "http://gitlab.yqxiu.cn"
                  },
                  {
                    "name": "Agent_Gitlab_Token",
                    "value": "VD4E18spcC7_WFMyniAz"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 5,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T08:45:31Z",
            "lastTransitionTime": "2018-05-02T08:45:31Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-jenkins-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-jenkins-agent",
        "uid": "0ea84b7e-4dd4-11e8-b75e-52540018543c",
        "resourceVersion": "11129842158",
        "generation": 30,
        "creationTimestamp": "2018-05-02T06:43:00Z",
        "labels": {
          "qcloud-app": "devex-jenkins-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "30",
          "description": "Connect To Jenkins Server"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-jenkins-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-jenkins-agent",
              "qcloud-redeploy-timestamp": "1526024127"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-jenkins-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-jenkins-build-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "192.168.0.17:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "svcname",
                    "value": "devex-jenkins-agent"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "Agent_Jenkins_Endpoint",
                    "value": "http://build.yqxiu.cn:18080"
                  },
                  {
                    "name": "Agent_Jenkins_User",
                    "value": "admin"
                  },
                  {
                    "name": "Agent_Jenkins_Passwd",
                    "value": "admin"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 30,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-11T09:56:07Z",
            "lastTransitionTime": "2018-05-11T09:56:07Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-kong-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-kong-agent",
        "uid": "d76dbec1-4eb8-11e8-b75e-52540018543c",
        "resourceVersion": "11033681978",
        "generation": 8,
        "creationTimestamp": "2018-05-03T10:00:43Z",
        "labels": {
          "qcloud-app": "devex-kong-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "8",
          "description": "Update Kong Info"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-kong-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-kong-agent",
              "qcloud-redeploy-timestamp": "1525919762"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-kong-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-kong-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "mdb1.yqxiu.cn:27010,mdb.yqxiu.cn:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "Agent_Gitlab_Agent_Endpoint",
                    "value": "http://192.168.0.4:18000/api"
                  },
                  {
                    "name": "Agent_Kong_Endpoint",
                    "value": "http://192.168.0.4:8001"
                  },
                  {
                    "name": "svcname",
                    "value": "devex-kong-agent"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 8,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-03T10:00:46Z",
            "lastTransitionTime": "2018-05-03T10:00:46Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-resource-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-resource-agent",
        "uid": "1ebef947-504b-11e8-b75e-52540018543c",
        "resourceVersion": "10913887877",
        "generation": 13,
        "creationTimestamp": "2018-05-05T10:00:20Z",
        "labels": {
          "qcloud-app": "devex-resource-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "13",
          "description": "Collect Container Resource"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-resource-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-resource-agent",
              "qcloud-redeploy-timestamp": "1525771236"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-resource-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-resource-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "svcname",
                    "value": "devex-resource-agent"
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "Agent_Influx_Endpoint",
                    "value": "192.168.1.14:8089"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "mdb1.yqxiu.cn:27010,mdb.yqxiu.cn:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "Agent_Influx_TCP_Endpoint",
                    "value": "http://192.168.1.14:8086"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 13,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-05T12:13:40Z",
            "lastTransitionTime": "2018-05-05T12:13:40Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-status-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-status-agent",
        "uid": "25bfaf46-4dcd-11e8-b75e-52540018543c",
        "resourceVersion": "10913917164",
        "generation": 19,
        "creationTimestamp": "2018-05-02T05:53:33Z",
        "labels": {
          "qcloud-app": "devex-status-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "19",
          "description": "Receive Message From DDog. Update Request Deploy Status"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-status-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-status-agent",
              "qcloud-redeploy-timestamp": "1525771271"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-status-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/devex-status-agent:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "Agent_Mongo_endpoint",
                    "value": "192.168.0.17:27010"
                  },
                  {
                    "name": "Agent_Mongo_DB",
                    "value": "data-mgr"
                  },
                  {
                    "name": "Agent_Nsq_Endpoint",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "svcname",
                    "value": "devex_status_agent"
                  },
                  {
                    "name": "Agent_Caas_Endpoint",
                    "value": "http://192.168.1.7:8000/api"
                  }
                ],
                "resources": {
                  "limits": {
                    "memory": "50Mi"
                  },
                  "requests": {
                    "memory": "20Mi"
                  }
                },
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File",
                "imagePullPolicy": "Always",
                "securityContext": {
                  "privileged": false
                }
              },
              {
                "name": "sidecar",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:latest",
                "env": [
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_AGENT_SPIDER_SVC",
                    "value": "scene-web"
                  },
                  {
                    "name": "DDOG_AGENT_NAME",
                    "value": "SpiderAgent"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "DDOG_AGENT_SPIDER_NS",
                    "value": "eqxiu"
                  },
                  {
                    "name": "svcname",
                    "value": "scene-web_sidecar"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt gelf-address=udp://123.206.176.139:9988; --log-opt buf=100; --log-opt env=svcname;\""
                  },
                  {
                    "name": "DDOG_AGENT_SPIDER_PORT",
                    "value": "29200"
                  }
                ],
                "resources": {
                  "limits": {
                    "memory": "20Mi"
                  },
                  "requests": {
                    "memory": "10Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 19,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T09:01:45Z",
            "lastTransitionTime": "2018-05-02T09:01:45Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "devex-ui",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/devex-ui",
        "uid": "6eab5af8-4508-11e8-b75e-52540018543c",
        "resourceVersion": "11313801594",
        "generation": 30,
        "creationTimestamp": "2018-04-21T02:05:15Z",
        "labels": {
          "qcloud-app": "devex-ui"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "30",
          "description": "The UI of DevEx platform"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "devex-ui"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "devex-ui",
              "qcloud-redeploy-timestamp": "1526259817"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "devex-ui",
                "image": "ccr.ccs.tencentyun.com/eqxiu/caas-devex-ui:latest",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "NODE_VERSION",
                    "value": "9.11.1"
                  },
                  {
                    "name": "YARN_VERSION",
                    "value": "1.5.1"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "50Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 30,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T09:01:27Z",
            "lastTransitionTime": "2018-05-02T09:01:27Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "influxdb-chronograf",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/influxdb-chronograf",
        "uid": "d186e0d0-4781-11e8-b75e-52540018543c",
        "resourceVersion": "10508483303",
        "generation": 2,
        "creationTimestamp": "2018-04-24T05:39:12Z",
        "labels": {
          "qcloud-app": "influxdb-chronograf"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "2"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "influxdb-chronograf"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "influxdb-chronograf",
              "qcloud-redeploy-timestamp": "1525253773"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "influxdb-chronograf",
                "image": "chronograf:alpine",
                "args": [
                  "--influxdb-url=http://192.168.1.14:8086"
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "256Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "128Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 2,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T08:45:39Z",
            "lastTransitionTime": "2018-05-02T08:45:39Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "scheduler",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/scheduler",
        "uid": "57444d1a-0576-11e8-b75e-52540018543c",
        "resourceVersion": "10508505556",
        "generation": 64,
        "creationTimestamp": "2018-01-30T04:30:46Z",
        "labels": {
          "qcloud-app": "scheduler"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "55"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "scheduler"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "scheduler",
              "qcloud-redeploy-timestamp": "1525253791"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "scheduler",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler:0.6.10",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "DDOG_DEBUG",
                    "value": "true"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_NAME_SPACE",
                    "value": "eqxiu"
                  },
                  {
                    "name": "DDOG_LOG_OPT",
                    "value": "\"--log-opt gelf-address=udp://123.206.176.139:9988; --log-opt buf=100; --log-opt env=svcname;\""
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt gelf-address=udp://123.206.176.139:9988; --log-opt buf=100; --log-opt env=svcname;\""
                  },
                  {
                    "name": "svcname",
                    "value": "scheduler"
                  },
                  {
                    "name": "DDOG_REGION",
                    "value": "sh"
                  },
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "20m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "10m",
                    "memory": "30Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 64,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T08:45:55Z",
            "lastTransitionTime": "2018-05-02T08:45:55Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "scheduler-agent-1517902615",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/scheduler-agent-1517902615",
        "uid": "8206562e-0b10-11e8-b75e-52540018543c",
        "resourceVersion": "11032072477",
        "generation": 17,
        "creationTimestamp": "2018-02-06T07:36:56Z",
        "labels": {
          "qcloud-app": "scheduler-agent-1517902615"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "15"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "scheduler-agent-1517902615"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "scheduler-agent-1517902615",
              "qcloud-redeploy-timestamp": "1525917782"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:0.6.9",
                "env": [
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "svcname",
                    "value": "scheduler-destory-agent"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10; --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "10m",
                    "memory": "20Mi"
                  },
                  "requests": {
                    "cpu": "10m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 17,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-10T01:43:41Z",
            "lastTransitionTime": "2018-05-10T01:43:41Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "scheduler-deploy-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/scheduler-deploy-agent",
        "uid": "852a6d59-4dd5-11e8-b75e-52540018543c",
        "resourceVersion": "10551936157",
        "generation": 4,
        "creationTimestamp": "2018-05-02T06:53:29Z",
        "labels": {
          "qcloud-app": "scheduler-deploy-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "4"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "scheduler-deploy-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "scheduler-deploy-agent",
              "qcloud-redeploy-timestamp": "1525310444"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "scheduler-deploy-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:deployagent-0.6.9",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_AGENT_RETRI_NAMESPACE",
                    "value": "eqxiu;"
                  },
                  {
                    "name": "DDOG_SUB_NET_ID",
                    "value": "subnet-ba0hwkov"
                  },
                  {
                    "name": "DDOG_AGENT_NAME",
                    "value": "DeployAgent"
                  },
                  {
                    "name": "DDOG_LOG_OPT",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=100;             --log-opt env=svcname;\""
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=100;             --log-opt env=svcname;\""
                  },
                  {
                    "name": "svcname",
                    "value": "scheduler-deploy-agent"
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 4,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-02T08:45:51Z",
            "lastTransitionTime": "2018-05-02T08:45:51Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    },
    {
      "metadata": {
        "name": "scheduler-monitor-agent",
        "namespace": "scheduler",
        "selfLink": "/apis/apps/v1beta1/namespaces/scheduler/deployments/scheduler-monitor-agent",
        "uid": "05df3ae2-4dd6-11e8-b75e-52540018543c",
        "resourceVersion": "10634454319",
        "generation": 3,
        "creationTimestamp": "2018-05-02T06:57:05Z",
        "labels": {
          "qcloud-app": "scheduler-monitor-agent"
        },
        "annotations": {
          "deployment.changecourse": "Updating",
          "deployment.kubernetes.io/revision": "3"
        }
      },
      "spec": {
        "replicas": 1,
        "selector": {
          "matchLabels": {
            "qcloud-app": "scheduler-monitor-agent"
          }
        },
        "template": {
          "metadata": {
            "creationTimestamp": null,
            "labels": {
              "qcloud-app": "scheduler-monitor-agent",
              "qcloud-redeploy-timestamp": "1525253747"
            }
          },
          "spec": {
            "containers": [
              {
                "name": "scheduler-monitor-agent",
                "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:monitoragent-0.6.9",
                "env": [
                  {
                    "name": "PATH",
                    "value": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
                  },
                  {
                    "name": "DDOG_MONGO_NAME",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_MONGO_PASSWD",
                    "value": "password"
                  },
                  {
                    "name": "DDOG_MONGO_DB",
                    "value": "cloud"
                  },
                  {
                    "name": "DDOG_NSQD_ENDPOINT",
                    "value": "192.168.1.12:4150"
                  },
                  {
                    "name": "DDOG_MONGO_ENDPOINT",
                    "value": "192.168.1.12:27017"
                  },
                  {
                    "name": "DDOG_AGENT_NAME",
                    "value": "MonitorAgent"
                  },
                  {
                    "name": "svcname",
                    "value": "scheduler-monitor-agent"
                  },
                  {
                    "name": "log_opt",
                    "value": "\"--log-opt influx-address=udp://192.168.1.14:8089; --log-opt buf=10;             --log-opt env=svcname;\""
                  },
                  {
                    "name": "LOGCHAIN_DRIVER",
                    "value": "influx"
                  }
                ],
                "resources": {
                  "limits": {
                    "cpu": "500m",
                    "memory": "50Mi"
                  },
                  "requests": {
                    "cpu": "250m",
                    "memory": "20Mi"
                  }
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
            "securityContext": {},
            "imagePullSecrets": [
              {
                "name": "qcloudregistrykey"
              }
            ],
            "schedulerName": "default-scheduler"
          }
        },
        "strategy": {
          "type": "RollingUpdate",
          "rollingUpdate": {
            "maxUnavailable": 0,
            "maxSurge": 1
          }
        },
        "minReadySeconds": 10,
        "revisionHistoryLimit": 5
      },
      "status": {
        "observedGeneration": 3,
        "replicas": 1,
        "updatedReplicas": 1,
        "readyReplicas": 1,
        "availableReplicas": 1,
        "conditions": [
          {
            "type": "Available",
            "status": "True",
            "lastUpdateTime": "2018-05-04T06:44:34Z",
            "lastTransitionTime": "2018-05-04T06:44:34Z",
            "reason": "MinimumReplicasAvailable",
            "message": "Deployment has minimum availability."
          }
        ]
      }
    }
  ]
}
`
	var kd K8sDeploymentInfo

	err := json.Unmarshal([]byte(data), &kd)
	assert.Nil(t, err, "Unmarshal Error")
	assert.Equal(t, "DeploymentList", kd.Kind)
	assert.Equal(t, "11335571561", kd.Metadata.ResourceVersion)
	assert.Equal(t, "caas-deploy-agent", kd.Items[0].Metadata.Labels["qcloud-app"])
	assert.Equal(t, "caas-deploy-agent", kd.Items[0].Spec.Selector.MatchLabels["qcloud-app"])
	assert.Equal(t, "password", kd.Items[0].Spec.Template.Spec.Containers[0].Env[2]["value"])
	assert.Equal(t, "20Mi", kd.Items[0].Spec.Template.Spec.Containers[0].Resources.Requests["memory"])
	assert.Equal(t, "Deployment has minimum availability.", kd.Items[len(kd.Items)-1].Status.Conditions[0].Message)
}
