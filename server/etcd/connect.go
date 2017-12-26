package etcd

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var endpoint = os.Getenv("DOG_ETCD_ENDPOINT")
var cli *clientv3.Client

func check() error {
	if endpoint == "" {
		return errors.New("DOG_ETCD_ENDPOINT Check Failed!")
	}

	return nil
}

func init() {
	if err := check(); err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	var err error

	ep := strings.Split(endpoint, ";")
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   ep,
		DialTimeout: 15 * time.Second,
	})
	if err != nil {
		log.Println("Etcd Init Failed! " + err.Error())
		os.Exit(-1)
	}
}

func Put(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := cli.Put(ctx, key, value)
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}

func Dele(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	_, err := cli.Delete(ctx, key)
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}
