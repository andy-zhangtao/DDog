package service

import "strings"

// HealthCheck 服务检查数据
type HealthCheck struct {
	Type         string `json:"type"`
	HealthNum    int    `json:"health_num"`
	UnhealthNum  int    `json:"unhealth_num"`
	IntervalTime int    `json:"interval_time"`
	TimeOut      int    `json:"time_out"`
	DelayTime    int    `json:"delay_time"`
	CheckMethod  string `json:"check_method"`
	Port         int    `json:"port"`
	Protocol     string `json:"protocol"`
	Path         string `json:"path"`
	Cmd          string `json:"cmd"`
}

const (
	ReadyCheck      = "readyCheck"
	LiveCheck       = "liveCheck"
	CheckMethodTCP  = "methodTcp"
	CheckMethodHTTP = "methodHttp"
	CheckMethodCmd  = "methodCmd"
)

// GenerateCheck 生成健康检测参数
func (ht *HealthCheck) GenerateDefaultCheck() {
	if ht.Type == "" {
		ht.Type = ReadyCheck
	}

	if ht.HealthNum == 0 {
		ht.HealthNum = 1
	}

	if ht.UnhealthNum == 0 {
		ht.UnhealthNum = 3
	}

	if ht.IntervalTime == 0 {
		ht.IntervalTime = 10
	}

	if ht.TimeOut == 0 {
		ht.TimeOut = 5
	}

	if ht.DelayTime == 0 {
		ht.DelayTime = 30
	}
}

// GenerateTCPCheck 生成TCP检查
// port 端口
func (ht *HealthCheck) GenerateTCPCheck(port int) {
	ht.GenerateDefaultCheck()
	ht.Protocol = CheckMethodTCP
	ht.Port = port
}

// GenerateHttpCheck 生成Http检查
// port 端口
// url 请求路径
// isHttps 是否为https
func (ht *HealthCheck) GenerateHttpCheck(port int, url string, isHttps bool) {
	ht.GenerateDefaultCheck()
	ht.Protocol = CheckMethodHTTP
	ht.Port = port
	ht.Path = strings.TrimSpace(url)
	if isHttps {
		ht.Protocol = "HTTPS"
	} else {
		ht.Protocol = "HTTP"
	}
}

// GenerateCmdCheck 生成命令检查
// cmd 需要执行的命令
func (ht *HealthCheck) GenerateCmdCheck(cmd string) {
	ht.GenerateDefaultCheck()
	ht.Protocol = CheckMethodCmd
	ht.Cmd = cmd
}
