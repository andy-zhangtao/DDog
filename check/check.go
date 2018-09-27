package check

import (
	"errors"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"os"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/18.
//检测当前是否满足服务启动条件

func CheckMongo() error {
	if os.Getenv(_const.EnvMongo) == "" {
		return errors.New(fmt.Sprintf("%s Empty", _const.EnvMongo))
	}

	if os.Getenv(_const.EnvMongoDB) == "" {
		return errors.New(fmt.Sprintf("%s Empty", _const.EnvMongoDB))
	}

	return nil
}

//Checknamespace 初始化默认命名空间
func CheckNamespace() error {

	if os.Getenv(_const.EnvClusterID) == "" {
		return errors.New(fmt.Sprintf("%s Emtpy", _const.EnvClusterID))
	}

	if os.Getenv(_const.EnvRegion) == "" {
		return errors.New(fmt.Sprintf("%s Emtpy", _const.EnvRegion))
	}

	_const.DefaultNameSpace = os.Getenv(_const.EnvDefaultNS)
	if _const.DefaultNameSpace == "" {
		_const.DefaultNameSpace = "default"
	}

	_const.DefaultDevNameSpace = os.Getenv(_const.EnvDefaultDevNs)
	if _const.DefaultDevNameSpace == "" {
		_const.DefaultDevNameSpace = "eqxiu-dev"
	}

	_const.DefaultTestNameSpace = os.Getenv(_const.EnvDefaultTestNs)
	if _const.DefaultTestNameSpace == "" {
		_const.DefaultTestNameSpace = "eqxiu-test"
	}

	_const.DefaultPreProduceNameSpace = os.Getenv(_const.EnvDefaultPreProduceNs)
	if _const.DefaultPreProduceNameSpace == "" {
		_const.DefaultPreProduceNameSpace = "eqxiu-pre-pro"
	}

	_const.DefaultProduceNameSpace = os.Getenv(_const.EnvDefaultProduceNs)
	if _const.DefaultProduceNameSpace == "" {
		_const.DefaultProduceNameSpace = "eqxiu-pro"
	}
	return nil
}

func CheckNsq() error {
	if os.Getenv(_const.EnvNsqdEndpoint) == "" {
		return errors.New(fmt.Sprintf("%s Empty", _const.EnvNsqdEndpoint))
	}

	return nil
}

// CheckLogOpt 检查默认日志参数
// 用于docker plugin发送日志
func CheckLogOpt() error {
	if os.Getenv(_const.EnvDefaultLogOpt) == "" {
		return errors.New(fmt.Sprintf("%s Empty", _const.EnvDefaultLogOpt))
	}

	return nil
}
