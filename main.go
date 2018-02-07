package main

import (
	"net/http"
	//"github.com/andy-zhangtao/DDog/server/dns"
	//_ "github.com/andy-zhangtao/DDog/server/etcd"
	_ "github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/client/handler"
	"github.com/andy-zhangtao/DDog/server/metadata"
	//_ "github.com/andy-zhangtao/DDog/const"
	_ "github.com/andy-zhangtao/DDog/pre-check"
	"os"
	"github.com/andy-zhangtao/DDog/server/container"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"fmt"
	"github.com/Sirupsen/logrus"
)

var _VERSION_ = "-unknown-"
var _APIVERSION_ = "/v1"
var _INNER_VERSION_ = "v0.6.4"

const (
	ModuleName = "DDog Main"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	region := os.Getenv(_const.EnvRegion)
	if region == "" {
		logrus.WithFields(logrus.Fields{"Region Not Found": _const.EnvRegion}).Panic(ModuleName)
	}

	fmt.Println("===================")
	logrus.WithFields(logrus.Fields{"version": getVersion(),}).Info("DDOG VERSION")

	//go watch.Go(region)
	r := mux.NewRouter()
	r = dnsAPI(r)
	r = cloudAPI(r)
	r = metadataAPI(r)
	r = serviceAPI(r)
	r = namespaceAPI(r)
	r = containerAPI(r)
	r = svcorAPI(r)
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	logrus.Println(http.ListenAndServe(":8000", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}

func getVersion() string {
	return _INNER_VERSION_ +"-"+ _VERSION_
}

func getApiPath(url string) string {
	return _APIVERSION_ + url
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getVersion()))
}

func cloudAPI(r *mux.Router) *mux.Router {
	r.HandleFunc("/_ping", ping).Methods(http.MethodGet)
	r.HandleFunc(getApiPath("/_ping"), ping).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.GetNodeInfo), qcloud.GetClusterNodes).Methods(http.MethodGet)
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
	r.HandleFunc(getApiPath(_const.QuerySvcConfigStatus), svcconf.QuerySvcConf).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.DeleteSvcConfig), svcconf.DeleteSvcConf).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.UpgradeService), svcconf.UpgradeSvcConf).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.RunService), qcloud.RunService).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DeleteService), qcloud.DeleteService).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.ReinstallService), qcloud.ReinstallService).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.CheckSvcConfig), svcconf.CheckSvcConf).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DeploySvcConfig), qcloud.DeployService).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.QuerySvcStatus), qcloud.GetSampleSVCInfo).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.UpdateSvcConfig), svcconf.UpdateNetPort).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.RollUpService), qcloud.RollingUpServiceWithSvc).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.RollBackService), qcloud.RollBackService).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.ConfirmService), qcloud.ConfirmRollService).Methods(http.MethodPost)
	return r
}

func namespaceAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.NewNameSpace), qcloud.CreateNamespace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.DeleteNameSpace), qcloud.Deletenamespace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetNSInfo), handler.QueryNameSpace).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetNSInfo), handler.QueryNamespaceByName).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.CheckNameSpace), qcloud.CheckNamespace).Methods(http.MethodPost)
	return r
}

func dnsAPI(r *mux.Router) *mux.Router {
	//r.HandleFunc(getApiPath(_const.DnsMetaData), dns.SaveDNS).Methods(http.MethodPost)
	//r.HandleFunc(getApiPath(_const.DnsMetaData), dns.DeleDNS).Methods(http.MethodDelete)
	//r.HandleFunc(getApiPath(_const.DnsMetaData), dns.GetDNS).Methods(http.MethodGet)
	//r.HandleFunc(getApiPath(_const.AddSvcIP), handler.AddSvcDnsAR).Methods(http.MethodPost)
	return r
}

func containerAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.NewContainer), container.CreateContainer).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetContainer), container.GetContainer).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.DeleteContainer), container.DeleteContainer).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.UpgradeContainer), container.UpgradeContainer).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.CheckContainer), container.CheckContainer).Methods(http.MethodPost)
	return r
}

func svcorAPI(r *mux.Router) *mux.Router {
	r.HandleFunc(getApiPath(_const.AddSvcGroup), svcconf.AddSvcConfGroup).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.GetSvcGroup), svcconf.GetSvcConfGroup).Methods(http.MethodGet)
	r.HandleFunc(getApiPath(_const.DeleSvcGroup), svcconf.DeleteSvcConfGroup).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.RunSvcGroup), qcloud.RunSvcGroup).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.ReinstallSvcGroup), qcloud.ReinstallSvcGroup).Methods(http.MethodPost)
	r.HandleFunc(getApiPath(_const.UninstallSvcGroup), qcloud.UninstallSvcGroup).Methods(http.MethodPost)
	return r
}
