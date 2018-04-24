package event

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/const/v1"
	"net/url"
	"github.com/sirupsen/logrus"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"fmt"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/24.
//获取服务事件数据
type SEvent struct {
	FirstSeen string `json:"firstseen"`
	LastSeen  string `json:"lastseen"`
	Count     int    `json:"count"`
	Level     string `json:"level"`
	ObjType   string `json:"objtype"`
	ObjName   string `json:"objname"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
}

type EventResult struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	CodeDesc string `json:"codedesc"`
	Data struct {
		Events []SEvent `json:"eventList"`
	} `json:"data"`
}

type ServiceEventRequest struct {
	Svcname   string        `json:"svcname"`
	Namespace string        `json:"namespace"`
	ClusterId string        `json:"cluster_id"`
	Pub       public.Public `json:"pub"`
	SecretKey string        `json:"secret_key"`
	sign      string
	Debug     bool
}

const (
	EventAction = "DescribeServiceEvent"
)

func (this *ServiceEventRequest) GetServiceEvent() (events []SEvent, err error) {
	if this.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	pubMap := public.PublicParam(EventAction, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(append(public.PubilcField, []string{"clusterId", "namespace", "serviceName"}...), map[string]string{
		"clusterId":   this.ClusterId,
		"namespace":   this.Namespace,
		"serviceName": this.Svcname,
	}, pubMap)

	signStr := "GET" + v1.QCloudApiEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	reqURL := this.sign + "&Signature=" + url.QueryEscape(sign)

	logrus.WithFields(logrus.Fields{"Request URL": reqURL}).Debug(EventAction)

	resp, err := http.Get(public.API_URL + reqURL)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var es EventResult

	err = json.Unmarshal(data, &es)
	if err != nil {
		err = errors.New(fmt.Sprintf("Query Service [%s] Event Faile! Error [%s]", this.Svcname, err.Error()))
		return nil, err
	}

	if es.Code != 0 {
		err = errors.New(fmt.Sprintf("Query Service [%s] Event Faile! Result Code [%d] Result Message [%s]", this.Svcname, es.Code, es.Message))
		return
	}

	events = es.Data.Events
	return
}
