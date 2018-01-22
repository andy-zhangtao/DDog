package _const

const (
	SvcIPEmpty              = "The Service IP is Empty!"
	HttpSvcEmpty            = "Svc Field Is Empty!"
	SvcNotFound             = "This Svc Not Found!"
	EnvNotFound             = " Env Not Found!"
	EnvDomainNotFound       = "The " + EnvDomain + EnvNotFound
	EnvEtcdNotFound         = "The " + EnvEtcd + EnvNotFound
	EnvUpStreamNotFound     = "The " + EnvUpStream + EnvNotFound
	EnvMongoNotFound        = "The " + EnvMongo + EnvNotFound
	EnvMongoNameNotFound    = "The " + EnvMongoName + EnvNotFound
	EnvMongoPasswdNotFound  = "The " + EnvMongoPasswd + EnvNotFound
	EnvMongoDBNotFound      = "The " + EnvMongoDB + EnvNotFound
	EnvRegionNotFound       = "The " + EnvRegion + EnvNotFound
	MetaDataDupilcate       = "This Region MetaData Exist!"
	RegionNotFound          = "Region Field Cannot Be Empty!"
	ClusterNotFound         = "ClusterID Field Cannot Be Empty!"
	NamespaceNotFound       = "Namespace Field Cannot Be Empty!"
	SvcConfNotFound         = "Svcname Field Cannot Be Empty!"
	SvcGroupNotFound        = "SvcGroup Field Cannot Be Empty!"
	NameNotFound            = "Name Field Cannot Be Empty!"
	ImageNotFounc           = "Img Field Cannot Be Empty!"
	SVCNoExist              = "This svc doesn't exist! Please check namespace or svc"
	IdxNotFound             = "The Idx Field Cannot Be Empty!"
	IdxVlaueError           = "The Idx Must Be Great Than Zero!"
	SvcConfExist            = "The same name svc configure has exist!"
	ConConfExist            = "The same name container configure has exist!"
	SvcIDNotFound           = "The SvcId Cannot Be Empty!"
	LbProtocolError         = "The protocol type wrong. The value must be 0 or 1"
	LbPortError             = "The Container Port(in_port) Error or LB Port(out_port) Error!"
	AccessTypeError         = "The AccessType Error. The value must in 0,1 and 2"
	IDNotFound              = "The ID Field Cannot Be Empty!"
	SvcHasExist             = "The Svc Has Exist In This Group!"
	ContainerNotFound       = "Cannot find container in this svc"
	NotRollingUP            = "This Svc Doesn't Under Rolling Up Status"
	NotAllInstanceRollingUP = "Not All Instances Complete Rolling Up"
)
