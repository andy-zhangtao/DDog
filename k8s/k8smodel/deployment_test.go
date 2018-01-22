package k8smodel

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshal(t *testing.T) {
	var jsons = `{
  "kind": "Deployment",
  "apiVersion": "apps/v1beta1",
  "metadata": {
    "name": "test1",
    "namespace": "devenv",
    "selfLink": "/apis/apps/v1beta1/namespaces/devenv/deployments/test1",
    "uid": "95b45c54-fc34-11e7-a645-52540018543c",
    "resourceVersion": "5346066798",
    "generation": 2,
    "creationTimestamp": "2018-01-18T09:47:23Z",
    "labels": {
      "qcloud-app": "test1"
    },
    "annotations": {
      "deployment.changecourse": "Updating",
      "deployment.kubernetes.io/revision": "2"
    }
  },
  "spec": {
    "replicas": 2,
    "selector": {
      "matchLabels": {
        "qcloud-app": "test1"
      }
    },
    "template": {
      "metadata": {
        "creationTimestamp": null,
        "labels": {
          "qcloud-app": "test1",
          "qcloud-redeploy-timestamp": "1516269377"
        }
      },
      "spec": {
        "containers": [
          {
            "name": "eqxiu-base-server",
            "image": "ccr.ccs.tencentyun.com/eqxiu/eqxiu-base-server:1.1.1",
            "resources": {},
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
    "replicas": 2,
    "updatedReplicas": 2,
    "readyReplicas": 2,
    "availableReplicas": 2,
    "conditions": [
      {
        "type": "Available",
        "status": "True",
        "lastUpdateTime": "2018-01-18T09:48:03Z",
        "lastTransitionTime": "2018-01-18T09:48:03Z",
        "reason": "MinimumReplicasAvailable",
        "message": "Deployment has minimum availability."
      }
    ]
  }
}`

	k8d, err := K8dUnmarshal([]byte(jsons))

	assert.Nil(t, err)
	assert.Equal(t, "test1", k8d.Metadata.Name)
	assert.Equal(t, "devenv", k8d.Metadata.Namespace)
	assert.Equal(t, 2, k8d.Spec.Replicas)
	assert.Equal(t, 2, k8d.Status.ReadyReplicas)
	assert.Equal(t, 2, k8d.Status.AvailableReplicas)
}
