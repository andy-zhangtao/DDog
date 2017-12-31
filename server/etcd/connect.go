package etcd

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/client"
	"github.com/andy-zhangtao/DDog/const"
)

var endpoint = os.Getenv(_const.EnvEtcd)
var cliv3 *clientv3.Client
var cliv2 client.Client
var kapi client.KeysAPI
var isV2 bool
var debug = false

func check() error {
	if endpoint == "" {
		return errors.New(_const.EnvEtcdNotFound)
	}

	return nil
}

func init() {
	if os.Getenv("ETCDCTL_API") == "3" {
		if err := check(); err != nil {
			log.Println(err.Error())
			os.Exit(-1)
		}

		var err error

		ep := strings.Split(endpoint, ";")
		cliv3, err = clientv3.New(clientv3.Config{
			Endpoints:   ep,
			DialTimeout: 15 * time.Second,
		})

		if err != nil {
			log.Printf("Etcd[%s] Init Failed [%s]! \n", ep, err.Error())
			os.Exit(-1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		list, err := cliv3.MemberList(ctx)
		defer cancel()
		if err != nil {
			log.Printf("Etcd[%s] Get Member List Failed [%s]! \n", ep, err.Error())
			os.Exit(-1)
		}

		for _, m := range list.Members {
			log.Printf("[%s][%s]\n", m.Name, m.ClientURLs)
		}

	} else {
		isV2 = true
		if err := check(); err != nil {
			log.Println(err.Error())
			os.Exit(-1)
		}

		var err error

		ep := strings.Split(endpoint, ";")
		for i, e := range ep {
			ep[i] = "http://" + e
		}

		cliv2, err = client.New(client.Config{
			Endpoints: ep,
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: 15 * time.Second,
		})

		if err != nil {
			log.Printf("Etcd[%s] Init Failed [%s]! \n", ep, err.Error())
			os.Exit(-1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		version, err := cliv2.GetVersion(ctx)
		defer cancel()
		if err != nil {
			log.Printf("Etcd[%s] Get Member List Failed [%s]! \n", ep, err.Error())
			os.Exit(-1)
		}

		kapi = client.NewKeysAPI(cliv2)
		log.Printf("当前ETCD Server[%s] Cluster[%s]\n", version.Server, version.Cluster)
	}

	log.Printf("当前使用的Etcd版本为[%s]\n", os.Getenv("ETCDCTL_API"))
}

func Put(key, value string) error {
	if debug {
		log.Printf("Etcd/PUT操作记录Key=[%s]Value=[%s]\n", key, value)
	}

	if isV2 {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		_, err := kapi.Set(ctx, key, value, nil)
		defer cancel()
		if err != nil {
			return err
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		_, err := cliv3.Put(ctx, key, value)
		defer cancel()
		if err != nil {
			return err
		}
	}

	return nil
}

func Dele(key string) error {
	if debug {
		log.Printf("Etcd/Dele操作记录Key=[%s]\n", key)
	}
	if isV2 {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		_, err := kapi.Delete(ctx, key, nil)
		defer cancel()
		if err != nil {
			return err
		}
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		_, err := cliv3.Delete(ctx, key)
		defer cancel()
		if err != nil {
			return err
		}
	}

	return nil
}

func Get(key string, opts []string) (map[string]string, error) {
	if debug {
		log.Printf("Etcd/Get操作记录Key=[%s]opts=[%s]\n", key, opts)
	}
	data := make(map[string]string)
	if isV2 {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		var ops []clientv3.OpOption
		if len(opts) > 0 {
			for _, o := range opts {
				switch o {
				case "--from-key":
					ops = append(ops, clientv3.WithFromKey())
				}
			}
		}

		resp, err := kapi.Get(ctx, key, nil)
		defer cancel()
		if err != nil {
			return data, err
		}

		data[resp.Node.Key] = resp.Node.Value
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		var ops []clientv3.OpOption
		if len(opts) > 0 {
			for _, o := range opts {
				switch o {
				case "--from-key":
					ops = append(ops, clientv3.WithFromKey())
				}
			}
		}

		resp, err := cliv3.Get(ctx, key, ops...)
		defer cancel()
		if err != nil {
			return data, err
		}

		for _, ev := range resp.Kvs {
			data[string(ev.Key)] = string(ev.Value)
		}
	}

	if debug {
		log.Printf("Etcd Get数据为[%v]\n", data)
	}
	return data, nil
}

func SetDebug(debug bool) {
	debug = debug
}
