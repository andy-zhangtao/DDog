package caas

import (
	"github.com/graphql-go/graphql"
	"github.com/andy-zhangtao/DDog/model/caasmodel"
	"github.com/andy-zhangtao/DDog/model/svcconf"
	"github.com/andy-zhangtao/DDog/model/container"
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
