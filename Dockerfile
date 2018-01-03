FROM golang:1.9-alpine3.7
LABEL MAINTAINER=ztao@gmail.com
RUN apk --update add build-base
ADD . /go/src/github.com/andy-zhangtao/DDog
WORKDIR /go/src/github.com/andy-zhangtao/DDog
#COPY /etc/localtime /etc/localtime
RUN GOARCH=amd64 go build  -a -o bin/ddog