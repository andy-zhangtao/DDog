package _const

const Corefile = `
. {
	debug
	errors
	whoami
	log

	etcd {{.Domain}} {
		stubzones
		path /
		endpoint {{.Etcd}}
	}

	proxy . {{.Upstream}}
}
`