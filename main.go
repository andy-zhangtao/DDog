package main

import (
	"log"
	"net/http"

	"github.com/andy-zhangtao/DDog/server/dns"
	_ "github.com/andy-zhangtao/DDog/server/etcd"
	_ "github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/client/handler"
	"github.com/andy-zhangtao/DDog/watch"
	"github.com/andy-zhangtao/DDog/server/metadata"
	_ "github.com/andy-zhangtao/DDog/const"
)

var _VERSION_ = "unknown"
var _APIVERSION_ = "/v1"

func main() {
	log.Println(getVersion())
	go watch.Go()
	r := mux.NewRouter()
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.SaveDNS).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.DeleDNS).Methods(http.MethodDelete)
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.GetDNS).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.GetNodeInfo), qcloud.GetClusterNodes).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.AddSvcIP), handler.AddSvcDnsAR).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetClusterInfo), handler.QueryClusterInfo).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.GetNSInfo), handler.QueryNameSpace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetSvcSampleInfo), qcloud.GetSampleSVCInfo).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.AddMetaData), metadata.Startup).Methods(http.MethodPost)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	log.Println(http.ListenAndServe(":8000", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}

func getVersion() string {
	return _VERSION_
}

func getApiPath(url string) string {
	return _APIVERSION_ + url
}
