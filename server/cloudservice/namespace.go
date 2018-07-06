package cloudservice

import (
	"github.com/andy-zhangtao/DDog/model/caasmodel"
	"strings"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/qcloud_api/v1/namespace"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"net/url"
		"github.com/andy-zhangtao/DDog/server/dbservice"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/19.

// CheckNamespace 验证指定的命名空间是否存在
// 如果不存在,则创建
func CheckNamespace(ns caasmodel.NameSpace) (err error) {
	name := strings.Replace(strings.ToLower(ns.Name), " ", "-", -1)
	_, err = dbservice.GetNamespaceByOwnerAndName(name, ns.Owner)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			var md *metadata.MetaData
			if name == "proenv" {
				//	预发布环境
				md, err = metadata.GetMetaDataByRegion("", "proenv")
				if err != nil {
					return errors.New(_const.RegionNotFound)
				}
			} else {
				md, err = metadata.GetMetaDataByRegion("")
				if err != nil {
					return errors.New(_const.RegionNotFound)
				}
			}

			q := namespace.NSpace{
				Pub: public.Public{
					Region:   md.Region,
					SecretId: md.Sid,
				},
				SecretKey: md.Skey,
				ClusterId: md.ClusterID,
				Name:      url.QueryEscape(name),
				Desc:      url.QueryEscape(ns.Desc),
			}

			q.SetDebug(true)
			if err = q.CreateNamespace(); err != nil {
				if strings.Contains(err.Error(), "NamespaceAlreadyExist") {
					return nil
				}
				return err
			}

			return dbservice.SaveNamespace(caasmodel.NameSpace{
				Name:  name,
				Owner: ns.Owner,
			})
		} else {
			return err
		}
	}

	return dbservice.UpdateNamespace(ns)
}

func DeleteNamespace(ns caasmodel.NameSpace) (err error) {
	name := strings.Replace(strings.ToLower(ns.Name), " ", "-", -1)
	oldNs, err := dbservice.GetNamespaceByOwnerAndName(name, ns.Owner)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}

		return err
	}

	var md *metadata.MetaData
	if name == "proenv" {
		//	预发布环境
		md, err = metadata.GetMetaDataByRegion("", "proenv")
		if err != nil {
			return errors.New(_const.RegionNotFound)
		}
	} else {
		md, err = metadata.GetMetaDataByRegion("")
		if err != nil {
			return errors.New(_const.RegionNotFound)
		}
	}

	q := namespace.NSpace{
		Pub: public.Public{
			Region:   md.Region,
			SecretId: md.Sid,
		},
		SecretKey: md.Skey,
		ClusterId: md.ClusterID,
		Name:      url.QueryEscape(name),
		Desc:      url.QueryEscape(ns.Desc),
		Rmname: []string{
			url.QueryEscape(name),
		},
	}

	if err = q.DeleteNamespace(); err != nil {
		return err
	}

	return dbservice.DeleteNamespaceByID(oldNs.ID)
}
