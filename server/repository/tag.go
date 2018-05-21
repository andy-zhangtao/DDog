package repository

import (
	"github.com/andy-zhangtao/qcloud_api/v1/repository"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"fmt"
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

func RenameMyTag(srcname, destname string) (err error) {
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

	q.Params = map[string]interface{}{
		"src_image":  srcname,
		"dest_image": destname,
	}

	resp, err := q.RenameMyTag()
	if err != nil {
		err = errors.New(fmt.Sprintf("Rename Image Tag Error [%s] SrcImage [%s] DestImage [%s]", err.Error(), srcname, destname))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("Rename Image Tag Error RespCode [%d] Message[%s] SrcImage [%s] DestImage [%s]", resp.Code, resp.Message, srcname, destname))
		return
	}

	return nil
}
