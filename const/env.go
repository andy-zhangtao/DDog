package _const

const (
	EnvDomain        = "DDOG_DOMAIN"
	EnvEtcd          = "DDOG_ETCD_ENDPOINT"
	EnvUpStream      = "DDOG_UP_STREAM"
	EnvConfPath      = "DDOG_CONF_PATH"
	EnvMongo         = "DDOG_MONGO_ENDPOINT"
	EnvMongoName     = "DDOG_MONGO_NAME"
	EnvMongoPasswd   = "DDOG_MONGO_PASSWD"
	EnvMongoDB       = "DDOG_MONGO_DB"
	EnvRegion        = "DDOG_REGION"
	EnvDefaultNS     = "DDOG_NAME_SPACE"
	EnvGoblin        = "DDOG_GOBLIN_ENDPOINT"
	EnvK8sEndpoint   = "DDOG_K8S_ENDPOINT"
	EnvK8sToken      = "DDOG_K8S_TOKEN"
	EnvDefaultLogOpt = "DDOG_LOG_OPT"
	EnvNsqdEndpoint  = "DDOG_NSQD_ENDPOINT"
	EnvSubNetID      = "DDOG_SUB_NET_ID"
)

var DEBUG = false
var DefaultNameSpace string
var Region string
var K8sEndpoint string
var K8sToken string

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
