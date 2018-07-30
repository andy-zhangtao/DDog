package main

import (
	"github.com/andy-zhangtao/_hulk_client"
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
	"github.com/andy-zhangtao/DDog/server/repository"
	"github.com/nsqio/go-nsq"
	"github.com/andy-zhangtao/DDog/server/tool"
	"github.com/openzipkin/zipkin-go"
	"time"
)

const ModuleName = "DDog-Server-GraphQL"
const ModuleVersion = "v0.1.0"
const ModuleResume = "Caas Graphql接口平台"

var producer *nsq.Producer

func init() {
	_hulk_client.Run()
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

	producer, _ = nsq.NewProducer(os.Getenv(_const.EnvNsqdEndpoint), nsq.NewConfig())
}

//使用GraphQL接口的DDog Server
func main() {

	router := mux.NewRouter()
	router.Path("/api").HandlerFunc(handleGraphQL)
	router.Path("/backup/{filename}").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		filename := mux.Vars(request)["filename"]

		http.Redirect(writer, request, "/download/"+filename, http.StatusMovedPermanently)
	})
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
		"repository": &graphql.Field{
			Type:        graphql.NewList(caas.RepositoryType),
			Description: "Kubentes Repository Data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return repository.QueryMyRepository()
			},
		},
		"tag": &graphql.Field{
			Type:        graphql.NewList(caas.TagType),
			Description: "Kubentes Tag Data",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				return repository.QueryMyTag(name)
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
				"traceid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"parentid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				owner, _ := p.Args["owner"].(string)
				name, _ := p.Args["name"].(string)
				desc, _ := p.Args["desc"].(string)

				//zipkin parameters
				traceid, _ := p.Args["traceid"].(string)
				id, _ := p.Args["id"].(string)
				parentid, _ := p.Args["parentid"].(string)

				var errmessage = ""
				if traceid != "" && id != "" {

					span, reporter, iscreate := tool.GetZipKinSpan(ModuleName, "AddNameSpace", traceid, id, parentid)
					if iscreate {
						defer func() {
							if errmessage != "" {
								span.Annotate(time.Now(), fmt.Sprintf("%s-AddNamespace Receive Check Namespace Error [%s]", ModuleName, errmessage))
							}

							span.Finish()
							reporter.Close()
						}()
					}
					//zipKinUrl := os.Getenv(_const.ENV_AGENT_ZIPKIN_ENDPOINT)
					//if zipKinUrl != "" {
					//	reporter := httpreporter.NewReporter(fmt.Sprintf("%s/api/v2/spans", zipKinUrl))
					//	endpoint, err := zipkin.NewEndpoint(ModuleName, tool.GetLocalIP())
					//	if err != nil {
					//		logrus.WithFields(logrus.Fields{"Create ZipKin Endpoint Error": fmt.Sprintf("unable to create local endpoint: %+v\n", err)}).Error(ModuleName)
					//	} else {
					//		tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
					//		if err != nil {
					//			logrus.WithFields(logrus.Fields{"Create ZipKin Tracer Error": fmt.Sprintf("unable to create tracer: %+v\n", err)}).Error(ModuleName)
					//		} else {
					//			//为了还原成正确的traceid,需要在traceid前后各添加一个0
					//			//具体原因，参考traceid UnmarshalJSON源码
					//			traceid = fmt.Sprintf("0%s0", traceid)
					//			id = fmt.Sprintf("0%s0", id)
					//			ctx := zmodel.SpanContext{}
					//			_tracid := new(zmodel.TraceID)
					//			_tracid.UnmarshalJSON([]byte(traceid))
					//			ctx.TraceID = *_tracid
					//
					//			_id := new(zmodel.ID)
					//			_id.UnmarshalJSON([]byte(id))
					//			ctx.ID = *_id
					//
					//			_parentid := new(zmodel.ID)
					//			_parentid.UnmarshalJSON([]byte(parentid))
					//			ctx.ParentID = _parentid
					//
					//			span := tracer.StartSpan(ModuleName, zipkin.Parent(ctx))
					//			logrus.WithFields(logrus.Fields{"span": span}).Info(ModuleName)
					//			span.Annotate(time.Now(), fmt.Sprintf("%s-AddNamespace Receive Check Namespace Request", ModuleName))
					//
					//		}
					//	}
					//}
				}

				ns := caasmodel.NameSpace{
					Name:  name,
					Owner: owner,
					Desc:  desc,
				}

				if name == "proenv" {
					ns.Owner = "admin"
				}

				if name == "release" {
					ns.Owner = "admin"
				}
				err := cloudservice.CheckNamespace(ns)
				if err != nil {
					errmessage = err.Error()
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

				if name == "proenv" {
					ns.Owner = "admin"
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
				"traceid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"parentid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},

			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)
				instance := 2

				//zipkin parameters
				traceid, _ := p.Args["traceid"].(string)
				id, _ := p.Args["id"].(string)
				parentid, _ := p.Args["parentid"].(string)

				var errmessage = ""
				var span zipkin.Span
				if traceid != "" && id != "" {
					s, reporter, iscreate := tool.GetZipKinSpan(ModuleName, "StartDeploy", traceid, id, parentid)
					if iscreate {
						span = s
						defer func() {
							if errmessage != "" {
								span.Annotate(time.Now(), fmt.Sprintf("%s-StartDeploy Error [%s]", ModuleName, errmessage))
							}

							span.Finish()
							reporter.Close()
						}()
					}
				}

				if i, ok := p.Args["instance"]; ok {
					instance, _ = i.(int)
				}

				cf, err := svcconf.GetSvcConfByName(name, namespace)
				if err != nil {
					errmessage = err.Error()
					return nil, err
				}

				if cf == nil {
					errmessage = err.Error()
					return nil, errors.New(_const.SVCNoExist)
				}

				msg := agent.DeployMsg{
					Span:      span.Context(),
					SvcName:   name,
					NameSpace: namespace,
					Upgrade:   true,
					Replicas:  instance,
				}

				cf.Deploy = _const.DeployIng
				err = svcconf.UpdateSvcConf(cf)
				if err != nil {
					errmessage = err.Error()
					return nil, err
				}

				data, err := json.Marshal(msg)
				if err != nil {
					errmessage = err.Error()
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
				"replicas": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"traceid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"parentid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)
				replicas, _ := p.Args["replicas"].(int)

				//zipkin parameters
				traceid, _ := p.Args["traceid"].(string)
				id, _ := p.Args["id"].(string)
				parentid, _ := p.Args["parentid"].(string)

				var errmessage = ""
				if traceid != "" && id != "" {
					span, reporter, iscreate := tool.GetZipKinSpan(ModuleName, "AddService", traceid, id, parentid)
					if iscreate {
						defer func() {
							if errmessage != "" {
								span.Annotate(time.Now(), fmt.Sprintf("%s-AddService Error [%s]", ModuleName, errmessage))
							}

							span.Finish()
							reporter.Close()
						}()
					}
				}

				if replicas == 0 {
					replicas = 1
				}

				conf := &svcconf.SvcConf{
					Name:      name,
					Namespace: namespace,
					Replicas:  replicas,
				}

				cf, err := svcconf.GetSvcConfByName(conf.Name, conf.Namespace)
				if err != nil {
					if strings.Contains(err.Error(), "not found") {
						conf.Id = bson.NewObjectId()
						if err = mongo.SaveSvcConfig(conf); err != nil {
							errmessage = err.Error()
							return *conf, err
						}
					}
					return nil, err
				}
				//这里的判断不优雅，需要改掉
				if cf == nil {
					conf.Id = bson.NewObjectId()
					if err = mongo.SaveSvcConfig(conf); err != nil {
						errmessage = err.Error()
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
				"env": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"traceid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"parentid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["service"].(string)
				image, _ := p.Args["image"].(string)
				namespace, _ := p.Args["namespace"].(string)

				//zipkin parameters
				traceid, _ := p.Args["traceid"].(string)
				id, _ := p.Args["id"].(string)
				parentid, _ := p.Args["parentid"].(string)

				var errmessage = ""
				if traceid != "" && id != "" {
					span, reporter, iscreate := tool.GetZipKinSpan(ModuleName, "AddContainer", traceid, id, parentid)
					if iscreate {
						defer func() {
							if errmessage != "" {
								span.Annotate(time.Now(), fmt.Sprintf("%s-AddContainer Error [%s]", ModuleName, errmessage))
							}

							span.Finish()
							reporter.Close()
						}()
					}
				}

				var ps []int
				envs := make(map[string]string)

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

				if e, ok := p.Args["env"]; ok {
					et, _ := e.([]interface{})
					for _, ett := range et {
						if str, ok := ett.(string); ok {
							tm := strings.Split(str, "=")
							envs[strings.TrimSpace(tm[0])] = strings.TrimSpace(strings.Join(tm[1:], "="))
						}
					}
				}

				con := container.Container{
					Name: name,
					Img:  image,
					Port: ps,
					Svc:  name,
					Nsme: namespace,
					Env:  envs,
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
					errmessage = err.Error()
					return nil, err
				}

				logrus.WithFields(logrus.Fields{"find old container": isExist}).Info(ModuleName)
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
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				endpoint, _ := p.Args["endpoint"].(string)
				token, _ := p.Args["token"].(string)
				region, _ := p.Args["region"].(string)
				namespace, _ := p.Args["namespace"].(string)

				if err := k8service.UpdateK8sCluster(k8sconfig.K8sCluster{
					Name:      name,
					Region:    region,
					Endpoint:  endpoint,
					Token:     token,
					Namespace: namespace,
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
		"renametag": &graphql.Field{
			Type:        graphql.String,
			Description: "Rename Image Tag",
			Args: graphql.FieldConfigArgument{
				"src": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"dest": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				src, _ := p.Args["src"].(string)
				dest, _ := p.Args["dest"].(string)

				if err := repository.RenameMyTag(src, dest); err != nil {
					return err.Error(), nil
				}
				return nil, nil
			},
		},
		"replica": &graphql.Field{
			Type:        graphql.Int,
			Description: "Modify the service instance number",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"namespace": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"replica": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				name, _ := p.Args["name"].(string)
				namespace, _ := p.Args["namespace"].(string)
				replica, _ := p.Args["replica"].(int)

				if err := qcloud.ModifyInstancesReplica(name, namespace, replica); err != nil {
					return nil, err
				}

				data, _ := json.Marshal(agent.DeployMsg{
					SvcName:   name,
					NameSpace: namespace,
					Replicas:  replica,
				})

				return replica, producer.Publish(_const.SvcReplicaMsg, data)
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
