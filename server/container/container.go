package container

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"github.com/andy-zhangtao/DDog/server/tool"
	"log"
	"github.com/andy-zhangtao/DDog/model/container"
	"strings"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"encoding/base64"
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

	if err = checkContainer(&con); err != nil {
		tool.ReturnError(w, err)
		return
	}

	if _const.DEBUG {
		log.Printf("[CreateContainer] Check Container Data :[%v] \n", con)
	}

	sv, err := svcconf.GetSvcConfByName(con.Svc, con.Nsme)
	if sv == nil {
		tool.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	tcon, err := container.GetContainerByName(con.Name, con.Svc, con.Nsme)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if tcon == nil {
		var callback = func(err error) {
			fmt.Println("svc conf", err)
			if err != nil {
				sv.Status = 3
			} else {
				sv.Status = 0
			}
			svcconf.UpdateSvcConf(sv)
			return
		}
		sv.Status = 1
		go func(callback func(error)) {
			img := base64.StdEncoding.EncodeToString([]byte(con.Img))
			err := tool.InspectImgInfo(con.Svc, con.Nsme, img, callback)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}
		}(callback)
		sv.Status = 2
		fmt.Printf("svc conf [%v]\n", sv)
		err = svcconf.UpdateSvcConf(sv)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		container.SaveContainer(&con)
		return
	} else {
		tool.ReturnError(w, errors.New(_const.ConConfExist))
	}
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

	if err = checkContainer(&con); err != nil {
		tool.ReturnError(w, err)
		return
	}

	sv, err := mongo.GetSvcConfByName(con.Svc, con.Nsme)
	if sv == nil {
		tool.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	rmall := r.URL.Query().Get("rmall")
	isRmall, err := strconv.ParseBool(rmall)
	if err != nil {
		isRmall = false
	}

	if isRmall {
		err := mongo.DeleteAllContainer(con.Svc, con.Nsme)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
	}
	con.ID = bson.NewObjectId()
	if err = mongo.SaveContainer(con); err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Write([]byte(con.ID.Hex()))
	return
}
