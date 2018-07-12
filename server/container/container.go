package container

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"strconv"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/DDog/model/container"
	"strings"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"os"
	"github.com/sirupsen/logrus"
	"fmt"
)

const (
	ModuleName = "Container Operation"
)

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var con container.Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{"Request Body": string(data)}).Info(ModuleName)

	_, isExist, err := isExistContainer(&con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	if isExist {
		tool.ReturnError(w, errors.New(_const.ConConfExist))
	} else {
		err = createContainer(&con)
	}
	return
}

func GetContainer(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svc := r.URL.Query().Get("svc")
	if svc == "" {
		tool.ReturnError(w, errors.New(_const.HttpSvcEmpty))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("cid")
	if id == "" {
		con, err := mongo.GetContaienrBySvc(svc, ns)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(con)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		w.Write(data)
	} else {
		con, err := mongo.GetContainerByID(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(con)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		w.Write(data)
	}
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("cid")
	if id == "" {
		ns := r.URL.Query().Get("namespace")
		if ns == "" {
			tool.ReturnError(w, errors.New(_const.NamespaceNotFound))
			return
		}

		svc := r.URL.Query().Get("svc")
		if svc == "" {
			tool.ReturnError(w, errors.New(_const.HttpSvcEmpty))
			return
		}
		err := mongo.DeleteAllContainer(svc, ns)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	} else {
		err := mongo.DeleteContainerById(id)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}
}

// checkContainer 检查容器配置参数合法性
// Img 必填项
// Svc 必填项
// Name 非必填项, 如果为空，则根据镜像名称生成。 规则如下:
//		1. 如果镜像名为 domain，则名称为domain
//		2. 如果镜像名为 domain/name 则名称为name
//		3. 如果镜像名为 domain/name/img 则名称为img
// Nsme 命名空间 非必填项 若为空，则取默认值
// Idx 非必填项 若为空，则默认为1
func checkContainer(con *container.Container) error {

	if con.Img == "" {
		return errors.New(_const.ImageNotFounc)
	}

	if con.Name == "" {
		img := strings.Split(con.Img, "/")
		if len(img) == 1 {
			con.Name = img[0]
		} else if len(img) == 2 {
			con.Name = img[1]
		} else if len(img) == 3 {
			con.Name = img[2]
		}

		if strings.Contains(con.Name, ":") {
			tn := strings.Split(con.Name, ":")
			con.Name = tn[0]
		}
	}

	if con.Svc == "" {
		return errors.New(_const.HttpSvcEmpty)
	}

	if con.Nsme == "" {
		con.Nsme = _const.DefaultNameSpace
	}

	if con.Idx == 0 {
		con.Idx = 1
	}

	//获取不同环境的日志参数
	//日志参数规则如下：
	//"--log-opt gelf-address=default:udp://192.168.0.8:12201|+|proenv:udp://xxxxxxxxxxx|+|; --log-opt buf=1; --log-opt env=svcname;"
	var logOpt string
	var logOpts []string

	lops := strings.Split(os.Getenv(_const.EnvDefaultLogOpt), ";")
	if len(lops) > 0 {
		logOpts = lops[1 : len(lops)-1]

		for _, l := range lops {
			if strings.HasPrefix(strings.TrimSpace(l), "\"--log-opt") {
				ops := strings.Split(l, "=")
				if len(ops) > 1 {
					logOpt = ops[1]
				} else {
					logOpt = os.Getenv(_const.EnvDefaultLogOpt)
					logrus.WithFields(logrus.Fields{"GelfAddress Value Error": l}).Error(ModuleName)
				}
				break
			}
		}
	} else {
		logOpt = os.Getenv(_const.EnvDefaultLogOpt)
		logrus.WithFields(logrus.Fields{"LogOpt Value Error": os.Getenv(_const.EnvDefaultLogOpt)}).Error(ModuleName)
	}

	switch con.Nsme {
	case "proenv":
		_logOpt := strings.Split(logOpt, "|+|")
		for _, l := range _logOpt {
			if strings.HasPrefix(l, "proenv:") {
				logOpt = fmt.Sprintf("\"--log-opt gelf-address=%s;%s;\"", strings.Split(l, "proenv:")[1], strings.Join(logOpts, ";"))
			}
		}
	default:
		_logOpt := strings.Split(logOpt, "|+|")
		for _, l := range _logOpt {
			if strings.HasPrefix(l, "default:") {
				logOpt = fmt.Sprintf("\"--log-opt gelf-address=%s;%s;\"", strings.Split(l, "default:")[1], strings.Join(logOpts, ";"))
			}
		}
	}

	logrus.WithFields(logrus.Fields{"namespace": con.Nsme, "log-opt": logOpt}).Info(ModuleName)
	if con.Env == nil {
		con.Env = map[string]string{
			"LOGCHAIN_DRIVER": os.Getenv(_const.EnvDefaultLogDriver),
			"svcname":         con.Svc,
			"log_opt":         logOpt,
		}
	} else {
		con.Env["LOGCHAIN_DRIVER"] = os.Getenv(_const.EnvDefaultLogDriver)
		con.Env["svcname"] = con.Svc
		con.Env["log_opt"] = logOpt
	}

	return nil
}

func UpgradeContainer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var con container.Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{"Request Body": string(data), "Operation": "UpgradeContainer"}).Info(ModuleName)

	_, isExist, err := isExistContainer(&con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if isExist {
		rmall := r.URL.Query().Get("rmall")
		isRmall, err := strconv.ParseBool(rmall)
		if err != nil {
			isRmall = false
		}

		if isRmall {
			err = container.DeleteAllContaienrUnderSvc(con.Svc, con.Nsme)
			if err != nil {
				return
			}
		}
		err = upgreadeContaienr(&con)
	} else {
		err = errors.New(_const.SVCNoExist)
	}

	if err != nil {
		tool.ReturnError(w, err)
		return
	}
	return
}

func IsExistContainer(con *container.Container) (old *container.Container, isExist bool, err error) {
	return isExistContainer(con)
}

// isExistContainer 检查容器配置是否存在
// 如果存在则返回TRUE，否则返回FALSE
func isExistContainer(con *container.Container) (old *container.Container, isExist bool, err error) {
	logrus.WithFields(logrus.Fields{"Request Body": con, "Operation": "IsExistContainer"}).Info(ModuleName)

	if err = checkContainer(con); err != nil {
		return
	}

	logrus.WithFields(logrus.Fields{"After Check": con, "Operation": "IsExistContainer"}).Info(ModuleName)

	tcon, err := container.GetContainerByName(con.Name, con.Svc, con.Nsme)
	if err != nil {
		return
	}

	return tcon, !(tcon == nil), nil
}

// CheckContainer 确认容器配置数据是否存在
// 如果存在则更新，如果不存在则创建
func CheckContainer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var con container.Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	logrus.WithFields(logrus.Fields{"Request Body": string(data), "Operation": "CheckContainer"}).Info(ModuleName)

	var nt []container.NetConfigure
	for _, p := range con.Port {
		nt = append(nt, container.NetConfigure{
			AccessType: 0,
			InPort:     p,
			OutPort:    p,
			Protocol:   0,
		})
	}

	con.Net = nt
	_, isExist, err := isExistContainer(&con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if isExist {
		isChange, err := container.UpgradeContainerNetByName(con.Name, con.Svc, con.Nsme, con.Net)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		logrus.WithFields(logrus.Fields{"Compare NetConfigure Is change?": isChange, "Operation": CheckContainer}).Info(ModuleName)

		if isChange {
			scf, err := svcconf.GetSvcConfByName(con.Svc, con.Nsme)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}
			err = svcconf.GenerateNetconifg(scf)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}
		}

		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		err = upgreadeContaienr(&con)
	} else {
		err = createContainer(&con)
	}

	if err != nil {
		tool.ReturnError(w, err)
	} else {
		tool.ReturnResp(w, []byte("Container Check Succ!"))
	}

	return
}

func CreateContainerForGraphQL(con *container.Container) (err error) {
	return createContainer(con)
}

// createContainer 创建容器配置信息
// 在创建的同时会调用Goblin来解析网络配置数据
func createContainer(con *container.Container) (err error) {
	err = upgreadeSvcConf(*con, true)
	if err != nil {
		return
	}

	return container.SaveContainer(con)
}

// upgreadeSvcConf 异步更新服务配置信息
// updataNet 是否需要同步更新网络信息，在更新容器配置信息时，已经更新过网络信息了，这里就不需要再更新了
func upgreadeSvcConf(con container.Container, updateNet bool) (err error) {
	sv, err := svcconf.GetSvcConfByName(con.Svc, con.Nsme)
	if sv == nil {
		err = errors.New(_const.SVCNoExist)
		return
	}
	sv.Status = 0
	if updateNet {
		portMap := make(map[int]int)
		for _, n := range sv.Netconf {
			portMap[n.InPort] = n.InPort
		}

		for _, c := range con.Net {
			if _, exists := portMap[c.InPort]; !exists {
				sv.Netconf = append(sv.Netconf, c)
			}
		}
	}
	err = svcconf.UpdateSvcConf(sv)
	return

}

// upgreadeContaienr 更新容器配置数据
func upgreadeContaienr(con *container.Container) (err error) {
	err = backupContainer(*con)
	if err != nil {
		return
	}

	err = container.UpgradeContaienrByName(con)
	if err != nil {
		return
	}

	err = upgreadeSvcConf(*con, false)
	return
}

// backupContainer 备份容器配置数据
func backupContainer(con container.Container) (err error) {
	scf, err := svcconf.GetSvcConfByName(con.Svc, con.Nsme)
	if err != nil {
		return
	}

	if scf.BackID == "" {
		scf.BackupSvcConf()
	} else {

		bscf, err := scf.GetBackSvcConf()
		if err != nil {
			return err
		}

		for _, cn := range bscf.BackContainer {
			if cn.Img == con.Img {
				cn.Img = con.Img
				cn.Cmd = con.Cmd
				cn.Env = con.Env
				cn.Idx = con.Idx
				cn.Net = con.Net
				return svcconf.UpdateSvcConf(bscf)
			}
		}
	}

	return
}
