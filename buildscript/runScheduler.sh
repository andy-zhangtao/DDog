#!/bin/bash
version=0.6.9
docker run -it --rm --name scheduler -e DDOG_DEBUG="true" -e DDOG_MONGO_ENDPOINT="192.168.1.12:27017" -e DDOG_MONGO_NAME="cloud" -e DDOG_MONGO_PASSWD="password" -e DDOG_MONGO_DB="cloud" -e DDOG_NAME_SPACE="devenv" -e DDOG_REGION="sh" -e DDOG_NSQD_ENDPOINT="192.168.1.12:4150" -p 8000:8000 ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns:$version
