package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"encoding/json"
	"io/ioutil"
	"github.com/graphql-go/graphql"
	"fmt"
	"github.com/andy-zhangtao/DDog/check"
	"os"
	"github.com/andy-zhangtao/DDog/server/caas"
	"github.com/andy-zhangtao/DDog/server/dbservice"
	"github.com/andy-zhangtao/DDog/model/caasmodel"

	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/server/cloudservice"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/agent"
	"github.com/andy-zhangtao/DDog/bridge"
	"github.com/andy-zhangtao/DDog/server/mongo"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"github.com/andy-zhangtao/DDog/model/container"
	cs "github.com/andy-zhangtao/DDog/server/container"
)

const ModuleName = "DDog-Server-GraphQL"

func init() {

	if err := check.CheckMongo(); err != nil {
		logrus.WithFields(logrus.Fields{"Check Mongo Error": err}).Error(ModuleName)
		os.Exit(-1)
	}

	if err := check.CheckNamespace(); err != nil {
		logrus.WithFields(logrus.Fields{"Check Env Error": err}).Error(ModuleName)
		os.Exit(-1)
	}

	if err := check.CheckNsq(); err != nil {
		logrus.WithFields(logrus.Fields{"Check Nsq Error": err}).Error(ModuleName)
		os.Exit(-1)
	}

	if err := check.CheckLogOpt(); err != nil {
		logrus.WithFields(logrus.Fields{"Check Log Opt Error": err}).Error(ModuleName)
		os.Exit(-1)
	}
}

//使用GraphQL接口的DDog Server
func main() {
	router := mux.NewRouter()
	router.Path("/api").HandlerFunc(handleGraphQL)
	handler := cors.AllowAll().Handler(router)
	logrus.Fatal(http.ListenAndServe(":8000", handler))
}

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		//返回所有命名空间数据
		"namespace": &graphql.Field{
			Type:        graphql.NewList(caas.CaasNameSpaceType),
			Description: "The NameSpace In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				if name == "" {
					ns, err := dbservice.GetAllNamesapce()
					if err != nil {
						return nil, err
					}
					return ns, err
				}

				ns, err := dbservice.GetNamespaceByName(name)
				if err != nil {
					return nil, err
				}

				return []caasmodel.NameSpace{ns}, nil
			},
		},
		//返回指定命名空间的服务信息
		"service": &graphql.Field{
			Type:        graphql.NewList(caas.CaasServiceConfType),
			Description: "The Service Info Which Deploy In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"deploy": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)

				if name != "" {
					scf, err := svcconf.GetSvcConfByName(name, namespace)
					if err != nil {
						return nil, err
					}

					return *scf, nil
				}

				if d, ok := p.Args["deploy"]; ok {
					deploy, _ := d.(int)
					scf, err := svcconf.GetSvcConfByDeployStatus(deploy)
					if err != nil {
						return nil, err
					}

					return *scf, nil
				}

				if namespace != "" {
					scf, err := svcconf.GetSvcConfByNamespace(namespace)
					if err != nil {
						return nil, err
					}

					return scf, nil
				}

				return nil, nil
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		//对命名空间的新增操作,操作幂等
		"addNamespace": &graphql.Field{
			Type:        caas.CaasNameSpaceType,
			Description: "Create A New Namespace",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"owner": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"desc": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				owner, _ := p.Args["owner"].(string)
				name, _ := p.Args["name"].(string)
				desc, _ := p.Args["desc"].(string)

				ns := caasmodel.NameSpace{
					Name:  name,
					Owner: owner,
					Desc:  desc,
				}

				err := cloudservice.CheckNamespace(ns)
				if err != nil {
					return nil, err
				}

				return ns, nil
			},
		},
		"delNamespace": &graphql.Field{
			Type:        caas.CaasNameSpaceType,
			Description: "Delete A Exist Namespace In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"owner": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				owner, _ := p.Args["owner"].(string)
				name, _ := p.Args["name"].(string)

				ns := caasmodel.NameSpace{
					Name:  name,
					Owner: owner,
				}

				err := cloudservice.DeleteNamespace(ns)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
		"deployService": &graphql.Field{
			Type:        caas.CaasServiceConfType,
			Description: "Deploy A New Service In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"instance": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)
				instance := 2

				if i, ok := p.Args["instance"]; ok {
					instance, _ = i.(int)
				}

				cf, err := svcconf.GetSvcConfByName(name, namespace)
				if err != nil {
					return nil, err
				}

				if cf == nil {
					return nil, errors.New(_const.SVCNoExist)
				}

				msg := agent.DeployMsg{
					SvcName:   name,
					NameSpace: namespace,
					Upgrade:   true,
					Replicas:  instance,
				}

				cf.Deploy = _const.DeployIng
				err = svcconf.UpdateSvcConf(cf)
				if err != nil {
					return nil, err
				}

				data, err := json.Marshal(msg)
				if err != nil {
					return nil, err
				}

				bridge.SendDeployMsg(string(data))

				return *cf, nil
			},
		},
		"addService": &graphql.Field{
			Type:        caas.CaasServiceConfType,
			Description: "Add A New Service In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)

				conf := &svcconf.SvcConf{
					Name:      name,
					Namespace: namespace,
					Replicas:  2,
				}

				cf, err := svcconf.GetSvcConfByName(conf.Name, conf.Namespace)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						conf.Id = bson.NewObjectId()
						if err = mongo.SaveSvcConfig(conf); err != nil {
							return *conf, err
						}
					}
					return nil, err
				}

				return *cf, nil
			},
		},
		"delService": &graphql.Field{
			Type:        caas.CaasServiceConfType,
			Description: "Add A New Service In Caas",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)

				conf := &svcconf.SvcConf{
					Name:      name,
					Namespace: namespace,
					Replicas:  2,
				}

				cf, err := svcconf.GetSvcConfByName(conf.Name, conf.Namespace)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						return nil, nil
					}
					return nil, err
				}

				return nil, cf.DeleteMySelf()
			},
		},
		"addContainer": &graphql.Field{
			Type:        caas.CaasContainerType,
			Description: "Add A Container Configure Under Specify Service",
			Args: graphql.FieldConfigArgument{
				"service": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"ports": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["service"].(string)
				image, _ := p.Args["image"].(string)
				namespace, _ := p.Args["namespace"].(string)

				var ps []int

				if pt, ok := p.Args["ports"]; ok {
					ptt, _ := pt.([]interface{})
					for _, pi := range ptt {
						if vp, ok := pi.(int); ok {
							ps = append(ps, vp)
						}
					}
				}
				con := container.Container{
					Name: name,
					Img:  image,
					Port: ps,
					Svc:  name,
					Nsme: namespace,
				}

				var nt []container.NetConfigure
				for _, p := range con.Port {
					nt = append(nt, container.NetConfigure{
						AccessType: 0,
						InPort:     p,
						OutPort:    p,
						Protocol:   0,
					})
				}

				con.Net = nt

				oldCon, isExist, err := cs.IsExistContainer(&con)
				if err != nil {
					return nil, err
				}

				if isExist {
					if len(con.Port) > 0 {
						oldCon.Port = con.Port
						var nt []container.NetConfigure
						for _, p := range oldCon.Port {
							nt = append(nt, container.NetConfigure{
								AccessType: 0,
								InPort:     p,
								OutPort:    p,
								Protocol:   0,
							})
						}
						oldCon.Net = nt

						return *oldCon, container.UpgradeContaienrByName(oldCon)
					}
					return *oldCon, nil
				}

				return con, container.SaveContainer(&con)
			},
		},
		"delContainer": &graphql.Field{
			Type:        caas.CaasContainerType,
			Description: "Add A Container Configure Under Specify Service",
			Args: graphql.FieldConfigArgument{
				"service": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["service"].(string)
				image, _ := p.Args["image"].(string)
				namespace, _ := p.Args["namespace"].(string)
				con := container.Container{
					Name: name,
					Img:  image,
					Svc:  name,
					Nsme: namespace,
				}

				return nil, container.DeleteContainerByName(con.Name, con.Svc, con.Nsme)
			},
		},
	},
})
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func executeQuery(query map[string]interface{}, schema graphql.Schema) *graphql.Result {

	params := graphql.Params{
		Schema:        schema,
		RequestString: query["query"].(string),
	}

	if query["variables"] != nil {
		params.VariableValues = query["variables"].(map[string]interface{})
	}

	result := graphql.Do(params)

	if len(result.Errors) > 0 {
		fmt.Println("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func handleGraphQL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var g map[string]interface{}
	if r.Method == http.MethodGet {
		g = make(map[string]interface{})
		g["query"] = r.URL.Query().Get("query")
		result := executeQuery(g, schema)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(ModuleName)
		json.NewEncoder(w).Encode(result)
	}

	if r.Method == http.MethodPost {
		data, _ := ioutil.ReadAll(r.Body)
		logrus.WithFields(logrus.Fields{"body": string(data)}).Info(ModuleName)

		err := json.Unmarshal(data, &g)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
		}
		logrus.WithFields(logrus.Fields{"graph": g}).Info(ModuleName)
		result := executeQuery(g, schema)
		logrus.WithFields(logrus.Fields{"result": result.Data}).Info(ModuleName)
		json.NewEncoder(w).Encode(result)
	}
}
