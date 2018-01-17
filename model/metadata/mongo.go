package metadata

import (
	"github.com/andy-zhangtao/DDog/server/mongo"
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/tool"
)

type MetaData struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Sid       string        `json:"secret_id"`
	Skey      string        `json:"secret_key"`
	Region    string        `json:"region"`
	ClusterID string        `json:"cluster_id"`
}

// GetMetaDataByRegion 查询指定区域的元数据
// 如果region为空， 则默认输出第一条MetaData数据
func GetMetaDataByRegion(region string) (md *MetaData, err error) {

	if region == "" {
		imds, err := mongo.FindAllMetaData()
		if err != nil {
			return nil, err
		}
		md, err = unmarshal(imds[0])
		if err != nil {
			return nil, err
		}
		return md, nil
	}

	imd, err := mongo.FindMetaDataByRegion(region)
	if err != nil {
		if tool.IsNotFound(err) {
			md = new(MetaData)
			err = nil
			return nil, err
		}
		return
	}
	md, err = unmarshal(imd)
	return
}

// DelteMetaData 删除指定区域的元数据
func DelteMetaData(md MetaData) (err error) {
	err = mongo.DeleteMetaData(md.ID.Hex())
	return
}

func DeleteMetaDataByRegion(region string) (err error) {
	md, err := GetMetaDataByRegion(region)
	if err != nil {
		return
	}

	if md.ID == "" {
		return
	}

	err = DelteMetaData(*md)
	return
}

func SaveMetaData(md MetaData) error {
	return mongo.SaveMetaData(md)
}

func unmarshal(imd interface{}) (md *MetaData, err error) {
	if imd == nil {
		return
	}
	data, err := bson.Marshal(imd)
	if err != nil {
		return
	}

	var m MetaData
	err = bson.Unmarshal(data, &m)
	if err != nil {
		return
	}

	md = &m
	return
}
