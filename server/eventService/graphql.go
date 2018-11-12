/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package eventService

import (
	"github.com/andy-zhangtao/DDog/model/events"
	"github.com/andy-zhangtao/DDog/server/svcconf"
	"github.com/graphql-go/graphql"
)

var WatchEventType = graphql.NewObject(graphql.ObjectConfig{
	Name: "watchEvents",
	Fields: graphql.Fields{
		"time": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(events.K8sWatchEvent); ok {
					return e.Time, nil
				}

				return nil, nil
			},
		},
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(events.K8sWatchEvent); ok {
					return e.Name, nil
				}

				return nil, nil
			},
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(events.K8sWatchEvent); ok {
					return e.Type, nil
				}

				return nil, nil
			},
		},
		"message": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if e, ok := p.Source.(events.K8sWatchEvent); ok {
					return e.Message, nil
				}

				return nil, nil
			},
		},
	},
})

var WatchEvent = &graphql.Field{
	Type:        graphql.NewList(WatchEventType),
	Description: "Query Speicfy Service Events",
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"namespace": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"desc": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		name, _ := p.Args["name"].(string)
		namespace, _ := p.Args["namespace"].(string)
		desc, _ := p.Args["desc"].(bool)

		scf, err := svcconf.GetSvcConfByName(name, namespace)
		if err != nil {
			return nil, err
		}

		events, err := QueryServiceEvents(name, scf.SvcName, namespace, desc)
		if err != nil {
			return nil, err
		}

		return events, nil
	},
}
