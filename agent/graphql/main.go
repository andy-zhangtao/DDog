package main

import (
	"github.com/sirupsen/logrus"
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
	"github.com/andy-zhangtao/DDog/server/qcloud"
	"github.com/andy-zhangtao/DDog/server/k8service"
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
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
	router.Handle("/download/{filename}", http.StripPrefix("/download/", http.FileServer(http.Dir("/tmp"))))
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
						err = errors.New(fmt.Sprintf("Get SvcConf Error [%s] Filter name:[%s] namespace:[%s]", err.Error(), name, namespace))
						return nil, err
					}

					if scf != nil {
						return []svcconf.SvcConf{*scf}, nil
					}

					return nil, nil
				}

				if d, ok := p.Args["deploy"]; ok {
					deploy, _ := d.(int)
					scf, err := svcconf.GetSvcConfByDeployStatus(deploy)
					if err != nil {
						return nil, err
					}

					return []svcconf.SvcConf{*scf}, nil
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
		//	返回服务实例信息
		"instance": &graphql.Field{
			Type:        graphql.NewList(caas.InstanceType),
			Description: "All instances info of service",
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

				return qcloud.GetInstanceInfo(name, namespace)
			},
		},
		"k8scluster": &graphql.Field{
			Type:        graphql.NewList(caas.K8sClusterTYpe),
			Description: "Kubentes Cluster Data",
			Args: graphql.FieldConfigArgument{
				"region": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				region, _ := p.Args["region"].(string)

				if region == "" {
					return k8service.GetALlK8sCluster()
				}

				return k8service.GetK8sCluster(region)
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
				//这里的判断不优雅，需要改掉
				if cf == nil {
					conf.Id = bson.NewObjectId()
					if err = mongo.SaveSvcConfig(conf); err != nil {
						return nil, err
					}
					return *conf, nil
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
							if vp > 0 {
								ps = append(ps, vp)
							}
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
					container.DeleteContainerByName(oldCon.Name, oldCon.Svc, oldCon.Nsme)
					return *oldCon, cs.CreateContainerForGraphQL(&con)
				}

				return con, cs.CreateContainerForGraphQL(&con)
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
		"k8scluster": &graphql.Field{
			Type:        caas.K8sClusterTYpe,
			Description: "Add A K8s Cluster Data",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"endpoint": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"region": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				endpoint, _ := p.Args["endpoint"].(string)
				token, _ := p.Args["token"].(string)
				region, _ := p.Args["region"].(string)

				if err := k8service.UpdateK8sCluster(k8sconfig.K8sCluster{
					Name:     name,
					Region:   region,
					Endpoint: endpoint,
					Token:    token,
				}); err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
		"delK8sCluster": &graphql.Field{
			Type:        caas.K8sClusterTYpe,
			Description: "Delete A K8s Cluster Data",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(string)
				if err := k8service.DeleteK8sClusterByID(id); err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
		"k8sbackup": &graphql.Field{
			Type:        graphql.String,
			Description: "Kubernetes Backup",
			Args: graphql.FieldConfigArgument{
				"region": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				region, _ := p.Args["region"].(string)
				namespace, _ := p.Args["namespace"].(string)
				return k8service.BackupK8sCluster(region, namespace)
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
