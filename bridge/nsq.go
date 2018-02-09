package bridge

import (
	"os"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/nsqio/go-nsq"
	"github.com/Sirupsen/logrus"
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

func SendMonitorMsg(msg string) error {
	return makeMsg(_const.SvcMonitorMsg, msg)
}
func makeMsg(topic, msg string) error {
	logrus.WithFields(logrus.Fields{"Topic": topic, "Msg": msg}).Info(ModuleName)
	return nb.producer.Publish(topic, []byte(msg))
}
