package qcloud

import (
	"net/http"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"errors"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/gogather/zsort"
	"fmt"
	"github.com/andy-zhangtao/DDog/k8s"
	"github.com/andy-zhangtao/DDog/k8s/k8smodel"
	gt "github.com/andy-zhangtao/gogather/time"
	"github.com/andy-zhangtao/DDog/model/container"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/andy-zhangtao/DDog/bridge"
)

var globalChan chan int
var globalMap map[string]chan int

const (
	ModuleName = "QCloud Service"
)

func getChan(gc string) (chan int) {
	/*需要判断是否有不存在的chan，否则有可能会产生阻塞*/
	if c, ok := globalMap[gc]; ok {
		return c
	} else {
		tc := make(chan int)
		go func() {
			<-tc
		}()
		return tc
	}
}

func setChan(gc string) {
	if globalMap == nil {
		globalMap = make(map[string]chan int)
	}

	globalMap[gc] = make(chan int)
}

// closeChan 关闭并且删除
func closeChan(gc string) {
	if c, ok := globalMap[gc]; ok {
		close(c)
		delete(globalMap, gc)
	}

}
func GetSampleSVCInfo(w http.ResponseWriter, r *http.Request) {

	var id string
	var scf *svcconf.SvcConf
	var err error
	var nsme string
	name := r.URL.Query().Get("svcname")
	if name != "" {
		//	如果上传服务名称，则直接重新部署此服务
		nsme = r.URL.Query().Get("namespace")
		if nsme == "" {
			nsme = _const.DefaultNameSpace
			if nsme == "" {
				tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
				return
			}
		}
	} else {
		id = r.URL.Query().Get("id")
		if id == "" {
			tool.ReturnError(w, errors.New(_const.IDNotFound))
			return
		}
		scf, err = svcconf.GetSvcConfByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}
	name = strings.TrimSpace(name)
	nsme = strings.TrimSpace(nsme)

	scf, err = svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Error": err, "svcname": name, "namespace": nsme, "Operation": "GetSampleSVCInfo"}).Info(ModuleName)
		tool.ReturnError(w, err)
		return
	}

	type SvcStatus struct {
		Name   string           `json:"name"`
		Status string           `json:"status"`
		ULb    map[string][]int `json:"ulb"`
		MLb    map[string][]int `json:"lb"`
		//LbPort   []int  `json:"lb_port"`
		Replicas int    `json:"replicas"`
		Msg      string `json:"msg"`
	}

	sm := make(map[string][]int)
	for _, s := range scf.SvcNameBak {
		sm[s.IP] = s.Port
	}

	ss := SvcStatus{
		Name: scf.Name,
		MLb: map[string][]int{
			scf.LbConfig.IP: scf.LbConfig.Port,
		},
		Replicas: len(scf.Instance),
		Msg:      scf.Msg,
	}

	switch scf.Deploy {
	case 0:
		ss.Status = "ready for deploy"
	case 1:
		ss.Status = "normal"
		ss.ULb = sm
	case 2:
		ss.Status = "updating"
	case 3:
		ss.Status = "rolling fully complete"
		ss.ULb = sm
	case 4:
		ss.Status = "failed"
	case 5:
		ss.Status = "rolling partially complete"
		ss.ULb = sm
	case 6:
		ss.Status = "rolling"
	}

	data, err := json.Marshal(&ss)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Marshal Error": err, "origin data": ss, "Operation": "GetSampleSVCInfo"}).Info(ModuleName)
		tool.ReturnError(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// RunService 在K8s集群中创建服务
// svcname 服务配置名称, 此名称应该在创建服务之前首先创建。
// namespace 命名空间名称, 如果为空则为默认值
// upgrade 是否为升级操作. 默认为false。
// 当在创建服务时，会使用以下默认参数
// 1. 默认启用健康检查和准备就绪检查。
// 2. 上述两种检查使用TCP端口检查方式
// 3. 均针对容器对外暴露的端口进行检查，如果镜像构建未对外暴露端口，则不会对此镜像启用检查
// 4. 延时30秒后启动检查
// 5. 连续三次，间隔10秒，健康均失败则检查失败
// 6. 每次检查超时时间为5秒
// 7. 每个服务默认存在2个实例
func RunService(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("svcname")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		if nsme == "" {
			nsme = _const.DefaultNameSpace
		}
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	isUpgrade := false

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	replicas := 0
	if r.URL.Query().Get("replicas") == "" {
		replicas = cf.Replicas
	} else {
		replicas, err = strconv.Atoi(r.URL.Query().Get("replicas"))
		if err != nil {
			replicas = cf.Replicas
		}
	}

	logrus.WithFields(logrus.Fields{"svc_conf": cf,}).Info("RunService")

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	sn := cf.Name + "-" + gt.GetTimeStamp(10)
	isUpgrade, err = strconv.ParseBool(r.URL.Query().Get("upgrade"))
	if err != nil {
		isUpgrade = false
	}

	if isUpgrade {
		// 服务直接升级,不需要通过蓝绿发布
		if cf.SvcName != "" {
			data, _ := json.Marshal(_const.DestoryMsg{
				Svcname:   cf.SvcName,
				Namespace: cf.Namespace,
			})
			bridge.SendDestoryMsg(string(data))
		}
		cf.SvcName = sn
	} else {
		//蓝绿发布
		if cf.SvcName != "" {
			/*当前存在正式服务，则此操作应该是升级操作*/
			isUpgrade = true
			if len(cf.SvcNameBak) == 0 {
				cf.SvcNameBak = map[string]svcconf.LoadBlance{
					sn: svcconf.LoadBlance{},
				}
			} else {
				cf.SvcNameBak[sn] = svcconf.LoadBlance{}
			}
		} else {
			cf.SvcName = sn
		}
	}

	//if cf.SvcName != "" {
	//	go func() {
	//		q := service.Service{
	//			Pub: public.Public{
	//				SecretId: md.Sid,
	//				Region:   md.Region,
	//			},
	//			ClusterId:   md.ClusterID,
	//			ServiceName: cf.SvcName,
	//			Namespace:   cf.Namespace,
	//			SecretKey:   md.Skey,
	//		}
	//
	//		resp, err := q.DeleteService()
	//		if err != nil {
	//			logrus.WithFields(logrus.Fields{
	//				"error": err.Error(),
	//			}).Error("Delete Service Error")
	//		}
	//
	//		if resp.Code != 0 {
	//			logrus.WithFields(logrus.Fields{
	//				"resp_code":   resp.Code,
	//				"resp_reason": resp.Message,
	//			}).Warn("Delete Service Failed")
	//		}
	//	}()
	//}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: sn,
		ServiceDesc: cf.Desc,
		Replicas:    replicas,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)
	if len(cf.Netconf) > 0 {
		var pm []service.PortMappings
		for _, n := range cf.Netconf {
			p := service.PortMappings{}
			switch n.Protocol {
			case 0:
				p.Protocol = "TCP"
			case 1:
				p.Protocol = "UDP"
			}
			p.ContainerPort = n.InPort
			p.LbPort = n.OutPort
			pm = append(pm, p)
		}
		q.PortMappings = pm
		switch cf.Netconf[0].AccessType {
		case 0:
			q.AccessType = "ClusterIP"
		case 1:
			q.AccessType = "LoadBalancer"
		case 2:
			q.AccessType = "SvcLBTypeInner"
		}
	} else {
		q.AccessType = "None"
	}

	var cons []service.Containers

	containers, err := container.GetAllContainersBySvc(cf.Name, cf.Namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if len(containers) == 0 {
		tool.ReturnError(w, errors.New(fmt.Sprintf("[Find Container Error][%s]svc[%s]namespace[%s]", _const.ContainerNotFound, cf.Name, cf.Namespace)))
		return
	}

	for _, cn := range containers {
		var cnns container.Container
		data, err := bson.Marshal(cn)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		err = bson.Unmarshal(data, &cnns)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		var hk []service.HealthCheck

		/*如果对外提供网络端口，则使用TCP的健康检测，否则使用CMD方式的健康检测*/
		if len(cn.Net) > 0 {
			for _, n := range cn.Net {
				shk := service.HealthCheck{
					Type:        service.LiveCheck,
					UnhealthNum: 5,
					DelayTime:   30,
					CheckMethod: service.CheckMethodTCP,
				}
				shk.GenerateTCPCheck(n.InPort)

				hk = append(hk, shk)
				shk.Type = service.ReadyCheck
				hk = append(hk, shk)
			}
		} else {
			/*当前默认使用ps -ef |grep svcname来作为*/
			shk := service.HealthCheck{
				Type:        service.LiveCheck,
				UnhealthNum: 5,
				DelayTime:   30,
				CheckMethod: service.CheckMethodCmd,
			}
			cmd := fmt.Sprintf("/bin/sh -c \"ps -ef | grep %s |grep -v grep\"", cf.Name)
			logrus.WithFields(logrus.Fields{"Cmd Check": cmd, "Operation": "RunService"}).Info(ModuleName)
			shk.GenerateCmdCheck(cmd)
			hk = append(hk, shk)
			shk.Type = service.ReadyCheck
			hk = append(hk, shk)
		}

		cons = append(cons, service.Containers{
			ContainerName: cnns.Name,
			Image:         cnns.Img,
			HealthCheck:   hk,
			Envs: cnns.Env,
		})
	}

	q.Containers = cons

	logrus.WithFields(logrus.Fields{"QCloud Request": q, "Deploy Type": isUpgrade}).Info(ModuleName)

	resp, err := q.CreateNewSerivce()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if resp.Code != 0 {
		cf.Deploy = 4
		cf.Msg = resp.Message
		svcconf.UpdateSvcConf(cf)
		tool.ReturnError(w, errors.New(resp.Message))
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var plugin = func(scf *svcconf.SvcConf, data interface{}) {
		c, _ := data.(map[string]interface{})

		tcf, _ := c["tempSvc"].(*svcconf.SvcConf)
		isupgrade, _ := c["upgrade"].(bool)
		sn, _ := c["rollsvc"].(string) /*本次升级的服务名*/

		if isupgrade && scf.Deploy == 1 {
			scf.SvcNameBak[sn] = scf.LbConfig
			scf.LbConfig = tcf.LbConfig
			//scf.Instance = tcf.Instance
			scf.Netconf = tcf.Netconf
			scf.Deploy = 6
			svcconf.UpdateSvcConf(scf)
		}
		getChan(cf.SvcName) <- 1
	}

	tcf := new(svcconf.SvcConf)
	cf.Copy(tcf)
	go asyncQueryServiceStatus(sn, nsme, q, cf, map[string]interface{}{
		"upgrade": isUpgrade,
		"tempSvc": tcf,
		"rollsvc": sn,
	}, plugin)

	cf.Deploy = 2
	svcconf.UpdateSvcConf(cf)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("EQXC-Run-Svc", "200")
	w.Write(data)
}

func DeleteService(w http.ResponseWriter, r *http.Request) {

	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	var cf *svcconf.SvcConf
	var err error
	id := r.URL.Query().Get("id")
	if id == "" {
		name := r.URL.Query().Get("svcname")
		if name == "" {
			tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
			return
		}

		nsme := r.URL.Query().Get("namespace")
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}

		cf, err = svcconf.GetSvcConfByName(name, nsme)
	} else {
		cf, err = svcconf.GetSvcConfByID(id)
	}

	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   clusterid,
		ServiceName: cf.Name,
		Namespace:   cf.Namespace,
		SecretKey:   md.Skey,
	}

	resp, err := q.DeleteService()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("EQXC-Run-Svc", "200")
	w.Write(data)
}

func ReinstallService(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	var id string
	var cf *svcconf.SvcConf
	var err error
	var nsme string
	name := r.URL.Query().Get("svcname")
	if name != "" {
		//	如果上传服务名称，则直接重新部署此服务
		nsme = r.URL.Query().Get("namespace")
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	} else {
		id = r.URL.Query().Get("id")
		if id == "" {
			tool.ReturnError(w, errors.New(_const.IDNotFound))
			return
		}
		cf, err = svcconf.GetSvcConfByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}

	if name == "" {
		name = cf.Name
	}
	if nsme == "" {
		nsme = cf.Namespace
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   clusterid,
		ServiceName: name,
		Namespace:   nsme,
		SecretKey:   md.Skey,
	}

	resp, err := q.RedeployService()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func DeployService(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("svcname")
	if name == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		nsme = _const.DefaultNameSpace
		if nsme == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	cf, err := svcconf.GetSvcConfByName(name, nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if cf == nil {
		tool.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	//md, err := metadata.GetMetaDataByRegion("")
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//q := service.Svc{
	//	Pub: public.Public{
	//		SecretId: md.Sid,
	//		Region:   md.Region,
	//	},
	//	ClusterId: md.ClusterID,
	//	Namespace: cf.Namespace,
	//	SecretKey: md.Skey,
	//}
	//q.SetDebug(true)
	//resp, err := q.QuerySampleInfo()
	//if err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//isUpgrade := false
	//for _, r := range resp.Data.Services {
	//	if _const.DEBUG {
	//		log.Printf("[DeployService] Find Svc Dist:[%s] Current:[%s]\n", cf.Name, r.ServiceName)
	//	}
	//	if strings.Compare(r.ServiceName, cf.Name) == 0 {
	//		isUpgrade = true
	//		break
	//	}
	//}
	//
	oldPath := r.URL.RawQuery + "&namespace=" + cf.Namespace
	//
	//if isUpgrade {
	//	// 进行蓝绿发布
	//	r.URL.RawQuery = oldPath + "&upgrade=true"
	//} else {
	//	// 同时发布
	//	r.URL.RawQuery = oldPath + "&upgrade=false"
	//}

	r.URL.RawQuery = oldPath + "&upgrade=true"
	logrus.WithFields(logrus.Fields{
		"url":     r.URL.String(),
		"oldPath": oldPath,
	}).Info("DeployService")

	RunService(w, r)

}

func RollingUpService(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svcname")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
		if namespace == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	scp := 0.5
	scope := r.URL.Query().Get("percent")
	if scope != "" {
		scp, err := strconv.Atoi(scope)
		if err != nil {
			scp = scp / 100
		}
	}

	scf, err := svcconf.GetSvcConfByName(svc, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	rollCons, leftName, left := scf.CountInstances(scp)
	if left <= 0 && len(rollCons) == 0 {
		scf.Deploy = 3
		svcconf.UpdateSvcConf(scf)
		return
	}

	if len(rollCons) == 0 {
		tool.ReturnError(w, errors.New("No Instance Will Be Update!"))
		return
	}

	ins := scf.Instance
	for _, n := range rollCons {
		for i, is := range ins {
			if is.Name == n {
				ins[i].Status = 4
			}
		}
	}

	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		Namespace:   scf.Namespace,
		ServiceName: scf.SvcName,
		SecretKey:   md.Skey,
		Instance:    rollCons,
		Replicas:    scf.Replicas,
	}

	q.SetDebug(true)

	resp, err := q.DestoryInstance()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	data, err := json.Marshal(&resp)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var updateInstance = func(scf *svcconf.SvcConf, rollCons interface{}) {
		/* 如果当前Instance实例中有LefeName中的记录，则表明是本次没有参与升级的 */
		c, _ := rollCons.([]string)
		ins := scf.Instance

		for i, s := range ins {
			ins[i].Status = 4
			for _, n := range c {
				logrus.WithFields(logrus.Fields{"New Name": s.Name, "Old Name": n}).Info(ModuleName)
				if s.Name == n {
					ins[i].Status = 1
				}
			}
		}

		scf.Deploy = 5
		scf.Instance = ins
	}

	go asyncQueryServiceStatus(scf.Name, scf.Namespace, q, scf, leftName, updateInstance)

	scf.Instance = ins
	scf.Deploy = 6
	svcconf.UpdateSvcConf(scf)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// RollBackService 服务回滚
// 只有全部回滚，没有部分回滚
func RollBackService(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svcname")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
		if namespace == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	scf, err := svcconf.GetSvcConfByName(svc, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	/*首先恢复原服务*/
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: scf.SvcName,
		Namespace:   scf.Namespace,
		ScaleTo:     scf.Replicas,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)
	_, err = q.ModeifyInstance()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var plugin = func(scf *svcconf.SvcConf, data interface{}) {
		for {
			/*等待部署服务结束，统一更新数据*/
			tscf, err := svcconf.GetSvcConfByID(scf.Id.Hex())
			if err != nil {
				scf.Msg = err.Error()
				break
			}

			//scf.SvcNameBak[sn] = scf.LbConfig
			//scf.LbConfig = tcf.LbConfig
			////scf.Instance = tcf.Instance
			//scf.Netconf = tcf.Netconf
			if tscf.Deploy != 6 {
				time.Sleep(10 * time.Second)
			} else {
				tscf.Instance = scf.Instance
				scf = tscf
				scf.Deploy = 1
				break
			}
		}
	}

	go func() {
		/*延时5秒再开始查询，避免出现状态不准确的情况*/
		time.Sleep(15 * time.Second)
		asyncQueryServiceStatus(scf.SvcName, scf.Namespace, q, scf, nil, plugin)
	}()

	scf.Deploy = 7
	svcconf.UpdateSvcConf(scf)

	/*删除升级服务*/
	for key, _ := range scf.SvcNameBak {
		q := service.Service{
			Pub: public.Public{
				SecretId: md.Sid,
				Region:   md.Region,
			},
			ClusterId:   md.ClusterID,
			ServiceName: key,
			Namespace:   scf.Namespace,
			SecretKey:   md.Skey,
		}

		q.SetDebug(true)
		_, err = q.DeleteService()
		if err != nil {
			logrus.WithFields(logrus.Fields{"Delete Upgrade Service Error": err}).Error(ModuleName)
		}
	}

	scf.SvcNameBak = nil
	scf.Deploy = 6
	svcconf.UpdateSvcConf(scf)
}

// RollingUpServiceWithSvc 以Service为单位进行滚动升级
// 当触发升级操作时, 以当前服务配置为模板，创建另外一个新的服务
func RollingUpServiceWithSvc(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svcname")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
		if namespace == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	scp := 0.5
	scope := r.URL.Query().Get("percent")
	if scope != "" {
		sc, err := strconv.Atoi(scope)
		if err == nil {
			scp = float64(sc) / float64(100)
		}
	}
	logrus.WithFields(logrus.Fields{"Need Rolling Up": scp}).Info(ModuleName)
	scf, err := svcconf.GetSvcConfByName(svc, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	rollCons, _, left := scf.CountInstances(scp)
	if left <= 0 && len(rollCons) == 0 {
		scf.Deploy = 3
		svcconf.UpdateSvcConf(scf)
		return
	}

	if len(rollCons) == 0 {
		tool.ReturnError(w, errors.New("No Instance Will Be Update!"))
		return
	}

	/*创建升级服务*/
	r.URL.RawQuery += fmt.Sprintf("&namespace=%s&replicas=%d", scf.Namespace, len(rollCons))

	logrus.WithFields(logrus.Fields{"Request URL": r.URL.String()}).Info(ModuleName)

	setChan(scf.SvcName)
	RunService(w, r)

	<-getChan(scf.SvcName)
	closeChan(scf.SvcName)

	/*开始缩容*/
	/*在RunService中已经修改了部分状态,所以更新当前服务状态*/
	scf, _ = svcconf.GetSvcConfByName(svc, namespace)
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: scf.SvcName,
		Namespace:   scf.Namespace,
		ScaleTo:     left,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)
	_, err = q.ModeifyInstance()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	/*当前ScaleTo的值就是预期的实例数,如果不赋值，在异步采集状态时，就会直接退出*/
	q.Replicas = left
	var plugin = func(scf *svcconf.SvcConf, data interface{}) {
		tscf, err := svcconf.GetSvcConfByID(scf.Id.Hex())
		if err != nil {
			scf.Msg = err.Error()

		}
		tscf.Instance = scf.Instance
		tscf.Copy(scf)
		scf.Deploy = 5
	}

	go func() {
		/*延时5秒再开始查询，避免出现状态不准确的情况*/
		time.Sleep(5 * time.Second)
		asyncQueryServiceStatus(scf.SvcName, scf.Namespace, q, scf, nil, plugin)
	}()

	svcconf.UpdateSvcConf(scf)
	return
}

// ConfirmRollService 确认升级完成. 只有当前状态为滚动升级中，并且所有实例状态都是升级成功的情况下才可调用此API
func ConfirmRollService(w http.ResponseWriter, r *http.Request) {
	svc := r.URL.Query().Get("svcname")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.SvcConfNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		namespace = _const.DefaultNameSpace
		if namespace == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}
	}

	scf, err := svcconf.GetSvcConfByName(svc, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	/*只有当前可用实例为0，并且状态为滚动部署完成时才可以使用确认功能*/
	if scf.Deploy != 5 || len(scf.Instance) != 0 {
		tool.ReturnError(w, errors.New(_const.NotRollingUP))
		return
	}

	/*随机挑选一个服务，然后将实例数扩展为预定值*/
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	sn := ""
	for k, _ := range scf.SvcNameBak {
		sn = k
		break
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		ServiceName: sn,
		Namespace:   scf.Namespace,
		ScaleTo:     scf.Replicas,
		SecretKey:   md.Skey,
	}

	q.SetDebug(true)
	q.Replicas = scf.Replicas
	_, err = q.ModeifyInstance()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	go func(esn string, scf *svcconf.SvcConf) {
		// esn 是正在确认的服务,不需要删除
		for s, _ := range scf.SvcNameBak {
			if s != esn {
				q := service.Service{
					Pub: public.Public{
						SecretId: md.Sid,
						Region:   md.Region,
					},
					ClusterId:   md.ClusterID,
					ServiceName: s,
					Namespace:   scf.Namespace,
					ScaleTo:     scf.Replicas,
					SecretKey:   md.Skey,
				}

				q.SetDebug(true)
				q.DeleteService()
			}
		}

		q := service.Service{
			Pub: public.Public{
				SecretId: md.Sid,
				Region:   md.Region,
			},
			ClusterId:   md.ClusterID,
			ServiceName: scf.SvcName,
			Namespace:   scf.Namespace,
			ScaleTo:     scf.Replicas,
			SecretKey:   md.Skey,
		}
		q.SetDebug(true)
		q.DeleteService()
	}(sn, scf)
	scf.Deploy = 8
	svcconf.UpdateSvcConf(scf)
	asyncQueryServiceStatus(sn, scf.Namespace, q, scf, nil, nil)
	scf.Deploy = 1
	scf.SvcName = sn
	scf.SvcNameBak = nil
	svcconf.UpdateSvcConf(scf)
	tool.ReturnResp(w, []byte("Confirm Success"))
	return
}

func RunSvcGroup(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svcConfGroup := r.URL.Query().Get("svcgroup")
	if svcConfGroup == "" {
		tool.ReturnError(w, errors.New(_const.SvcGroupNotFound))
		return
	}

	logrus.WithFields(logrus.Fields{"clusterid": clusterid, "namespace": namespace, "svcgroup": svcConfGroup}).Info(ModuleName)

	sg, err := mongo.GetSvcConfGroupByName(svcConfGroup, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	svcg, err := svcconf.Unmarshal(sg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{"svcg": svcg}).Info(ModuleName)

	svcPair := zsort.SortByValue(svcg.SvcGroup)
	rawQuery := r.URL.RawQuery
	nd := strings.Index(rawQuery, "&svcname=")
	if nd > 0 {
		//clear query path
		rawQuery = rawQuery[:nd]
	}

	for i := len(svcPair) - 1; i >= 0; i -- {

		r.URL.RawQuery = rawQuery + "&svcname=" + svcPair[i].Key

		logrus.WithFields(logrus.Fields{"Deploy Svcname": svcPair[i].Key, "Header": r.URL.Query()}).Info(ModuleName)

		w.Header().Del("EQXC-Run-Svc")
		DeployService(w, r)
		logrus.WithFields(logrus.Fields{"Deploy SvcName": svcPair[i].Key, "Response": w.Header()}).Info(ModuleName)

		if w.Header().Get("EQXC-Run-Svc") != "200" {
			return
		}

	}
}

func ReinstallSvcGroup(w http.ResponseWriter, r *http.Request) {
	UninstallSvcGroup(w, r)
	RunSvcGroup(w, r)
}

func UninstallSvcGroup(w http.ResponseWriter, r *http.Request) {
	clusterid := r.URL.Query().Get("clusterid")
	if clusterid == "" {
		tool.ReturnError(w, errors.New(_const.ClusterNotFound))
		return
	}

	namespace := r.URL.Query().Get("namespace")
	if namespace == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svcConfGroup := r.URL.Query().Get("svcgroup")
	if svcConfGroup == "" {
		tool.ReturnError(w, errors.New(_const.SvcGroupNotFound))
		return
	}

	logrus.Printf("[UninstallSvcGroup]clusterid:[%s]namespace:[%s]svcgroup:[%s]\n", clusterid, namespace, svcConfGroup)

	sg, err := mongo.GetSvcConfGroupByName(svcConfGroup, namespace)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	svcg, err := svcconf.Unmarshal(sg)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	logrus.Printf("[UninstallSvcGroup]svcg:[%v]\n", svcg)

	svcPair := zsort.SortByValue(svcg.SvcGroup)
	rawQuery := r.URL.RawQuery

	nd := strings.Index(rawQuery, "&svcname=")
	if nd > 0 {
		//clear query path
		rawQuery = rawQuery[:nd]
	}

	for i := len(svcPair) - 1; i >= 0; i -- {
		r.URL.RawQuery = rawQuery + "&svcname=" + svcPair[i].Key

		logrus.Printf("[UninstallSvcGroup]Delete svcname :[%s] All header:[%v] \n", svcPair[i].Key, r.URL.Query())

		w.Header().Del("EQXC-Run-Svc")
		DeleteService(w, r)
		logrus.Printf("[UninstallSvcGroup]Delete svcname :[%s] Response:[%v] \n", svcPair[i].Key, w.Header())

		if w.Header().Get("EQXC-Run-Svc") != "200" {
			return
		}
	}
}

// queryInstance 查询指定服务的实例状态
func queryInstance(svc, namespace string) (instances []service.Instance, err error) {
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		return nil, err
	}

	q := service.Service{
		Pub: public.Public{
			SecretId: md.Sid,
			Region:   md.Region,
		},
		ClusterId:   md.ClusterID,
		Namespace:   namespace,
		SecretKey:   md.Skey,
		ServiceName: svc,
	}

	q.SetDebug(true)
	resp, err := q.QueryInstance()
	if err != nil {
		return
	}

	logrus.Printf("[queryInstance] Query Instance [%v] \n", resp)

	if resp.Code != 0 {
		err = errors.New(resp.Message)
		return
	}
	instances = resp.Data.Instance
	return
}

func queryInstanceUseK8s(svc, namespace string) (instances []service.Instance, err error) {
	k := k8s.K8sMetaData{
		Endpoint:  _const.K8sEndpoint,
		Namespace: namespace,
		Svcname:   svc,
		Version:   "1.7",
		Token:     _const.K8sToken,
	}

	k8p, err := k.GetPodsInNamespace()
	if err != nil {
		return nil, err
	}

	var k8sService []k8smodel.K8sPods_items
	for _, k := range k8p.Items {
		if k.Metadata.Labels.Qcloud_app == svc {
			k8sService = append(k8sService, k)
		}
	}

	for _, k := range k8sService {
		instances = append(instances, service.Instance{
			Name:   k.Metadata.Name,
			Status: k.Status.Phase,
		})
	}

	logrus.Printf("[queryInstanceUseK8s] K8s [%v]  instances [%v] \n", k8p, instances)
	return

}

// 异步查询服务状态
// 当查询结束时会更新服务状态
// plugin 是每次查询结束时的回调函数
func asyncQueryServiceStatus(svc, namespace string, q service.Service, scf *svcconf.SvcConf, para interface{}, plugin func(conf *svcconf.SvcConf, param interface{})) {
	logrus.WithFields(logrus.Fields{"ServiceConf": scf}).Info(ModuleName)

	errIdx := 0
	jinx := 0 /*意想不到的崩溃次数,也作为失败的一种判断指标*/
	scf.Msg = ""
	//scf.Deploy = 0
	// 轮询当前服务的运行状态
	for {
		resp, err := q.QuerySvcInfo()
		if err != nil {
			logrus.WithFields(logrus.Fields{"QueryViaQCloud Error": err}).Error(ModuleName)
			errIdx ++
		}

		if errIdx == 3 {
			scf.Deploy = 4
			break
		}

		/*需要判断K8s是否出错,以免出现无效查询*/
		if resp.Code != 0 {
			errIdx ++
			scf.Msg = resp.Message
			break
		}

		if jinx == 10 && strings.ToLower(resp.Data.ServiceInfo.Status) != "normal" {
			scf.Deploy = 4
			for key, _ := range resp.Data.ServiceInfo.ReasonMap {
				scf.Msg += key + ";"
			}
			break
		}

		if strings.ToLower(resp.Data.ServiceInfo.Status) != "normal" {
			for key, _ := range resp.Data.ServiceInfo.ReasonMap {
				if key == "容器进程崩溃" {
					jinx ++
				}
			}
		}

		logrus.WithFields(logrus.Fields{"Service": resp.Data.ServiceInfo.ServiceName, "Status": resp.Data.ServiceInfo.Status}).Info(ModuleName)

		if strings.ToLower(resp.Data.ServiceInfo.Status) == "normal" {
			//先解析负载数据
			lb := svcconf.LoadBlance{
				IP: resp.Data.ServiceInfo.ServiceIp,
			}

			var port []int
			for _, c := range resp.Data.ServiceInfo.PortMappings {
				port = append(port, c.LbPort)
			}
			lb.Port = port
			scf.LbConfig = lb

			//succIns 成功启动的实例数
			//正常情况下, 成功启动的实例数应该等于目标实例数(distIns)
			//succIns := 0

			for {
				instances, err := queryInstance(svc, namespace)
				if err != nil {
					instances, err = queryInstanceUseK8s(svc, namespace)
					if err != nil {
						logrus.WithFields(logrus.Fields{"QueryInstance Error": err, "svc": svc, "namespace": namespace}).Error(ModuleName)
						scf.Msg = err.Error()
						errIdx ++
					}
				}

				logrus.WithFields(logrus.Fields{"QueryInstance": instances}).Info(ModuleName)

				var sic []svcconf.SvcInstance
				for _, inc := range instances {
					statusFlag := 0
					switch strings.ToLower(inc.Status) {
					case "running":
						statusFlag = 1
						sic = append(sic, svcconf.SvcInstance{
							Name:   inc.Name,
							Msg:    inc.Reason,
							Status: statusFlag,
						})
					case "waiting":
						fallthrough
					case "terminating":
						fallthrough
					case "terminated":
						statusFlag = 2
					case "notready":
						statusFlag = 3
					}
				}
				if len(sic) == q.Replicas {
					scf.Instance = sic
					scf.Deploy = 1
					scf.Msg = ""
					break
				}
				if errIdx == 3 {
					scf.Deploy = 4
					break
				}
				time.Sleep(3 * time.Second)
			}
		}
		if scf.Deploy == 1 {
			break
		}
		if scf.Deploy == 4 {
			errIdx = 3
		}
		time.Sleep(5 * time.Second)
	}

	logrus.WithFields(logrus.Fields{"Update ServiceConf": scf}).Info(ModuleName)

	if plugin != nil {
		plugin(scf, para)
	}
	logrus.WithFields(logrus.Fields{"After Update ServiceConf": scf.ToString()}).Info(ModuleName)

	svcconf.UpdateSvcConf(scf)
}
