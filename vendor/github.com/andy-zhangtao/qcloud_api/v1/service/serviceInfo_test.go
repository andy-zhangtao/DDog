package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSvcInfoUnmarshal(t *testing.T) {
	var jsons=`{
	"serviceName": "name1",
	"serviceDesc": "des",
	"status": "Waiting",
	"reasonMap": {
		"下载镜像失败": 1
	},
	"reason": "ImagePullBackOff",
	"regionId": 1,
	"desiredReplicas": 1,
	"currentReplicas": 0,
	"lbId": "",
	"lbStatus": "None",
	"createdAt": "2016-12-08 12:44:21",
	"accessType": "LoadBalancer",
	"serviceIp": "100.71.0.60",
	"externalIp": "",
	"namespace": "default",
	"portMappings": [{
		"containerPort": 100,
		"lbPort": 900,
		"nodePort": 32191,
		"protocol": "TCP"
	}],
	"containers": [{
		"containerName": "test",
		"image": "nginx",
		"envs": null,
		"volumeMounts": null,
		"liveProbe": null,
		"readyProbe": null,
		"cpu": 0,
		"memory": 0,
		"command": "",
		"arguments": null
	}],
	"selector": {
		"qcloud-app": "xxx"
	},
	"labels": {
		"qcloud-app": "xxx"
	}
}`

        svcInfo, err := SvcInfoUnmarshal([]byte(jsons))
        assert.Nil(t, err)
        assert.Equal(t, "name1", svcInfo.ServiceName)
        assert.Equal(t, "Waiting", svcInfo.Status)
}
