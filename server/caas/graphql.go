package caas

import (
	"github.com/graphql-go/graphql"
	"github.com/andy-zhangtao/DDog/model/caasmodel"
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
