
.PHONY: build
name = ddog

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

run: build
	./$(name)

release: *.go *.md
	docker run -it --rm -e DDOG_MONGO_ENDPOINT=127.0.0.1:27017 -e DDOG_MONGO_DB=cloud -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:unit-test-mongo /go/src/github.com/andy-zhangtao/DDog ddog -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)"
	ls -ltr bin/
