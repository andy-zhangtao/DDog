package _hulk_client

import (
	"os"
	"github.com/sirupsen/logrus"
)

const (
	HULK_GO_SDK = "hulk-go-sdk"
	ENDPOINT    = "HULK_ENDPOINT"
	NAME        = "HULK_PROJECT_NAME"
	VERSION     = "HULK_PROJECT_VERSION"
)

var _endpoint string
var _name string
var _version string

func init() {
	//检查HULK后端Endpoint是否存在
	if os.Getenv(ENDPOINT) == "" {
		logrus.Error("HULK_ENDPOINT Empty")
		os.Exit(-1)
	}

	//检查HULK 工程名称是否存在
	if os.Getenv(NAME) == "" {
		logrus.Error("HULK_PROJECT_NAME Empty")
		os.Exit(-1)
	}

	//检查HULK 配置版本是否存在
	if os.Getenv(VERSION) == "" {
		logrus.Error("HULK_PROJECT_VERSION Empty")
		os.Exit(-1)
	}

	_endpoint = os.Getenv(ENDPOINT)
	_name = os.Getenv(NAME)
	_version = os.Getenv(VERSION)

	env, err := queryHulk(_name, _version)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(-1)
	}

	padding(env)
}

// 为了支持dep管理, 引入一个空函数
func Run(){}