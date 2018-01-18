package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateNewSerivce(t *testing.T) {
	svc := Service{}

	svc.ClusterId = "clustid"
	svc.Namespace = "namespace"
	svc.ServiceName = "svcname"
	svc.ServiceDesc = "desc"
	svc.Replicas = 3
	svc.AccessType = "acctype"
	svc.Containers = []Containers{
		Containers{
			ContainerName: "con1",
			Image:         "image1",
			Envs: map[string]string{
				"key1": "value1",
			},
			Command: "ls",
			HealthCheck: []HealthCheck{
				HealthCheck{
					Type:         LiveCheck,
					HealthNum:    2,
					UnhealthNum:  3,
					IntervalTime: 10,
					TimeOut:      10,
					DelayTime:    20,
					CheckMethod:  CheckMethodTCP,
					Port:         3000,
				},
			},
		},
	}
	svc.PortMappings = []PortMappings{
		PortMappings{
			LbPort:        8000,
			Protocol:      "tcp",
			ContainerPort: 9000,
		},
	}

	field, regmap := svc.createSvc()
	assert.Equal(t, 22, len(field), "The length of field error!")

	assert.Equal(t, "clustid", regmap["clusterId"], "clusterId:clustid Error")
	assert.Equal(t, "namespace", regmap["namespace"], "namespace:namespace Error")
	assert.Equal(t, "acctype", regmap["accessType"], "accessType:acctype Error")
	assert.Equal(t, "tcp", regmap["portMappings.0.protocol"], "portMappings.0.protocol:tcp Error")
	assert.Equal(t, "svcname", regmap["serviceName"], "serviceName:svcname  Error")
	assert.Equal(t, "key1", regmap["containers.0.envs.0.name"], "containers.0.envs.0.name:key1 Error")
	assert.Equal(t, "value1", regmap["containers.0.envs.0.value"], "containers.0.envs.0.value:value1 Error")
	assert.Equal(t, "liveCheck", regmap["containers.0.healthCheck.0.type"], "containers.0.healthCheck.0.type:liveCheck Error")
	assert.Equal(t, "desc", regmap["serviceDesc"], "serviceDesc:desc error!")
	assert.Equal(t, "con1", regmap["containers.0.containerName"], "containers.0.containerName:con1 Error")
	assert.Equal(t, "image1", regmap["containers.0.image"], "containers.0.image:image1 Error")
	assert.Equal(t, "methodTcp", regmap["containers.0.healthCheck.0.checkMethod"], "containers.0.healthCheck.0.checkMethod:methodTcp Error")
	assert.Equal(t, "ls", regmap["containers.0.command"], "containers.0.command:ls Error")
	assert.Equal(t, "2", regmap["containers.0.healthCheck.0.healthNum"], "containers.0.healthCheck.0.healthNum:2 Error")
	assert.Equal(t, "3", regmap["replicas"], "replicas:3 Error")
	assert.Equal(t, "3", regmap["containers.0.healthCheck.0.unhealthNum"], "containers.0.healthCheck.0.unhealthNum:3 Error")
	assert.Equal(t, "10", regmap["containers.0.healthCheck.0.timeOut"], "containers.0.healthCheck.0.timeOut:10 Error")
	assert.Equal(t, "10", regmap["containers.0.healthCheck.0.intervalTime"], "containers.0.healthCheck.0.intervalTime:10 Error")
	assert.Equal(t, "20", regmap["containers.0.healthCheck.0.delayTime"], "containers.0.healthCheck.0.delayTime:20  Error")
	assert.Equal(t, "3000", regmap["containers.0.healthCheck.0.port"], "containers.0.healthCheck.0.port:3000 Error")
	assert.Equal(t, "8000", regmap["portMappings.0.lbPort"], "portMappings.0.lbPort:8000 Error")
	assert.Equal(t, "9000", regmap["portMappings.0.containerPort"], "portMappings.0.containerPort:9000  Error")
}
