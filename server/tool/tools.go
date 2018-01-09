package tool

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/andy-zhangtao/DDog/const"
)

type HttpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func ReturnError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("EQXC-Run-Svc", "500")
	w.WriteHeader(500)
	data, er := json.Marshal(&HttpError{
		Code: 500,
		Msg:  err.Error(),
	})

	if er != nil {
		w.Write([]byte("{code:500,msg:" + err.Error() + "}"))
	} else {
		w.Write(data)
	}
}

//func UnMarshaSvcConf(conf interface{}) (svcconf svcconf.SvcConf, err error) {
//	data, err := bson.Marshal(conf)
//	if err != nil {
//		return
//	}
//
//	err = bson.Unmarshal(data, &svcconf)
//	if err != nil {
//		return
//	}
//
//	return
//}

// IsNotFound 判断返回的错误是否是数据库无记录
func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}

func ReturnResp(w http.ResponseWriter, data []byte) {
	msg := _const.RespMsg{
		Code: 1000,
		Msg:  _const.OperationSucc,
		Data: string(data),
	}
	if len(data) == 4 && strings.ToLower(string(data)) == "null" {
		msg.Code = 0
		msg.Msg = _const.DataNotFound
	}

	d, err := json.Marshal(&msg)
	if err != nil {
		ReturnError(w, err)
		return
	}

	w.Write(d)
}
