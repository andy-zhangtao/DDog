/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package eventService

import (
	"encoding/json"
	"fmt"
	"github.com/andy-zhangtao/DDog/const"
	"github.com/andy-zhangtao/DDog/model/events"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/sirupsen/logrus"
)

const ModelName = "EventService"

func QueryServiceEvents(name, service, namespace string, desc bool) (event []events.K8sWatchEvent, err error) {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: os.Getenv(_const.Env_AGENT_INFLUX_ENDPOINT),
	})

	if err != nil {
		logrus.WithFields(logrus.Fields{"Error creating InfluxDB Client": err, _const.Env_AGENT_INFLUX_ENDPOINT: os.Getenv(_const.Env_AGENT_INFLUX_ENDPOINT)}).Error(ModelName)
		return event, err
	}

	defer c.Close()

	var cmd string

	if desc {
		cmd = fmt.Sprintf("SELECT \"event\", \"message\", \"type\" FROM events WHERE \"name\"='%s' and \"namespace\"='%s' and \"service\"='%s' order by time desc", name, namespace, service)
	} else {
		cmd = fmt.Sprintf("SELECT \"event\", \"message\", \"type\" FROM events WHERE \"name\"='%s' and \"namespace\"='%s' and \"service\"='%s' order by time  asc", name, namespace, service)
	}

	logrus.WithFields(logrus.Fields{"influx cmd": cmd}).Info(ModelName)
	q := client.NewQuery(cmd, os.Getenv(_const.ENV_AGENT_INFLUX_DB), "ns")
	if response, err := c.Query(q); err == nil && response.Error() == nil {

		for _, r := range response.Results {
			for _, s := range r.Series {
				for _, v := range s.Values {
					i, err := strconv.ParseInt(string(v[0].(json.Number)), 10, 64)
					if err != nil {
						return event, err
					}
					tm := time.Unix(0, i)

					event = append(event, events.K8sWatchEvent{
						Time:    tm.String(),
						Name:    v[1].(string),
						Message: v[2].(string),
						Type:    v[3].(string),
					})
					//fmt.Println(fmt.Sprintf("%s event:%s message:%s type:%s ", tm, v[1], v[2], v[3]))
				}
			}
		}
	}

	return
}
