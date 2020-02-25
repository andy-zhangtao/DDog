package _const

const (
	EnvDomain      = "DDOG_DOMAIN"
	EnvEtcd        = "DDOG_ETCD_ENDPOINT"
	EnvUpStream    = "DDOG_UP_STREAM"
	EnvConfPath    = "DDOG_CONF_PATH"
	EnvMongo       = "DDOG_MONGO_ENDPOINT"
	EnvMongoName   = "DDOG_MONGO_NAME"
	EnvMongoPasswd = "DDOG_MONGO_PASSWD"
	EnvMongoDB     = "DDOG_MONGO_DB"
	//EnvRegion 默认机房区域
	EnvRegion = "DDOG_REGION"
	//EnvClusterID 默认集群名称
	EnvClusterID = "DDOG_CLUSTER_ID"
	EnvDefaultNS = "DDOG_NAME_SPACE"
	//EnvDefaultDevNs默认开发环境
	EnvDefaultDevNs = "DDOG_DEV_NAME_SPACE"
	//EnvDefaultTestNs 默认测试环境
	EnvDefaultTestNs = "DDOG_TEST_NAME_SPACE"
	//EnvDefaultPreProduceNs 默认预发布环境
	EnvDefaultPreProduceNs = "DDOG_PRE_PRODUCE_NAME_SPACE"
	//EnvDefaultProduceNs 默认发布环境
	EnvDefaultProduceNs        = "DDOG_PRODUCE_NAME_SPACE"
	EnvGoblin                  = "DDOG_GOBLIN_ENDPOINT"
	EnvK8sEndpoint             = "DDOG_K8S_ENDPOINT"
	EnvK8sToken                = "DDOG_K8S_TOKEN"
	EnvDefaultLogOpt           = "DDOG_LOG_OPT"
	EnvDefaultLogDriver        = "LOGCHAIN_DRIVER"
	EnvMyNsqEndpoint           = "DDOG_MY_NSQD_ENDPOINT"
	EnvNsqdEndpoint            = "DDOG_NSQD_ENDPOINT"
	EnvNsqdEndpointRelease     = "DDOG_NSQD_ENDPOINT_RELEASE"
	EnvSubNetID                = "DDOG_SUB_NET_ID"
	ENV_AGENT_ZIPKIN_ENDPOINT  = "Agent_ZipKin_Endpoint"
	ENV_AGENT_HULK_ENDPOINT    = "Agent_Hulk_Endpoint"
	ENV_DEVEX_GRAPHQL_ENDPOINT = "Devex_Graph_Endpoint"
	//服务部署标志位,用来区分部署集群.
	ENV_DEPLOY_ENV              = "DDOG_DEPLOY_ENV"
	ENV_AGENT_INFLUX_DB         = "DDOG_Influx_DB"
	Env_AGENT_INFLUX_ENDPOINT   = "DDOG_Influx_Endpoint"
	ENV_WATCH_MONITOR_NAMESPACE = "DDOG_WATCH_MON_NAMESPACE" //K8sMonitor 监控命名空间
)

var DEBUG = false
var DefaultNameSpace string
var DefaultDevNameSpace string
var DefaultTestNameSpace string
var DefaultPreProduceNameSpace string
var DefaultProduceNameSpace string
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

const PROENV = "proenv"
const PROENVB = "proenv-b"
const RELEASEENV = "release"
const RELEASEENVB = "release-b"
const RELEASEENVC = "release-c"
const RELEASEENVD = "release-d"
const DEVENV = "devenv"
const TESTENV = "testenv"
const TESTENVB = "testenv-b"
const AUTOENV = "autoenv"
const STUDY = "study"
