package _const

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/6.
//保存服务状态常量
const (
	NeedDeploy       = iota
	DeploySuc
	DeployIng
	BGDeployING
	DeployFailed
	RollingUpIng
	DeployStatusSync
	RollBack
	DeployConfirm
)
