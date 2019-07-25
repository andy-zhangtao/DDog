package repository

//Write by zhangtao<ztao8607@gmail.com> . In 2018/5/15.
type QCRepository_data_repoInfo struct {
	Reponame         string `json:"reponame"`
	Repotype         string `json:"repotype"`
	TagCount         int    `json:"tagCount"`
	Public           int    `json:"public"`
	IsUserFavor      bool   `json:"isUserFavor"`
	IsQcloudOfficial bool   `json:"isQcloudOfficial"`
	FavorCount       int    `json:"favorCount"`
	PullCount        int    `json:"pullCount"`
	Description      string `json:"description"`
	CreationTime     string `json:"creationTime"`
	UpdateTime       string `json:"updateTime"`
}
type QCRepository_data struct {
	PrivilegeFiltered bool                         `json:"privilegeFiltered"`
	RepoInfo          []QCRepository_data_repoInfo `json:"repoInfo"`
	Server            string                       `json:"server"`
	TotalCount        int                          `json:"totalCount"`
}
type QCRepository struct {
	Code     int               `json:"code"`
	Message  string            `json:"message"`
	CodeDesc string            `json:"codeDesc"`
	Data     QCRepository_data `json:"data"`
}

type QCTag_data_tagInfo struct {
	Id            int    `json:"id"`
	RepoName      string `json:"repo_name"`
	TagName       string `json:"tagName"`
	TagId         string `json:"tagId"`
	ImageId       string `json:"imageId"`
	Size          string `json:"size"`
	CreationTime  string `json:"creationTime"`
	UpdateTime    string `json:"updateTime"`
	Author        string `json:"author"`
	Architecture  string `json:"architecture"`
	DockerVersion string `json:"dockerVersion"`
	Os            string `json:"os"`
	PushTime      string `json:"pushTime"`
	SizeByte      int    `json:"sizeByte"`
}
type QCTag_data struct {
	Reponame string               `json:"reponame"`
	Server   string               `json:"server"`
	TagCount int                  `json:"tagCount"`
	TagInfo  []QCTag_data_tagInfo `json:"tagInfo"`
}
type QCTag struct {
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	CodeDesc string     `json:"codeDesc"`
	Data     QCTag_data `json:"data"`
}

type QCTagSimple struct {
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	CodeDesc string     `json:"codeDesc"`
}