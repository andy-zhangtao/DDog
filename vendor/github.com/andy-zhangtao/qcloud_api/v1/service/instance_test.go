package service

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInstanceUnmarshal(t *testing.T) {
	var jsons = `{
			"name": "test1-2249619043-46qpq",
			"status": "Running",
			"reason": "",
			"srcReason": "",
			"ip": "172.16.1.124",
			"restartCount": 0,
			"readyCount": 1,
			"nodeName": "192.168.1.16",
			"nodeIp": "192.168.1.16",
			"createdAt": "2018-01-19 13:39:05",
			"containers": [{
				"name": "eqxiu-base-server",
				"containerId": "cc7a3a78f98d37275b4dc808e64fd63d20f8b278603cbb0873b3cd7fdcc93250",
				"status": "Running",
				"reason": "",
				"image": "ccr.ccs.tencentyun.com\/eqxiu\/eqxiu-base-server:1.1.1"
			}]
		}`

	inst, err := InstanceUnmarshal([]byte(jsons))
	assert.Nil(t, err)
	assert.Equal(t, "Running", inst.Status)
	assert.Equal(t, "test1-2249619043-46qpq", inst.Name)
}
