package repository

import (
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"github.com/andy-zhangtao/qcloud_api/const/v1"
	"net/url"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
type Repository struct {
	Pub       public.Public `json:"pub"`
	SecretKey string
	sign      string
	offset    string
	reponame  string
	Params    map[string]interface{}
}

const (
	ModuleName         = "Qcloud-Repository-Agent"
	QueryAllRepository = iota
	QueryRepositoryTag
	RenameRepository
)

func (this *Repository) generatePubParam(kind int) (string, string) {
	var field []string
	reqmap := make(map[string]string)

	optKind := ""
	switch kind {
	case QueryAllRepository:
		optKind = "SearchUserRepository"
		field = append(field, []string{"offset", "limit"}...)
		reqmap["offset"] = this.offset
		reqmap["limit"] = "100"
	case QueryRepositoryTag:
		optKind = "GetTagList"
		field = append(field, []string{"offset", "limit", "reponame"}...)
		reqmap["offset"] = this.offset
		reqmap["limit"] = "100"
		reqmap["reponame"] = this.reponame
	case RenameRepository:
		optKind = "DuplicateImage"
		field = append(field, []string{"src.image", "dest.image"}...)
		reqmap["src.image"] = this.Params["src_image"].(string)
		reqmap["dest.image"] = this.Params["dest_image"].(string)

	}

	pubMap := public.PublicParam(optKind, this.Pub.Region, this.Pub.SecretId)
	this.sign = public.GenerateSignatureString(field, reqmap, pubMap)
	signStr := "GET" + v1.QCloudRepositoryEndpoint + this.sign
	sign := public.GenerateSignature(this.SecretKey, signStr)
	return signStr, this.sign + "&Signature=" + url.QueryEscape(sign)
}
