package _const

import (
	"os"
	"strconv"
	"log"
)

const (
	EnvDomain      = "DDOG_DOMAIN"
	EnvEtcd        = "DDOG_ETCD_ENDPOINT"
	EnvUpStream    = "DDOG_UP_STREAM"
	EnvConfPath    = "DDOG_CONF_PATH"
	EnvMongo       = "DDOG_MONGO_ENDPOINT"
	EnvMongoName   = "DDOG_MONGO_NAME"
	EnvMongoPasswd = "DDOG_MONGO_PASSWD"
	EnvMongoDB     = "DDOG_MONGO_DB"
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

	if DEBUG {
		log.Println("启动DEBUG模式")
	} else {
		log.Println("关闭DEBUG模式")
	}
}
