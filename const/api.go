package _const

const (
	AddSvcIP         = "/svc/dns/a-record"
	DnsMetaData      = "/dns/metadata"
	GetNodeInfo      = "/cloud/cluster/nodes"
	GetClusterInfo   = "/cloud/cluster/info"
	GetNSInfo        = "/cloud/namespace/info"
	GetSvcMoreInfo   = "/cloud/svc/info/more"
	GetSvcSampleInfo = "/cloud/svc/info/sample"
	DeploySvcConfig  = "/cloud/svc/deploy"
	MetaData         = "/cloud/metadata"
	NewNameSpace     = "/cloud/namespace/create"
	DeleteNameSpace  = "/cloud/namespace/delete"
	CheckNameSpace   = "/cloud/namespace/check"
	NewContainer     = "/cloud/container/create"
	GetContainer     = "/cloud/container/info"
	DeleteContainer  = "/cloud/container/delete"
	UpgradeContainer = "/cloud/contianer/upgrade"
	NewSvcConfig     = "/cloud/svcconf/create"
	GetSvcConfig     = "/cloud/svcconf/info"
	DeleteSvcConfig  = "/cloud/svcconf/delete"
	CheckSvcConfig   = "/cloud/svcconf/check"
	RunService       = "/cloud/svc/run"
	UpgradeService   = "/cloud/svc/upgrade"
	DeleteService    = "/cloud/svc/delete"
	ReinstallService = "/cloud/svc/reinstall"
	AddSvcGroup      = "/cluod/svc/group/add"
	GetSvcGroup      = "/cloud/svc/group/info"
	DeleSvcGroup     = "/cloud/svc/group/delete"
)

type RespMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
