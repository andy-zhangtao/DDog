package tool

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/andy-zhangtao/DDog/const"
	"os"
	"fmt"
	"errors"
	"io/ioutil"
	"net/url"
	"github.com/sirupsen/logrus"
)

const (
	ModuleName = "Tools"
)

type HttpError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// ReturnError 返回常规错误
// 响应码为500
// 如果传入的错误信息不是json格式，将按照字符串格式返回
// 如果传入的错误信息符合json格式，将按照json格式返回
func ReturnError(w http.ResponseWriter, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("EQXC-Run-Svc", "500")
	w.WriteHeader(500)
	data, er := json.Marshal(&HttpError{
		Code: 500,
		Msg:  err.Error(),
	})

	if er != nil {
		w.Write([]byte("{\"code\":500,\"msg\":\"" + err.Error() + "\"}"))
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
	if err != nil {
		return strings.Contains(err.Error(), "not found")
	}
	return false
}

// ReturnResp 返回非代码层面的错误. 此错误经常用在返回服务执行成功
// 但执行的结果是失败的情况
// 例如服务执行成功，但对方没有返回任何数据，此时返回的响应码为200，但
// 追加一个没有查询到数据的返回信息
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

// InspectImgInfo 获取镜像信息
// conname 容器名称, 用于回调时查找容器配置
// svcname 服务配置名称，用于回调时确定服务配置
// namespace 命名空间名称, 回调时使用
// imgname 镜像名称,用于检索镜像数据
func InspectImgInfo(conname, svcname, namespace, imgname string, callback func(error)) error {
	goblin := os.Getenv(_const.EnvGoblin)
	if goblin == "" {
		return errors.New(fmt.Sprintf("[%s]Empty!\n", _const.EnvGoblin))
	}

	errChan := make(chan error)

	go func() {
		imgname = url.QueryEscape(imgname)
		path := fmt.Sprintf("http://%s/v1/inspect?name=%s&svc=%s&namespace=%s&img=%s", goblin, conname, svcname, namespace, imgname)
		fmt.Printf("Invoke[%s]\n ", path)
		resp, err := http.Get(path)
		if err != nil {
			errChan <- err
			return
		}

		if resp.StatusCode != 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf("Goblin Resp [%d] [%s] \n", resp.StatusCode, err))
			}
			errChan <- errors.New(fmt.Sprintf("Goblin Resp [%d] [%s] \n", resp.StatusCode, string(body)))
			return
		}
		errChan <- nil

	}()
	err := <-errChan
	logrus.WithFields(logrus.Fields{"Ready to invoke callback": svcname}).Info(ModuleName)
	callback(err)
	return nil
}
