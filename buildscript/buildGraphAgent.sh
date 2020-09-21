#!/bin/sh
VERSION=v1.5.18
echo "cp $GOPATH/src/github.com/andy-zhangtao/DDog/bin/ddog-graph-srv caas-deploy-agent"
cp $GOPATH/src/github.com/andy-zhangtao/DDog/bin/ddog-graph-srv caas-deploy-agent 
echo "docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION -f Dockerfile.$1 ."
docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION -f Dockerfile.$1 .
#echo "docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$1-$VERSION"
#docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$1-$VERSION
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-deploy-agent:$VERSION 
