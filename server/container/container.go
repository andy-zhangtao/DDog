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
	"log"
	"github.com/andy-zhangtao/DDog/model/container"
	"strings"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"fmt"
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

	if _const.DEBUG {
		log.Printf("[CreateContainer] Receive Request Body:[%s] \n", string(data))
	}

	_,isExist, err := isExistContainer(&con)
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

	if _const.DEBUG {
		log.Printf("[CreateContainer] Receive Request Body:[%s] \n", string(data))
	}

	_,isExist, err := isExistContainer(&con)
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

	return
}

// isExistContainer 检查容器配置是否存在
// 如果存在则返回TRUE，否则返回FALSE
func isExistContainer(con *container.Container) (old *container.Container, isExist bool, err error) {
	if _const.DEBUG {
		log.Printf("[IsExistContainer] Receive Request Body:[%v] \n", con)
	}

	if err = checkContainer(con); err != nil {
		return
	}

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

	if _const.DEBUG {
		log.Printf("[CreateContainer] Receive Request Body:[%s] \n", string(data))
	}

	old, isExist, err := isExistContainer(&con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if isExist {
		con.Net = old.Net
		err = upgreadeContaienr(&con)
	} else {
		err = createContainer(&con)
	}

	return
}

// createContainer 创建容器配置信息
// 在创建的同时会调用Goblin来解析网络配置数据
func createContainer(con *container.Container) (err error) {
	err = upgreadeSvcConf(*con)
	if err != nil {
		return
	}
	err = container.SaveContainer(con)
	return
}

// upgreadeSvcConf 异步更新服务配置信息
func upgreadeSvcConf(con container.Container) (err error) {
	sv, err := svcconf.GetSvcConfByName(con.Svc, con.Nsme)
	if sv == nil {
		err = errors.New(_const.SVCNoExist)
		return
	}
	var callback = func(err error) {
		//fmt.Println("svc conf", err)
		if err != nil {
			sv.Status = 3
			sv.Msg = err.Error()
			svcconf.UpdateSvcConf(sv)
		}

		return
	}

	sv.Status = 1
	go func(callback func(error)) {
		err := tool.InspectImgInfo(con.Name, con.Svc, con.Nsme, con.Img, callback)
		if err != nil {
			fmt.Printf("********[%s] CallBack Error[%s]********\n", con.Img, err.Error())
			return
		}
	}(callback)
	sv.Status = 2
	//fmt.Printf("svc conf [%v]\n", sv)
	err = svcconf.UpdateSvcConf(sv)
	return

}

// upgreadeContaienr 更新容器配置数据
func upgreadeContaienr(con *container.Container) (err error) {
	err = container.UpgradeContaienrByName(con)
	if err != nil {
		return
	}

	err = upgreadeSvcConf(*con)
	return
}
