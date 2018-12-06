
.PHONY: build
name = ddog
version = v1.1.1-DEBUG

client: agent/*.go
	rm -rf bin/ddog-agent
	cd agent;go build -ldflags "-X main._VERSION_=$(version) -X main._BUILD_=$(shell date +%Y%m%d_%H%M%S)" -o $(name)-agent
	mv agent/ddog-agent bin/ddog-agent

runclient: client
	bin/ddog-agent

graphql: agent/graphql/*.go
	cd agent/graphql;     go build -ldflags "-X main._VERSION_=$(version) -X main._BUILD_=$(shell date +%Y%m%d_%H%M%S)" -o ddog-graph-srv
	mv agent/graphql/ddog-graph-srv bin/ddog-graph-srv

graphql-release: agent/graphql/*.go
	cd agent/graphql;GOOS=linux GOARCH=amd64 go build -ldflags "-X main._VERSION_=$(version) -X main._BUILD_=$(shell date +%Y%m%d_%H%M%S)" -a -o ddog-graph-srv
	mv agent/graphql/ddog-graph-srv bin/ddog-graph-srv

build: *.go
	go build -ldflags "-X main._VERSION_=$(version) -X main._BUILD_=$(shell date +%Y%m%d_%H%M%S)" -o $(name)

run: build
	./$(name)

agent-release: agent/*.go
	docker run -it --rm -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:onbuild-v1.0.5 /go/src/github.com/andy-zhangtao/DDog/agent ddog-agent
	ls -ltr agent/bin
	@echo "############"
	@echo "ddog-agent build complete"

srv-release: *.go *.md
	docker run -it --rm -e DDOG_MONGO_ENDPOINT=127.0.0.1:27017 -e DDOG_MONGO_DB=cloud -e DDOG_NSQD_ENDPOINT=127.0.0.1:4150 -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:unitTest-v1.0.7 /go/src/github.com/andy-zhangtao/DDog ddog -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)"
	ls -ltr bin/
	@echo "############"
	@echo "ddog build complete"

release: agent-release srv-release
	@echo "############"
