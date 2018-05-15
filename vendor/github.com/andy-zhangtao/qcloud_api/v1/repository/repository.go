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
func (this *Repository) QueryMyRepository() (repoInfo []QCRepository_data_repoInfo, err error) {
	offset := 0
	for {
		qcr, err := this.queryMyRepository(offset)
		if err != nil {
			err = errors.New(fmt.Sprintf("Query All Repostiry Error [%s]", err.Error()))
			return repoInfo, err
		}

		repoInfo = append(repoInfo, qcr.Data.RepoInfo...)
		if len(repoInfo) >= qcr.Data.TotalCount {
			return repoInfo, err
		}

		offset = qcr.Data.TotalCount - len(repoInfo)
	}
}

func (this *Repository) queryMyRepository(offset int) (qcr QCRepository, err error) {
	this.offset = strconv.Itoa(offset)
	signStr, sign := this.generatePubParam(QueryAllRepository)

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
