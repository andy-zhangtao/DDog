package _const

const Corefile = `
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