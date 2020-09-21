#!/bin/bash
version=0.6.8
docker run -it --rm -e DDOG_MONGO_NAME="cloud" -e DDOG_MONGO_PASSWD="password" -e DDOG_AGENT_NAME="DeployAgent" -e DDOG_AGENT_RETRI_NAMESPACE="devenv;" -e DDOG_MONGO_DB="cloud" -e DDOG_MONGO_ENDPOINT="192.168.1.12:27017" -e DDOG_NSQD_ENDPOINT="192.168.1.12:4150" ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:$version
