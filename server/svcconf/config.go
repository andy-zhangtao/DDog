package svcconf

import (
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/server/mongo"
)

type SvcConf struct {
	Id         bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name"`
	Desc       string        `json:"desc"`
	Replicas   int           `json:"replicas"`
	AccessType int           `json:"access_type"`
	Namespace  string        `json:"namespace"`
}

func CreateSvcConf(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var conf SvcConf

	err = json.Unmarshal(data, &conf)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	if err = checkConf(conf); err != nil {
		server.ReturnError(w, err)
		return
	}

	cf, err := mongo.GetSvcConfByName(conf.Name, conf.Namespace)
	if cf != nil {
		server.ReturnError(w, errors.New(_const.SvcConfExist))
		return
	}

	if conf.Replicas == 0 {
		conf.Replicas = 1
	}

	conf.Id = bson.NewObjectId()
	if err = mongo.SaveSvcConfig(conf); err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Write([]byte(conf.Id.Hex()))
	return
}

func GetSvcConf(w http.ResponseWriter, r *http.Request) {
	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		server.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		conf, err := mongo.GetSvcConfNs(nsme)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(conf)
		if err != nil {
			server.ReturnError(w, err)
			return
		}
		w.Write(data)
	} else {
		conf, err := mongo.GetSvcConfByID(id)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(conf)
		if err != nil {
			server.ReturnError(w, err)
			return
		}
		w.Write(data)
	}

}

func DeleteSvcConf(w http.ResponseWriter, r *http.Request) {
	nsme := r.URL.Query().Get("namespace")
	if nsme == "" {
		server.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	if id == "" {
		err := mongo.DeleteSvcConfByNs(nsme)
		if err != nil {
			server.ReturnError(w, err)
			return
		}
	} else {
		err := mongo.DeleteSvcConfById(id)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

	}
}
func checkConf(conf SvcConf) error {
	if conf.Name == "" {
		return errors.New(_const.NameNotFound)
	}
	if conf.Namespace == "" {
		return errors.New(_const.NamespaceNotFound)
	}

	return nil
}
