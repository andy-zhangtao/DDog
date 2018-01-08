package dns

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/andy-zhangtao/DDog/server/etcd"
	gg "github.com/andy-zhangtao/gogather/strings"
	"strings"
	"github.com/andy-zhangtao/DDog/server/tool"
)

const (
	SaveMethod = iota
	DeleMethod
)

type DnsMeteData struct {
	Key   string `json:"key"`
	Value vsvc   `json:"value"`
}

type vsvc struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// SaveDNS 保存DNS数据
// Key为需要解析的域名
// Value为相对应的A记录
func SaveDNS(w http.ResponseWriter, r *http.Request) {

	dmd, err := getDNS(r, SaveMethod)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	value, err := json.Marshal(dmd.Value)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	err = etcd.Put(parseDomain(dmd.Key), string(value))
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
}

// DeleDNS 删除指定DNS记录
// Key为需要删除的域名
func DeleDNS(w http.ResponseWriter, r *http.Request) {
	dmd, err := getDNS(r, DeleMethod)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	err = etcd.Dele(parseDomain(dmd.Key))
	if err != nil {
		tool.ReturnError(w, err)
		return
	}
}

// GetDNS 获取指定DNS记录
// Key为需要查询的域名,支持模糊查询
func GetDNS(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		tool.ReturnError(w, errors.New("Domain can not be empty"))
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
			tool.ReturnError(w, err)
			return
		}

		value := data[parseDomain(domain)]

		var v vsvc

		err = json.Unmarshal([]byte(value), &v)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		resp, _ = json.Marshal(&DnsMeteData{
			Key:   domain,
			Value: v,
		})
	} else {
		var td []DnsMeteData
		data, err := etcd.Get(parseDomain(domain), []string{"--from-key"})
		if err != nil {
			tool.ReturnError(w, err)
			return
		}

		for key, value := range data {
			var v vsvc

			err = json.Unmarshal([]byte(value), &v)
			if err != nil {
				tool.ReturnError(w, err)
				return
			}

			k := gg.ReverseWithSeg(key, "/", ".")
			td = append(td, DnsMeteData{
				Key:   k[0:len(k)-1],
				Value: v,
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

		if dmd.Value.Host == "" {
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
