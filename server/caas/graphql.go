package caas

import (
	"github.com/graphql-go/graphql"
	"github.com/andy-zhangtao/DDog/model/caasmodel"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/model/container"
	"github.com/andy-zhangtao/qcloud_api/v1/event"
	"github.com/andy-zhangtao/DDog/model/metadata"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/qcloud_api/v1/public"
	"fmt"
	"github.com/andy-zhangtao/qcloud_api/v1/service"
	"github.com/andy-zhangtao/DDog/model/k8sconfig"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/4/18.
//提供GraphQL接口

var CaasServiceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "caas",
	Fields: graphql.Fields{
		"namespace": &graphql.Field{
			Type: graphql.NewList(CaasNameSpaceType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var ns []caasmodel.NameSpace
				if n, ok := p.Source.([]interface{}); ok {
					for _, nt := range n {
						if nn, ok := nt.(caasmodel.NameSpace); ok {
							ns = append(ns, nn)
						}
					}
				}
				return ns, nil
			},
		},
	},
})

var CaasServiceConfType = graphql.NewObject(graphql.ObjectConfig{
	Name: "service",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					return s.Id.Hex(), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					return s.Name, nil
				}
				return nil, nil
			},
		},
		"svc_name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					return s.SvcName, nil
				}
				return nil, nil
			},
		},
		"namespace": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					return s.Namespace, nil
				}
				return nil, nil
			},
		},
		"deploy": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					return s.Deploy, nil
				}
				return nil, nil
			},
		},
		"loadbalance": &graphql.Field{
			Type: graphql.NewList(graphql.String),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					var lb []string
					for _, port := range s.LbConfig.Port {
						lb = append(lb, fmt.Sprintf("%s:%d", s.LbConfig.IP, port))
					}
					return lb, nil
				}
				return nil, nil
			},
		},
		"events": &graphql.Field{
			Type: graphql.NewList(CaasServiceEventType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if s, ok := p.Source.(svcconf.SvcConf); ok {
					md, err := metadata.GetMetaDataByRegion("")
					if err != nil {
						return nil, errors.New(_const.RegionNotFound)
					}

					ser := event.ServiceEventRequest{
						Svcname:   s.SvcName,
						Namespace: s.Namespace,
						ClusterId: md.ClusterID,
						SecretKey: md.Skey,
						Pub: public.Public{
							Region:   md.Region,
							SecretId: md.Sid,
						},
						Debug: true,
					}

					return ser.GetServiceEvent()
				}

				return nil, nil
			},
		},
	},
})

var CaasServiceEventType = graphql.NewObject(graphql.ObjectConfig{
	Name: "events",
	Fields: graphql.Fields{
		"count": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(event.SEvent); ok {
					return e.Count, nil
				}

				return nil, nil
			},
		},
		"level": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(event.SEvent); ok {
					return e.Level, nil
				}

				return nil, nil
			},
		},
		"reason": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(event.SEvent); ok {
					return e.Reason, nil
				}

				return nil, nil
			},
		},
		"message": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(event.SEvent); ok {
					return e.Message, nil
				}

				return nil, nil
			},
		},
	},
})

var CaasNameSpaceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "namespace",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if n, ok := p.Source.(caasmodel.NameSpace); ok {
					return n.ID.Hex(), nil
				}
				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if n, ok := p.Source.(caasmodel.NameSpace); ok {
					return n.Name, nil
				}
				return nil, nil
			},
		},
		"owner": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if n, ok := p.Source.(caasmodel.NameSpace); ok {
					return n.Owner, nil
				}
				return nil, nil
			},
		},
		"desc": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if n, ok := p.Source.(caasmodel.NameSpace); ok {
					return n.Desc, nil
				}
				return nil, nil
			},
		},
		"creatime": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if n, ok := p.Source.(caasmodel.NameSpace); ok {
					return n.CreateTime, nil
				}
				return nil, nil
			},
		},
	},
})

var CaasContainerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "container",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(container.Container); ok {
						return c.ID.Hex(), nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(container.Container); ok {
						return c.Name, nil
					}
					return nil, nil
				},
			},
			"image": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(container.Container); ok {
						return c.Img, nil
					}
					return nil, nil
				},
			},
			"service": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(container.Container); ok {
						return c.Svc, nil
					}
					return nil, nil
				},
			},
			"ports": &graphql.Field{
				Type: graphql.NewList(graphql.Int),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if c, ok := p.Source.(container.Container); ok {
						return c.Port, nil
					}
					return nil, nil
				},
			},
		},
	},
)

var InstanceType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "instance",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if i, ok := p.Source.(service.Instance); ok {
						return i.Name, nil
					}
					return nil, nil
				},
			},
		},
	},
)

var K8sClusterTYpe = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "k8s",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.ID.Hex(), nil
					}
					return nil, nil
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.Name, nil
					}
					return nil, nil
				},
			},
			"region": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.Region, nil
					}
					return nil, nil
				},
			},
			"endpoint": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.Endpoint, nil
					}
					return nil, nil
				},
			},
			"token": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.Token, nil
					}
					return nil, nil
				},
			},
			"update": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if k, ok := p.Source.(k8sconfig.K8sCluster); ok {
						return k.UpdateTime, nil
					}
					return nil, nil
				},
			},
		},
	},
)
