
.PHONY: build
name = ddog

build:
	go build -ldflags "-X main._VERSION_=$(shell date +%Y%m%d-%H%M%S)" -o $(name)

run: build
	./$(name)

release: *.go *.md
	docker build -t vikings/ddog .
	mv ./ddog ./bin
	docker rmi -f vikings/ddog