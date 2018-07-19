package _hulk_client

import (
	"github.com/shurcooL/graphql"
	"os"
	"context"
	"github.com/sirupsen/logrus"
	"fmt"
	"encoding/json"
)

func queryHulk(name, version string) (env map[string]interface{}, err error) {

	variables := map[string]interface{}{
		"name":    graphql.String(name),
		"version": graphql.String(version),
	}

	var query struct {
		QueryHulk []struct {
			Name      graphql.String
			Version   graphql.String
			Configure graphql.String
		} `graphql:"queryHulk(name: $name, version: $version)"`
	}

	logrus.WithFields(logrus.Fields{"variables": variables}).Info(HULK_GO_SDK)
	client := graphql.NewClient(os.Getenv(ENDPOINT), nil)
	err = client.Query(context.Background(), &query, variables)
	if err != nil {
		logrus.Error(fmt.Sprintf("Query Hulk Error [%s]", err))
	}

	logrus.WithFields(logrus.Fields{"hulk": query.QueryHulk}).Info(HULK_GO_SDK)

	env = make(map[string]interface{})
	err = json.Unmarshal([]byte(query.QueryHulk[0].Configure), &env)
	if err != nil {
		logrus.Error(fmt.Sprintf("Unmarshal Configure Error [%s]", err))
	}

	logrus.WithFields(env).Info(HULK_GO_SDK)

	return
}

func padding(env map[string]interface{}) {
	for key, value := range env {
		os.Setenv(key, value.(string))
	}

}
