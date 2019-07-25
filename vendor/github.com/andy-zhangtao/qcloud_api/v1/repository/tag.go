package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/sirupsen/logrus"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
//QueryMyTag 查询镜像下的所有镜像
func (this *Repository) QueryMyTag(name string) (repoInfo []QCTag_data_tagInfo, err error) {
	offset := 0
	this.reponame = name
	for {
		qcr, err := this.queryMyTag(offset)
		if err != nil {
			err = errors.New(fmt.Sprintf("Query All Repostiry Error [%s]", err.Error()))
			return repoInfo, err
		}

		repoInfo = append(repoInfo, qcr.Data.TagInfo...)
		if len(repoInfo) >= qcr.Data.TagCount {
			return repoInfo, err
		}

		offset = qcr.Data.TagCount - len(repoInfo)
	}
}

//RenameMyTag 重命名指定镜像
//srcName 源镜像名称 xxxx/xx:xx
//destName 目标镜像名称 xxxx/xx:xx
func (this *Repository) RenameMyTag() (response QCRepository, err error) {
	signStr, sign := this.generatePubParam(RenameRepository)

	logrus.WithFields(logrus.Fields{"url": public.Repostiory_API_URL, "Key": this.SecretKey, "Body": signStr, "Sing": sign}).Info(ModuleName)

	resp, err := http.Get(public.Repostiory_API_URL + sign)
	if err != nil {
		err = errors.New(fmt.Sprintf("Rename Repository Error [%s]", err.Error()))
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &response); err != nil {
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s] Body [%s]", err.Error(), string(data)))
	}

	return
}

func (this *Repository) queryMyTag(offset int) (qcr QCTag, err error) {
	this.offset = strconv.Itoa(offset)
	signStr, sign := this.generatePubParam(QueryRepositoryTag)

	logrus.WithFields(logrus.Fields{"url": public.Repostiory_API_URL, "Key": this.SecretKey, "Body": signStr, "Sing": sign}).Info(ModuleName)

	resp, err := http.Get(public.Repostiory_API_URL + sign)
	if err != nil {
		err = errors.New(fmt.Sprintf("Query ALl Repository Error [%s]", err.Error()))
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &qcr); err != nil {
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s] Body [%s]", err.Error(), string(data)))
	}

	return
}

func (this *Repository) DeleteMyTag(name string) (qcr QCTagSimple, err error) {
	this.reponame = name
	signStr, sign := this.generatePubParam(DeleteRepositoryTag)

	logrus.WithFields(logrus.Fields{"url": public.Repostiory_API_URL, "Key": this.SecretKey, "Body": signStr, "Sing": sign}).Info(ModuleName)

	resp, err := http.Get(public.Repostiory_API_URL + sign)
	if err != nil {
		return qcr, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &qcr); err != nil {
		err = errors.New(fmt.Sprintf("Unmarshal Error [%s] Body [%s]", err.Error(), string(data)))
	}

	return
}
