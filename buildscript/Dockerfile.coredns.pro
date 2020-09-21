FROM 	coredns/coredns
LABEL 	maintainer=ztao8607@gmail.com
COPY	corefile.pro /Corefile
EXPOSE  53/tcp 53/udp
