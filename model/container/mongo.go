package container

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"github.com/andy-zhangtao/DDog/server/tool"
)

type Container struct {
	ID   bson.ObjectId     `json:"id,omitempty" bson:"_id,omitempty"`
	Name string            `json:"name"`
	Img  string            `json:"img"`
	Cmd  []string          `json:"cmd"`
	Env  map[string]string `json:"env"`
	Svc  string            `json:"svc"`
	Nsme string            `json:"namespace"`
	Idx  int               `json:"idx"`
}

func GetContainerByName(conname, svcname, namespace string) (con *Container, err error) {
	c, err := mongo.GetContaienrByName(conname, svcname, namespace)
	if err != nil {
		if !tool.IsNotFound(err) {
			return
		}
	}

	if tool.IsNotFound(err) {
		err = nil
		return
	}

	err = nil
	con, err = unmarshal(c)
	return
}

func SaveContainer(con *Container)(err error){
	con.ID = bson.NewObjectId()
	if err = mongo.SaveContainer(con); err != nil {
		return
	}
	return
}
func unmarshal(icon interface{}) (con *Container, err error) {
	if icon == nil {
		return
	}
	data, err := bson.Marshal(icon)
	if err != nil {
		return
	}

	var c Container
	err = bson.Unmarshal(data, &c)
	if err != nil {
		return
	}

	con = &c
	return
}
