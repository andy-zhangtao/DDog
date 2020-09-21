. {
	debug
	errors
	whoami
	log
	
	proxy . /etc/resolv.conf {
        	except eqxiu.com eqxiu.cc
    }

	#proxy . /etc/resolv.conf {
    #		except eqxiu.cc
	#}	

	etcd eqxiu.com eqxiu.cc {
		stubzones
		path /pro
		endpoint http://192.168.0.8:2379 http://192.168.0.9:2379 http://192.168.0.18:2379
	}

}
