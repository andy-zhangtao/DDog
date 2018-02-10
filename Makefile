
.PHONY: build
name = ddog

client: agent/*.go
	cd agent;go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)-agent
	mv agent/ddog-agent bin/ddog-agent

runclient: client
	bin/ddog-agent

build: *.go
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

run: build
	./$(name)

agent-release: agent/*.go
	docker run -it --rm -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:onbuild-v1.0.5 /go/src/github.com/andy-zhangtao/DDog/agent ddog-agent
	ls -ltr agent/bin
	@echo "############"
	@echo "ddog-agent build complete"

srv-release: *.go *.md
	docker run -it --rm -e DDOG_MONGO_ENDPOINT=127.0.0.1:27017 -e DDOG_MONGO_DB=cloud -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:unit-test-mongo /go/src/github.com/andy-zhangtao/DDog ddog -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)"
	ls -ltr bin/
	@echo "############"
	@echo "ddog build complete"

release: agent-release srv-release
	@echo "############"
