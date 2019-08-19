package repository

import (
	"errors"
	"fmt"

	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/v1/repository"
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

func RmMyTag(name, tag string) (err error) {
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
		"tags.0": tag,
	}

	resp, err := q.DeleteMyTag(name)
	if err != nil {
		err = errors.New(fmt.Sprintf("Remove Image Tag Error [%s] SrcImage [%s] ", err.Error(), name, tag))
		return
	}

	if resp.Code != 0 {
		err = errors.New(fmt.Sprintf("Remove Image Tag Error StatusCode[%d] Error: [%s] ", resp.Code, resp.CodeDesc))
		return
	}

	return nil
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
		err = errors.New(fmt.Sprintf("Rename Image Tag Error RespCode [%d] Message[%s] SrcImage [%s] DestImage [%s]", resp.Code, resp.Message+resp.CodeDesc, srcname, destname))
		return
	}

	return nil
}
