package bridge

import (
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"os"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/5.

type NsqBridge struct {
	producer *nsq.Producer
}

var nb *NsqBridge

const (
	ModuleName = "Bridge Nsq"
)

func init() {
	nsq_endpoint := os.Getenv(_const.EnvNsqdEndpoint)
	logrus.WithFields(logrus.Fields{
		"Connect NSQ": nsq_endpoint,
	}).Info(ModuleName)

	producer, err := nsq.NewProducer(nsq_endpoint, nsq.NewConfig())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Connect Nsq Error": err,
		}).Panic(ModuleName)
	}

	nb = &NsqBridge{
		producer: producer,
	}

	err = producer.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Ping Nsq Error": err,
		}).Panic(ModuleName)
	}

	logrus.WithFields(logrus.Fields{
		"Connect Nsq Succes": producer.String(),
	}).Info(ModuleName)

}

// SendDestoryMsg 发布销毁服务的消息
func SendDestoryMsg(msg string) error {
	return makeMsg(_const.SvcDestroyMsg, msg)
}

// SendMonitorMsg 发布监控信息消息
func SendMonitorMsg(msg string) error {
	logrus.WithFields(logrus.Fields{"Send Message": msg}).Info(ModuleName)
	if os.Getenv("DDOG_AGENT_SPIDER_NS") != "" {
		return makeMsg(fmt.Sprintf("%s_%s", _const.SvcK8sMonitorMsg, os.Getenv("DDOG_AGENT_SPIDER_NS")), msg)
	}
	return makeMsg(_const.SvcK8sMonitorMsg, msg)
	//return makeMsg(_const.SvcMonitorMsg, msg)
}

func SendDeployMsg(msg string) error {
	return makeMsg(_const.SvcDeployMsg, msg)
}

func makeMsg(topic, msg string) error {
	logrus.WithFields(logrus.Fields{"Topic": topic, "Msg": msg}).Info(ModuleName)
	return nb.producer.Publish(topic, []byte(msg))
}
