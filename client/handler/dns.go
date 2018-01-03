package handler

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/etcd"
	"github.com/andy-zhangtao/gogather/strings"
	"net/http"
	"github.com/andy-zhangtao/DDog/server"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/andy-zhangtao/DDog/server/mongo"
)

type Svc struct {
	Svcname      string `json:"svcname"`
	SecretId     string `json:"secret_id"`
	SecretKey    string `json:"secret_key"`
	Region       string `json:"region"`
	Clusterid    string `json:"clusterid"`
	Namespace    string `json:"namespace"`
	Allnamespace string `json:"allnamespace"`
	svcip        string
	Domain       string `json:"domain"`
}

type vsvc struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// WatchDNS 维护服务名称与DNS的对应关系
// 例如:
// 服务名称为: user-manager-domain-com
// 将转换为: /com/domain/manager/user
// 实际在DNS中对应的域名为 user.manager.domain.com
func (this Svc) WatchDNS() error {
	q := service.Svc{
		Pub: public.Public{
			Region:   this.Region,
			SecretId: this.SecretId,
		},
		ClusterId:    this.Clusterid,
		Namespace:    this.Namespace,
		Allnamespace: this.Allnamespace,
		SecretKey:    this.SecretKey,
	}
	q.SetDebug(_const.DEBUG)
	ssmd, err := q.QuerySampleInfo()
	if err != nil {
		return err
	}

	if ssmd.Code != 0 {
		return errors.New(ssmd.Message)
	}

	for _, svc := range ssmd.Data.Services {
		if svc.ServiceIp == "" {
			continue
		}

		err = mongo.SaveService(svc)
		if err != nil {
			return err
		}
		
		s := vsvc{
			Host: svc.ServiceIp,
			Port: 80,
		}

		data, err := json.Marshal(&s)
		if err != nil {
			return err
		}
		if _const.DEBUG {
			log.Printf("[UpDns][%s]\n", string(data))
		}

		err = etcd.Put("/"+strings.ReverseWithSeg(svc.ServiceName, "-", "/"), string(data))
		if err != nil {
			return err
		}
	}

	return nil
}

// ChangeDns 修改指定服务名称的DNS记录
// 当前支持修改A记录,如果修改成功则返回True, 如果找不到指定的服务名称则返回False
// 如果出现其他错误，则返回error
func (this Svc) ChangeDns() (bool, error) {

	q := service.Svc{
		Pub: public.Public{
			Region:   this.Region,
			SecretId: this.SecretId,
		},
		ClusterId:    this.Clusterid,
		Namespace:    this.Namespace,
		Allnamespace: this.Allnamespace,
		SecretKey:    this.SecretKey,
	}

	q.SetDebug(_const.DEBUG)
	ssmd, err := q.QuerySampleInfo()
	if err != nil {
		return false, err
	}

	if ssmd.Code != 0 {
		return false, errors.New(ssmd.Message)
	}

	for _, svc := range ssmd.Data.Services {
		if this.Svcname == svc.ServiceName {
			if svc.ServiceIp == "" {
				return false, errors.New(_const.SvcIPEmpty)
			}
			this.svcip = svc.ServiceIp
			break
		}
	}

	if this.svcip == "" {
		return false, errors.New(_const.SvcIPEmpty)
	}

	s := vsvc{
		Host: this.svcip,
		Port: 80,
	}

	data, err := json.Marshal(&s)
	if err != nil {
		return false, err
	}
	if _const.DEBUG {
		log.Printf("[UpDns][%s]\n", string(data))
	}

	err = etcd.Put("/"+strings.ReverseWithSeg(this.Domain, ".", "/"), string(data))
	return true, err
}

// AddSvcDnsAR 增加指定服务的A记录
func AddSvcDnsAR(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var svc Svc

	err = json.Unmarshal(data, &svc)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	isSuc, err := svc.ChangeDns()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	if !isSuc {
		server.ReturnError(w, errors.New(_const.SvcNotFound))
		return
	}

}
