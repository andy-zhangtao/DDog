package dns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/andy-zhangtao/DDog/server"
	"github.com/andy-zhangtao/DDog/server/etcd"
	gg "github.com/andy-zhangtao/gogather/strings"
	"strings"
)

const (
	SaveMethod = iota
	DeleMethod
)

type DnsMeteData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SaveDNS 保存DNS数据
// Key为需要解析的域名
// Value为相对应的A记录
func SaveDNS(w http.ResponseWriter, r *http.Request) {

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

// DeleDNS 删除指定DNS记录
// Key为需要删除的域名
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

// GetDNS 获取指定DNS记录
// Key为需要查询的域名,支持模糊查询
func GetDNS(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		server.ReturnError(w, errors.New("Domain can not be empty"))
		return
	}

	isFuzzy := false
	fuzzy := r.URL.Query().Get("fuzzy")
	if strings.ToLower(fuzzy) == "true" {
		isFuzzy = true
	}

	w.Header().Set("Content-Type", "application/json")
	var resp []byte
	if !isFuzzy {
		data, err := etcd.Get(parseDomain(domain), []string{})
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		resp, _ = json.Marshal(&DnsMeteData{
			Key:   domain,
			Value: data[parseDomain(domain)],
		})
	} else {
		var td []DnsMeteData
		data, err := etcd.Get(parseDomain(domain), []string{"--from-key"})
		if err != nil {
			server.ReturnError(w, err)
			return
		}

		for key, value := range data {
			k := gg.ReverseWithSeg(key, "/", ".")
			td = append(td, DnsMeteData{
				Key:   k[0:len(k)-1],
				Value: value,
			})
		}
		resp, _ = json.Marshal(td)
	}
	w.Write(resp)
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
		if dmd.Key == "" {
			return dmd, errors.New("Key can not be empty!")
		}

		if dmd.Value == "" {
			return dmd, errors.New("Value can not be empty!")
		}
	case DeleMethod:
		if dmd.Key == "" {
			return dmd, errors.New("Key can not be empty!")
		}
	}

	return dmd, nil
}

func parseDomain(domain string) string {
	return "/" + gg.ReverseWithSeg(domain, ".", "/")
}