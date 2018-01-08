package metadata

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"errors"
	"github.com/andy-zhangtao/qcloud_api/v1/cvm"
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/tool"
)

type metaData struct {
	Sid    string `json:"secret_id"`
	Skey   string `json:"secret_key"`
	Region string `json:"region"`
}

func Startup(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var md metaData
	err = json.Unmarshal(data, &md)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	//if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo, md.Sid); err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo, md.Skey); err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}
	//
	//if err = etcd.Put(_const.CloudEtcdRootPath+_const.CloudEtcdRegionInfo, md.Region); err != nil {
	//	tool.ReturnError(w, err)
	//	return
	//}

	if count, err := mongo.FindMetaDataByRegion(md.Region); err != nil {
		tool.ReturnError(w, err)
		return
	} else if count > 0 {
		tool.ReturnError(w, errors.New(_const.MetaDataDupilcate))
		return
	} else if err = mongo.SaveMetaData(md); err != nil {
		tool.ReturnError(w, err)
		return
	}

	bridge.GetMetaChan() <- 1
	return
}

// GetMetaData 获取存储在etcd中的密钥数据
func GetMetaData(region string) (metaData, error) {
	var md metaData

	//if keys, err := etcd.Get(_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo, nil); err != nil {
	//	return md, err
	//} else {
	//	md.Sid = keys[_const.CloudEtcdRootPath+_const.CloudEtcdSidInfo]
	//}
	//
	//if keys, err := etcd.Get(_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo, nil); err != nil {
	//	return md, err
	//} else {
	//	md.Skey = keys[_const.CloudEtcdRootPath+_const.CloudEtcdSkeyInfo]
	//}

	err := mongo.GetMetaDataByRegion(region, &md)
	if err != nil {
		return md, err
	}
	if md.Sid == "" || md.Skey == "" || md.Region == "" {
		return md, errors.New(region + " Metadata 获取为空")
	}
	return md, nil
}

func GetMetaDataWithHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	region := r.URL.Query().Get("region")
	if region != "" {
		md, err := GetMetaData(region)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		data, err := json.Marshal(&md)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		w.Write(data)
	} else {
		mds, err := mongo.GetALlMetaData()
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		data, err := json.Marshal(mds)
		if err != nil {
			tool.ReturnError(w, err)
			return
		}
		w.Write(data)
	}

}

// GetMdByClusterID 通过clusterid获取密钥数据
// 返回密钥ID，密钥值和所在区域
func GetMdByClusterID(clusterid string) (md metaData, err error) {
	var cluster cvm.ClusterInfo_data_clusters

	cs, err := mongo.GetClusterById(clusterid)
	if err != nil {
		return
	}

	data, err := bson.Marshal(cs)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &cluster)
	if err != nil {
		return
	}

	md, err = GetMetaData(_const.RegionMap[cluster.Region])
	if err != nil {
		return
	}
	return
}
