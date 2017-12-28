package _const

import (
	"os"
	"strconv"
)

const (
	EnvDomain   = "DDOG_DOMAIN"
	EnvEtcd     = "DDOG_ETCD_ENDPOINT"
	EnvConfPath = "DDOG_CONF_PATH"
)

var DEBUG = false

func init() {
	isDebug := os.Getenv("DDOG_DEBUG")
	debug, err := strconv.ParseBool(isDebug)
	if err != nil {
		DEBUG = false
	} else {
		DEBUG = debug
	}
}
