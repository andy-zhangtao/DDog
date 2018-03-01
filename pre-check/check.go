package pre_check

import (
	"strconv"
	"os"
	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/DDog/const"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/6.
//检查当前环境是否满足运行需要

const (
	ModuleName = "Pre-Check"
)

func init() {
	debug, err := strconv.ParseBool(os.Getenv("DDOG_DEBUG"))
	if err != nil {
		_const.DEBUG = false
	} else {
		_const.DEBUG = debug
	}

	_const.DefaultNameSpace = os.Getenv(_const.EnvDefaultNS)
	if _const.DefaultNameSpace == "" {
		_const.DefaultNameSpace = "default"
	}

	_const.K8sEndpoint = os.Getenv(_const.EnvK8sEndpoint)
	if _const.K8sEndpoint == "" {
		logrus.Println("DDOG_K8S_ENDPOINT Empty!")
	}

	_const.K8sToken = os.Getenv(_const.EnvK8sToken)
	if _const.K8sToken == "" {
		logrus.Println("DDOG_K8S_TOKEN Empty! ")
	}

	logrus.Printf("默认命名空间[%s]", _const.DefaultNameSpace)
	if _const.DEBUG {
		logrus.Println("启动DEBUG模式")
	} else {
		logrus.Println("关闭DEBUG模式")
	}

	if os.Getenv(_const.EnvNsqdEndpoint) == "" {
		logrus.WithFields(logrus.Fields{"Env Empty": _const.EnvNsqdEndpoint,}).Panic(ModuleName)
	}

}
