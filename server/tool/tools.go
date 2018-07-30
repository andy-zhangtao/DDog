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
	"github.com/andy-zhangtao/gogather/znet"
	"time"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	zmodel "github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
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

// GetLocalIP  获取本地IP地址. 如果获取失败，则使用默认的127.0.0.1
func GetLocalIP() (string) {
	ip, err := znet.LocallIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	return ip
}

//GetZipKinSpan 获取ZipKin跟踪Span
//servername和funcname是服务名称和函数名称,用来显示定位信息
//traceid, id, parentid 对应于上游span SpanContext中的traceid, id, parentid
//此函数只能用在跟踪链的中游,即跟踪链路(Span)已经建立，此函数所处的服务是跟踪链接中间的一环
//如果需要构建一个新的span,需要将span = tracer.StartSpan(servername, zipkin.Parent(ctx))修改为span = tracer.StartSpan(servername)
func GetZipKinSpan(servername, funcname, traceid, id, parentid string) (span zipkin.Span, reporter reporter.Reporter, isCreate bool) {
	isCreate = false
	zipKinUrl := os.Getenv(_const.ENV_AGENT_ZIPKIN_ENDPOINT)
	if zipKinUrl != "" {
		reporter = httpreporter.NewReporter(fmt.Sprintf("%s/api/v2/spans", zipKinUrl))
		endpoint, err := zipkin.NewEndpoint(servername, GetLocalIP())
		if err != nil {
			logrus.WithFields(logrus.Fields{"Create ZipKin Endpoint Error": fmt.Sprintf("unable to create local endpoint: %+v\n", err)}).Error(ModuleName)
		} else {
			tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
			if err != nil {
				logrus.WithFields(logrus.Fields{"Create ZipKin Tracer Error": fmt.Sprintf("unable to create tracer: %+v\n", err)}).Error(ModuleName)
			} else {
				//为了还原成正确的traceid,需要在traceid前后各添加一个0
				//具体原因，参考traceid UnmarshalJSON源码
				traceid = fmt.Sprintf("0%s0", traceid)
				id = fmt.Sprintf("0%s0", id)
				ctx := zmodel.SpanContext{}
				_tracid := new(zmodel.TraceID)
				_tracid.UnmarshalJSON([]byte(traceid))
				ctx.TraceID = *_tracid

				_id := new(zmodel.ID)
				_id.UnmarshalJSON([]byte(id))
				ctx.ID = *_id

				_parentid := new(zmodel.ID)
				_parentid.UnmarshalJSON([]byte(parentid))
				ctx.ParentID = _parentid

				span = tracer.StartSpan(servername, zipkin.Parent(ctx))
				logrus.WithFields(logrus.Fields{"span": span}).Info(servername)
				span.Annotate(time.Now(), fmt.Sprintf("%s-%s Receive Request", servername, funcname))
				isCreate = true
			}
		}
	}
	return
}
