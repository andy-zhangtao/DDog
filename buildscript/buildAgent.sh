#!/bin/sh
VERSION=v1.1.0
echo "cp $GOPATH/src/github.com/andy-zhangtao/DDog/agent/bin/ddog-agent ."
cp $GOPATH/src/github.com/andy-zhangtao/DDog/agent/bin/ddog-agent .
echo "docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION -f Dockerfile.spider ."
docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION -f Dockerfile.spider .
echo "docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:$VERSION -f Dockerfile.agent ."
docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:$VERSION -f Dockerfile.agent .
echo "docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:latest"
docker tag  ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:latest
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:$VERSION
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:latest"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-spider-agent:latest
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:$VERSION"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-coredns-agent:$VERSION
