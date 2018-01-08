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
)

type Container struct {
	ID   bson.ObjectId     `json:"id,omitempty" bson:"_id,omitempty"`
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
	Svc  string            `json:"svc"`
	Nsme string            `json:"namespace"`
	Idx  int               `json:"idx"`
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var con Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkContainer(con); err != nil {
		tool.ReturnError(w, err)
		return
	}

	sv, err := mongo.GetSvcConfByName(con.Svc, con.Nsme)
	if sv == nil {
		tool.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	con.ID = bson.NewObjectId()
	if err = mongo.SaveContainer(con); err != nil {
		tool.ReturnError(w, err)
		return
	}

	w.Write([]byte(con.ID.Hex()))
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

func checkContainer(con Container) error {
	if con.Name == "" {
		return errors.New(_const.NameNotFound)
	}

	if con.Img == "" {
		return errors.New(_const.ImageNotFounc)
	}

	if con.Svc == "" {
		return errors.New(_const.HttpSvcEmpty)
	}

	if con.Nsme == "" {
		return errors.New(_const.NamespaceNotFound)
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

	var con Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if err = checkContainer(con); err != nil {
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
