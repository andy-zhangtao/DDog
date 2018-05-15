package repository

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/sirupsen/logrus"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
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
