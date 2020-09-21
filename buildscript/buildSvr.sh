#!/bin/sh
VERSION=0.6.11
echo "cp $GOPATH/src/github.com/andy-zhangtao/DDog/bin/ddog  ."
cp $GOPATH/src/github.com/andy-zhangtao/DDog/bin/ddog .
echo "docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler:$VERSION -f Dockerfile ."
docker build -t ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler:$VERSION -f Dockerfile .
echo "docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler:$VERSION"
docker push ccr.ccs.tencentyun.com/eqxiu/eqxiu-caas-scheduler:$VERSION
