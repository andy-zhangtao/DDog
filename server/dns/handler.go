package dns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/andy-zhangtao/DDog/server"
	"github.com/andy-zhangtao/DDog/server/etcd"
	gg "github.com/andy-zhangtao/gogather/strings"
)

const (
	SaveMethod = iota
	DeleMethod
)

type DnsMeteData struct {
	//Domain string `json:"domain"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SaveDns(w http.ResponseWriter, r *http.Request) {

	dmd, err := getDNS(r, SaveMethod)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	err = etcd.Put(parseDomain(dmd.Key), dmd.Value)
	if err != nil {
		server.ReturnError(w, err)
		return
	}
}

func DeleDNS(w http.ResponseWriter, r *http.Request) {
	dmd, err := getDNS(r, DeleMethod)
	if err != nil {
		server.ReturnError(w, err)
		return
	}

	err = etcd.Dele(parseDomain(dmd.Key))
	if err != nil {
		server.ReturnError(w, err)
		return
	}
}

func getDNS(r *http.Request, method int) (DnsMeteData, error) {
	var dmd DnsMeteData
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return dmd, err
	}

	err = json.Unmarshal(data, &dmd)
	if err != nil {
		return dmd, err
	}

	switch method {
	case SaveMethod:
		//if dmd.Domain == "" {
		//	return dmd, errors.New("Domain can not be empty!")
		//}

		if dmd.Key == "" {
			return dmd, errors.New("Key can not be empty!")
		}

		if dmd.Value == "" {
			return dmd, errors.New("Value can not be empty!")
		}
	case DeleMethod:
		//if dmd.Domain == "" {
		//	return dmd, errors.New("Domain can not be empty!")
		//}

		if dmd.Key == "" {
			return dmd, errors.New("Key can not be empty!")
		}
	}

	return dmd, nil
}

func parseDomain(domain string) string {
	return "/" + gg.ReverseWithSeg(domain, ".", "/")
}
