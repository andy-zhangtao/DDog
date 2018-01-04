
.PHONY: build
name = ddog

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

run: build
	./$(name)

release: *.go *.md
	docker run -it --rm -v ${PWD}:/go/src/github.com/andy-zhangtao/DDog vikings/golang:1.9-onbuild /go/src/github.com/andy-zhangtao/DDog ddog
	ls -ltr bin/