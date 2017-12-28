// goinit.go 用来初始化coredns配置文件
package watch

import (
	"os"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"log"
	"text/template"
	"github.com/andy-zhangtao/DDog/server/etcd"
	"strings"
	"github.com/andy-zhangtao/DDog/client/handler"
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"time"
	"github.com/andy-zhangtao/DDog/bridge"
)

type watchdog struct {
	sid    string
	skey   string
	region []string
}

func Go() {
	err := genConfigure()
	if err != nil {
		log.Panicln(err)
	}

	w, err := getMetaData()
	if err != nil {
		log.Panicln(err)
	}

	for {
		now := time.Now()
		next := now.Add(time.Minute * 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		log.Printf("下次采集时间为[%s]\n", next.Format("200601021504"))

		select {
		case <-t.C:
			if w.sid != "" {
				watch(w)
			}else{
				log.Println("当前MetaData数据为空")
			}
		case <- bridge.GetMetaChan():
			w, err = getMetaData()
			if err != nil{
				log.Printf("获取MetaData失败[%s]\n",err.Error())
			}
		}
	}

}

// genConfigure 生成DNS配置文件
func genConfigure() error {
	type conf struct {
		Domain string
		Etcd string
	}

	name := os.Getenv(_const.EnvDomain)
	if name == "" {
		return errors.New(_const.EnvDomainNotFound)
	}

	etcd := os.Getenv(_const.EnvEtcd)
	if etcd == "" {
		return errors.New(_const.EnvEtcdNotFound)
	}

	path := os.Getenv(_const.EnvConfPath)
	if path == "" {
		path = "/"
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	var mf = conf{
		Domain: name,
		Etcd: etcd,
	}

	t := template.Must(template.New("makefile").Parse(_const.Corefile))

	file, err := os.Create(path + "corefile")
	if err != nil {
		return err
	}

	err = t.Execute(file, mf)
	if err != nil {
		return err
	}

	return nil
}

func watch(w *watchdog) {
	watchCluster(w)
}

// watchCluster 监听集群数据变化
func watchCluster(w *watchdog) {
	ch := handler.Cluster{
		SecretId:  w.sid,
		SecretKey: w.skey,
	}
	for _, r := range w.region {
		ch.Region = r
		cinfo, err := ch.SaveClusterInfo(true)
		if err != nil {
			log.Printf("[%s]Get Cluster Info Error [%s]\n", ch.Region, err.Error())
			continue
		}

		for _, c := range cinfo.Data.Clusters {
			nsInfo, err := watchNS(c, w)
			if err != nil {
				log.Printf("[%s]Get Namespace Info Error [%s]\n", c.ClusterName, err.Error())
				continue
			}

			for _, n := range nsInfo.Data.Namespaces {
				err = watchSVC(n, c, w)
				if err != nil {
					log.Printf("[%s]Maintain Svc <-> DNS Error [%s]\n", n.Name, err.Error())
					return
				}
			}
		}
	}
}

// watchNS  监听命名空间变化
func watchNS(c cvm.ClusterInfo_data_clusters, w *watchdog) (*namespace.NSInfo, error) {
	ns := handler.NameSpace{
		SecretId:  w.sid,
		SecretKey: w.skey,
		Region:    c.Region,
		ClusterID: c.ClusterId,
	}

	return ns.SaveNSInfo(true)
}

// watchSVC  监听服务变化
func watchSVC(n namespace.NSInfo_data_namespaces, c cvm.ClusterInfo_data_clusters, w *watchdog) error {
	svc := handler.Svc{
		SecretId:  w.sid,
		SecretKey: w.skey,
		Region:    c.Region,
		Clusterid: c.ClusterId,
		Namespace: n.Name,
	}
	return svc.WatchDNS()
}

// getMetaData 读取密钥ID,密钥值和机房区域信息
// 当三者发生变化时，需要重新调用来刷新数据
func getMetaData() (*watchdog, error) {

	sid, err := etcd.Get(_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo, nil)
	if err != nil {
		return nil, err
	}

	skey, err := etcd.Get(_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo, nil)
	if err != nil {
		return nil, err
	}

	region, err := etcd.Get(_const.CloudEtcdRootPath+_const.CloudEtcdRegionInfo, nil)
	if err != nil {
		return nil, err
	}

	w := watchdog{
		sid:    sid[_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo],
		skey:   skey[_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo],
		region: strings.Split(region[_const.CloudEtcdRootPath+_const.CloudEtcdRegionInfo], ";"),
	}

	//if w.sid == "" {
	//	return nil, errors.New("Can not find secret_id in etcd")
	//}
	//
	//if w.skey == "" {
	//	return nil, errors.New("Can not find secret_key in etcd")
	//
	//}
	//
	//if len(w.region) == 0 {
	//	return nil, errors.New("The region number error!")
	//}

	return &w, nil
}
