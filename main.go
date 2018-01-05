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
	"github.com/andy-zhangtao/DDog/server/metadata"
	_ "github.com/andy-zhangtao/DDog/const"
	"os"
	"github.com/andy-zhangtao/DDog/server/container"
	"github.com/andy-zhangtao/DDog/server/svcconf"
)

var _VERSION_ = "unknown"
var _APIVERSION_ = "/v1"

func main() {
	log.Println(getVersion())
	region := os.Getenv(_const.EnvRegion)
	if region == "" {
		log.Panic(_const.EnvRegionNotFound)
	}
	//go watch.Go(region)
	r := mux.NewRouter()
	r = dnsAPI(r)
	r = cloudAPI(r)
	r = metadataAPI(r)
	r = serviceAPI(r)
	r = namespaceAPI(r)
	r = containerAPI(r)

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

func cloudAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.GetNodeInfo), qcloud.GetClusterNodes).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.AddSvcIP), handler.AddSvcDnsAR).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetClusterInfo), handler.QueryClusterInfo).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetClusterInfo), handler.GetClusterInfo).Methods(http.MethodGet)
	//r.HandleFunc(getApiPath(_const.GetSvcMoreInfo), qcloud.GetSampleSVCInfo).Methods(http.MethodGet)
	return r
}

func metadataAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.MetaData), metadata.Startup).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.MetaData), metadata.GetMetaDataWithHttp).Methods(http.MethodGet)
	return r
}
func serviceAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.GetSvcSampleInfo), handler.QueryService).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.NewSvcConfig), svcconf.CreateSvcConf).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetSvcConfig), svcconf.GetSvcConf).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.DeleteSvcConfig), svcconf.DeleteSvcConf).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.RunService), qcloud.RunService).Methods(http.MethodPost)
	return r
}
func namespaceAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.NewNameSpace), qcloud.CreateNamespace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DeleteNameSpace), qcloud.Deletenamespace).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.GetNSInfo), handler.QueryNameSpace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetNSInfo), handler.QueryNamespaceByName).Methods(http.MethodGet)
	return r
}
func dnsAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.SaveDNS).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.DeleDNS).Methods(http.MethodDelete)
	r.HandleFunc(getApiPath(_const.DnsMetaData), dns.GetDNS).Methods(http.MethodGet)
	return r
}

func containerAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.NewContainer), container.CreateContainer).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetContainer), container.GetContainer).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.DeleteContainer), container.DeleteContainer).Methods(http.MethodPost)
	return r
}
