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
	EnvRegion      = "DDOG_REGION"
	EnvDefaultNS   = "DDOG_NAME_SPACE"
	EnvGoblin      = "DDOG_GOBLIN_ENDPOINT"
)

var DEBUG = false
var DefaultNameSpace string
var Region string

var RegionMap = map[string]string{
	"ap-shanghai": "sh",
}

var ReverseRegionMap = map[string]string{
	"sh": "ap-shanghai",
}

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const OperationSucc = "Operation Succ!"
const OperationFaile = "Operation Faile!"
const DataNotFound = OperationSucc + " But donot find any data!"

func init() {
	isDebug := os.Getenv("DDOG_DEBUG")
	debug, err := strconv.ParseBool(isDebug)
	if err != nil {
		DEBUG = false
	} else {
		DEBUG = debug
	}

	DefaultNameSpace = os.Getenv(EnvDefaultNS)
	if DefaultNameSpace == "" {
		DefaultNameSpace = "default"
	}

	log.Printf("默认命名空间[%s]\n",DefaultNameSpace)
	if DEBUG {
		log.Println("启动DEBUG模式")
	} else {
		log.Println("关闭DEBUG模式")
	}
}
