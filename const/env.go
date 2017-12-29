package _const

import (
	"os"
	"strconv"
)

const (
	EnvDomain   = "DDOG_DOMAIN"
	EnvEtcd     = "DDOG_ETCD_ENDPOINT"
	EnvUpStream = "DDOG_UP_STREAM"
	EnvConfPath = "DDOG_CONF_PATH"
)

var DEBUG = false
var RegionMap = map[string]string{
	"ap-shanghai": "sh",
}

func init() {
	isDebug := os.Getenv("DDOG_DEBUG")
	debug, err := strconv.ParseBool(isDebug)
	if err != nil {
		DEBUG = false
	} else {
		DEBUG = debug
	}
}
