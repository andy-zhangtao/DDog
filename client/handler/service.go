package handler

import (
	"net/http"
	"github.com/andy-zhangtao/DDog/const"
	"errors"
	"github.com/andy-zhangtao/DDog/server"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"encoding/json"
)

func QueryService(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	if ns == "" {
		server.ReturnError(w, errors.New(_const.NamespaceNotFound))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	name := r.URL.Query().Get("name")
	if name == "" {
		svc, err := mongo.GetAllSvcByNs(ns)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(svc)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		w.Write(data)
	} else {
		svc, err := mongo.GetSvcByName(ns, name)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		data, err := json.Marshal(svc)
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		w.Write(data)
	}
}
