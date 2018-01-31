package metadata

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"errors"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/DDog/model/cluster"
)

// Startup 初始化MetaData数据
func Startup(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var tmd metadata.MetaData
	err = json.Unmarshal(data, &tmd)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	if md, err := metadata.GetMetaDataByRegion(tmd.Region); err != nil {
		tool.ReturnError(w, err)
		return
	} else if md != nil {
		tool.ReturnError(w, errors.New(_const.MetaDataDupilcate))
		return
	} else if err = metadata.SaveMetaData(tmd); err != nil {
		tool.ReturnError(w, err)
		return
	}

	bridge.GetMetaChan() <- 1
	return
}

// GetMetaData 获取存储在etcd中的密钥数据
func GetMetaData(region string) (metadata.MetaData, error) {
	//var md metadata.MetaData

	md, err := metadata.GetMetaDataByRegion(region)
	//err := mongo.GetMetaDataByRegion(region, &md)
	if err != nil {
		return *md, err
	}
	if md.Sid == "" {
		return *md, errors.New(region + " Metadata 获取为空")
	}

	return *md, nil
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
func GetMdByClusterID(clusterid string) (md metadata.MetaData, err error) {
	//var cluster cvm.ClusterInfo_data_clusters
	//
	//cs, err := mongo.GetClusterById(clusterid)
	//if err != nil {
	//	return
	//}
	//
	//data, err := bson.Marshal(cs)
	//if err != nil {
	//	return
	//}
	//
	//err = bson.Unmarshal(data, &cluster)
	//if err != nil {
	//	return
	//}

	cluster, err := cluster.GetClusterByID(clusterid)
	if err != nil {
		return
	}

	md, err = GetMetaData(_const.RegionMap[cluster.Region])
	if err != nil {
		return
	}
	return
}

func UpdataMetadata(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	var md metadata.MetaData
	err = json.Unmarshal(data, &md)
	if err != nil {
		tool.ReturnError(w, err)
		return
	}

	tmd, err := metadata.GetMetaDataByRegion(md.Region)
	if err != nil {
		if !tool.IsNotFound(err) {
			tool.ReturnError(w, err)
			return
		} else if err = mongo.SaveMetaData(md); err != nil {
			tool.ReturnError(w, err)
			return
		}
	} else if tmd.Sid != "" {
		if err = metadata.DelteMetaData(*tmd); err != nil {
			tool.ReturnError(w, err)
		}
		return
	}

	bridge.GetMetaChan() <- 1
	return
}
