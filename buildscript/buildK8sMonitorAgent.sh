#!/bin/sh
VERSION=1.3.5.1
echo "cp $GOPATH/src/github.com/andy-zhangtao/DDog/agent/bin/ddog-agent ."
cp $GOPATH/src/github.com/andy-zhangtao/DDog/agent/bin/ddog-agent .
echo "docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$VERSION -f Dockerfile.$1 ."
docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$VERSION -f Dockerfile.$1 .
echo "docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$1-$VERSION"
docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$1-$VERSION
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$1-$VERSION"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler-agent:$1-$VERSION 
