package check

import (
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"errors"
	"fmt"
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

