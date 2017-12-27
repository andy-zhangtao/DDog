// goinit.go 用来初始化coredns配置文件
package watch

import (
	"os"
	"errors"
	"github.com/andy-zhangtao/DDog/const"
	"log"
	"text/template"
)

const corefile = `
. {
	debug
	errors
	whoami
	log

	proxy . /etc/resolv.conf {
        	except {{.Domain}}
    }

	etcd {{.Domain}} {
		stubzones
		path /
		endpoint http://{{.Etcd}}
	}

}
`

func Go() {
	err := genConfigure()
	if err != nil {
		log.Panicln(err)
	}
}

// genConfigure 生成DNS配置文件
func genConfigure() error {
	type conf struct {
		Name string
		Etcd string
	}

	name := os.Getenv(_const.EnvDomain)
	if name == "" {
		return errors.New(_const.EnvDomainNotFound)
	}

	etcd := os.Getenv(_const.EnvEtcd)
	if etcd == "" {
		return errors.New(_const.EnvEtcdNotFound)
	}

	path := os.Getenv(_const.EnvConfPath)
	if path == "" {
		path = "/"
	}
	var mf = conf{
		Name: name,
		Etcd: etcd,
	}

	t := template.Must(template.New("makefile").Parse(corefile))

	file, err := os.Create(path + "corefile")
	if err != nil {
		return err
	}

	err = t.Execute(file, mf)
	if err != nil {
		return err
	}

	return nil
}
