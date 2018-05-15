package repository

import (
	"github.com/andy-zhangtao/qcloud_api/v1/repository"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
func QueryMyTag(name string) (repos []repository.QCTag_data_tagInfo, err error) {
	md, err := metadata.GetMetaDataByRegion("")
	if err != nil {
		err = errors.New(_const.RegionNotFound)
		return
	}

	q := repository.Repository{
		Pub: public.Public{
			Region:   md.Region,
			SecretId: md.Sid,
		},
		SecretKey: md.Skey,
	}

	return q.QueryMyTag(name)
}
