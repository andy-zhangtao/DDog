package container

import (
	"net/http"
	"io/ioutil"
	"github.com/andy-zhangtao/DDog/server"
	"encoding/json"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Container struct {
	ID   bson.ObjectId     `json:"id,omitempty" bson:"_id,omitempty"`
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
	Svc  string            `json:"svc"`
	Nsme string            `json:"namespace"`
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	var con Container
	err = json.Unmarshal(data, &con)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	if err = checkContainer(con); err != nil {
		server.ReturnError(w, err)
		return
	}

	sv, err := mongo.GetSvcByName(con.Nsme, con.Svc)
	if sv == nil {
		server.ReturnError(w, errors.New(_const.SVCNoExist))
		return
	}

	con.ID = bson.NewObjectId()
	if err = mongo.SaveContainer(con); err != nil {
		server.ReturnError(w, err)
		return
	}

	w.Write([]byte(con.ID.Hex()))
	return
}

func GetContainer(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		server.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svc := r.URL.Query().Get("svc")
	if svc == "" {
		server.ReturnError(w, errors.New(_const.HttpSvcEmpty))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("cid")
	if id == "" {
		con, err := mongo.GetContaienrBySvc(svc, ns)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(con)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		w.Write(data)
	} else {
		con, err := mongo.GetContainerByID(id)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(con)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		w.Write(data)
	}
}

func DeleteContainer(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		server.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	svc := r.URL.Query().Get("svc")
	if svc == "" {
		server.ReturnError(w, errors.New(_const.HttpSvcEmpty))
		return
	}

	id := r.URL.Query().Get("cid")
	if id == "" {
		err := mongo.DeleteAllContainer(svc, ns)
		if err != nil {
			server.ReturnError(w, err)
			return
		}
	} else {
		err := mongo.DeleteContainerById(id)
		if err != nil {
			server.ReturnError(w, err)
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

	return nil
}
